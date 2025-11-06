import json
import pika
from typing import Optional

from worker.src.parsers import parser
from worker.src.settings import settings
from worker.src.logger import logger


class MQ:
    """Message Queue handler for RabbitMQ with proper connection management"""

    def __init__(self):
        self.connection: Optional[pika.BlockingConnection] = None
        self.channel: Optional[pika.adapters.blocking_connection.BlockingChannel] = None
        self.result_channel: Optional[pika.adapters.blocking_connection.BlockingChannel] = None
        self._setup_connections()

    def _setup_connections(self):
        """Setup RabbitMQ connections and channels"""
        try:
            # Main connection for consuming
            credentials = pika.PlainCredentials(
                settings.RABBITMQ_USER,
                settings.RABBITMQ_PASSWORD
            )
            self.connection = pika.BlockingConnection(
                pika.ConnectionParameters(
                    host=settings.RABBITMQ_HOST,
                    port=settings.RABBITMQ_PORT,
                    virtual_host=settings.RABBITMQ_VHOST,
                    credentials=credentials,
                    heartbeat=settings.WORKER_HEARTBEAT,
                    blocked_connection_timeout=settings.WORKER_BLOCKED_CONNECTION_TIMEOUT
                )
            )
            self.channel = self.connection.channel()

            # Declare task queue
            self.channel.queue_declare(
                queue=settings.RABBITMQ_TASKS_QUEUE,
                durable=True
            )

            # Set QoS to process one message at a time
            self.channel.basic_qos(prefetch_count=settings.WORKER_PREFETCH_COUNT)

            # Setup result channel (reuse connection)
            self.result_channel = self.connection.channel()
            self.result_channel.queue_declare(
                queue=settings.RABBITMQ_RESULTS_QUEUE,
                durable=True
            )

            logger.info("RabbitMQ connections established successfully")
        except Exception as e:
            logger.error(f"Failed to setup RabbitMQ connections: {e}")
            raise

    def start_consuming(self):
        """Start consuming messages from the task queue"""
        try:
            self.channel.basic_consume(
                queue=settings.RABBITMQ_TASKS_QUEUE,
                on_message_callback=self.callback,
                auto_ack=False  # Manual ACK for reliability
            )

            logger.info('Waiting for messages. To exit press CTRL+C')
            self.channel.start_consuming()
        except KeyboardInterrupt:
            logger.info("Stopping consumer...")
            self.stop()
        except Exception as e:
            logger.error(f"Error while consuming: {e}")
            self.stop()
            raise

    def callback(self, ch, method, properties, body):
        """Process incoming task message"""
        task_id = None
        try:
            # Parse JSON message
            message = json.loads(body.decode())
            task_id = message.get("task_id")
            task_type = message.get("task_type")
            params = message.get("params", {})

            logger.info(f"Received task: {task_id}, type: {task_type}")

            # Process task
            df = parser(task_type=task_type, params=params)

            # Save to file (temporary)
            output_file = f"parsed_games_{task_id}.xlsx"
            df.to_excel(output_file, index=False)
            logger.info(f"Data parsed and saved to {output_file}")

            # TODO: Save to S3
            # s3_key = storage.upload(output_file, task_id)

            # Send success result
            self.send_result({
                "task_id": task_id,
                "status": "success",
                "output_file": output_file,
                # "s3_key": s3_key
            })

            # ACK message only after successful processing
            ch.basic_ack(delivery_tag=method.delivery_tag)
            logger.info(f"Task {task_id} completed successfully")

        except json.JSONDecodeError as e:
            logger.error(f"Invalid JSON message: {e}")
            # Reject and requeue malformed messages
            ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)

        except Exception as e:
            logger.error(f"Error processing task {task_id}: {e}")
            # Send failure result
            if task_id:
                self.send_result({
                    "task_id": task_id,
                    "status": "failed",
                    "error": str(e)
                })
            # Don't ACK, let RabbitMQ requeue the message
            ch.basic_nack(delivery_tag=method.delivery_tag, requeue=True)

    def send_result(self, message: dict):
        """Send result message to result queue"""
        try:
            self.result_channel.basic_publish(
                exchange='',
                routing_key=settings.RABBITMQ_RESULTS_QUEUE,
                body=json.dumps(message),
                properties=pika.BasicProperties(
                    delivery_mode=2,  # Make message persistent
                    content_type='application/json'
                )
            )
            logger.info(f"Sent result: {message}")
        except Exception as e:
            logger.error(f"Failed to send result: {e}")
            # Try to reconnect if channel is closed
            if self.result_channel.is_closed:
                logger.warning("Result channel closed, attempting to reconnect...")
                self._setup_connections()
                self.result_channel.basic_publish(
                    exchange='',
                    routing_key=settings.RABBITMQ_RESULTS_QUEUE,
                    body=json.dumps(message),
                    properties=pika.BasicProperties(
                        delivery_mode=2,
                        content_type='application/json'
                    )
                )

    def stop(self):
        """Gracefully stop consuming and close connections"""
        try:
            if self.channel and self.channel.is_open:
                self.channel.stop_consuming()
                self.channel.close()
            if self.connection and self.connection.is_open:
                self.connection.close()
            logger.info("RabbitMQ connections closed")
        except Exception as e:
            logger.error(f"Error closing connections: {e}")
