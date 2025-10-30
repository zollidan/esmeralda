import pika
import json
import logging
from typing import Dict, Any, Optional

logger = logging.getLogger(__name__)


class ResultPublisher:
    """Класс для публикации результатов парсинга в RabbitMQ"""

    def __init__(self, channel: pika.channel.Channel, results_queue: str):
        self.channel = channel
        self.results_queue = results_queue

        # Объявляем очередь для результатов (durable для сохранности)
        self.channel.queue_declare(queue=self.results_queue, durable=True)

    def publish_success(self, task_id: str, file_path: str, metadata: Optional[Dict[str, Any]] = None):
        """
        Публикует успешный результат

        Args:
            task_id: ID задачи
            file_path: Путь к файлу в S3
            metadata: Дополнительные метаданные (опционально)
        """
        result_data = {
            'task_id': task_id,
            'status': 'completed',
            'file_path': file_path,
            'metadata': metadata or {}
        }

        self._publish(result_data)
        logger.info(f"Результат задачи {task_id} успешно опубликован")

    def publish_failure(self, task_id: str, error: str, metadata: Optional[Dict[str, Any]] = None):
        """
        Публикует ошибку выполнения задачи

        Args:
            task_id: ID задачи
            error: Описание ошибки
            metadata: Дополнительные метаданные (опционально)
        """
        result_data = {
            'task_id': task_id,
            'status': 'failed',
            'error': error,
            'file_path': None,
            'metadata': metadata or {}
        }

        self._publish(result_data)
        logger.error(f"Ошибка задачи {task_id} опубликована: {error}")

    def _publish(self, data: Dict[str, Any]):
        """Внутренний метод для публикации в RabbitMQ"""
        try:
            self.channel.basic_publish(
                exchange='',
                routing_key=self.results_queue,
                body=json.dumps(data, ensure_ascii=False),
                properties=pika.BasicProperties(
                    delivery_mode=2,  # persistent message
                    content_type='application/json'
                )
            )
        except Exception as e:
            logger.error(f"Ошибка публикации в RabbitMQ: {e}")
            raise
