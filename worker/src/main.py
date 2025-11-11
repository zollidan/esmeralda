import os
import sys
import signal

from worker.src.logger import logger
from worker.src.mq import MQ
from worker.src.settings import settings
from worker.src.storage import S3Storage


# Global instances for signal handler
mq_instance = None
storage_instance = None


def signal_handler(sig, frame):
    """Handle shutdown signals gracefully"""
    logger.info(f"Received signal {sig}, shutting down gracefully...")
    if mq_instance:
        mq_instance.stop()
    sys.exit(0)


def main():
    """Main entry point for the worker application"""
    global mq_instance, storage_instance

    try:
        logger.info("=" * 50)
        logger.info("Worker microservice starting...")
        logger.info(f"Task queue: {settings.RABBITMQ_TASKS_QUEUE}")
        logger.info(f"Result queue: {settings.RABBITMQ_RESULTS_QUEUE}")
        logger.info("=" * 50)

        # Setup signal handlers for graceful shutdown
        signal.signal(signal.SIGINT, signal_handler)
        signal.signal(signal.SIGTERM, signal_handler)

        # Initialize storage
        storage_instance = S3Storage()

        # Initialize and start MQ consumer
        mq_instance = MQ(storage=storage_instance)  # Передаем storage в MQ
        logger.info("Worker microservice started successfully")

        # Start consuming (this blocks until interrupted)
        mq_instance.start_consuming()

    except Exception as e:
        logger.error(f"Fatal error in main: {e}")
        if mq_instance:
            mq_instance.stop()
        sys.exit(1)


if __name__ == '__main__':
    main()