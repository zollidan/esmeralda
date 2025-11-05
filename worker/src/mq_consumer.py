import pika
import json
import logging
import uuid
from typing import Dict, Any
from config import get_settings
from parser.s3_uploader import S3Uploader
from parser.result_publisher import ResultPublisher
from worker.src.parser.core.temp_parser import parse_soccerway, parse_another_site

# Инициализируем настройки
settings = get_settings()

logging.basicConfig(
    level=settings.app.log_level,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Словарь парсеров - легко добавлять новые
PARSERS = {
    'soccerway': parse_soccerway,
    'another_site': parse_another_site,
    # Добавляй свои парсеры сюда
}


def process_task(task_data: Dict[str, Any], s3_uploader: S3Uploader, publisher: ResultPublisher):
    """
    Основная функция обработки задачи

    Args:
        task_data: Данные задачи из RabbitMQ
        s3_uploader: Экземпляр S3Uploader
        publisher: Экземпляр ResultPublisher
    """
    task_id = task_data.get('task_id', 'unknown')
    parser_id = task_data.get('parser_id', 'default')

    logger.info(f"Обработка задачи {task_id}, парсер: {parser_id}")

    try:
        # 1. Получаем нужный парсер
        parser_func = PARSERS.get(parser_id)
        if not parser_func:
            raise ValueError(f"Неизвестный тип парсера: {parser_id}")

        # 2. Выполняем парсинг
        df = parser_func(task_data)

        if df is None or df.empty:
            raise ValueError("Парсер вернул пустой результат")

        # 3. Генерируем имя файла
        filename = f"{parser_id}-{task_id}-{uuid.uuid4()}.xlsx"

        # 4. Загружаем в S3
        upload_result = s3_uploader.upload_dataframe(df, filename)

        # 5. Публикуем результат ТОЛЬКО после успешной загрузки
        if upload_result['success']:
            publisher.publish_success(
                task_id=task_id,
                file_path=upload_result['path'],
                metadata={
                    'parser_id': parser_id,
                    'rows_count': len(df),
                    'columns': list(df.columns)
                }
            )
            logger.info(f"✅ Задача {task_id} успешно завершена")
        else:
            # Если загрузка в S3 не удалась - публикуем ошибку
            error_msg = upload_result.get(
                'error', 'Неизвестная ошибка загрузки в S3')
            publisher.publish_failure(
                task_id=task_id,
                error=error_msg,
                metadata={'parser_id': parser_id}
            )
            logger.error(f"❌ Задача {task_id} провалена: {error_msg}")

    except Exception as e:
        # Любая ошибка парсинга - публикуем в RabbitMQ
        error_msg = f"Ошибка обработки задачи: {str(e)}"
        logger.exception(f"❌ Задача {task_id} провалена с исключением")

        publisher.publish_failure(
            task_id=task_id,
            error=error_msg,
            metadata={'parser_id': parser_id}
        )


def main():
    """Главная функция - запуск consumer"""

    logger.info(f"🚀 Запуск worker: {settings.app.worker_name}")
    logger.info(f"📋 Уровень логирования: {settings.app.log_level}")

    logger.info("Инициализация S3 uploader...")
    s3_uploader = S3Uploader(
        endpoint_url=settings.s3.endpoint_url,
        access_key=settings.s3.access_key,
        secret_key=settings.s3.secret_key,
        bucket_name=settings.s3.bucket_name
    )

    logger.info("Подключение к RabbitMQ...")
    credentials = pika.PlainCredentials(
        settings.rabbitmq.username,
        settings.rabbitmq.password
    )
    connection = pika.BlockingConnection(
        pika.ConnectionParameters(
            host=settings.rabbitmq.host,
            port=settings.rabbitmq.port,
            credentials=credentials
        )
    )
    channel = connection.channel()

    # Объявляем очередь для задач
    channel.queue_declare(queue=settings.rabbitmq.tasks_queue, durable=True)

    # Создаем publisher для результатов
    publisher = ResultPublisher(channel, settings.rabbitmq.results_queue)

    def callback(ch, method, properties, body):
        """Callback для обработки сообщений из очереди"""
        try:
            task_data = json.loads(body)
            process_task(task_data, s3_uploader, publisher)
        except json.JSONDecodeError as e:
            logger.error(f"Ошибка парсинга JSON: {e}")
        except Exception as e:
            logger.exception(f"Неожиданная ошибка в callback: {e}")
        finally:
            # Подтверждаем обработку сообщения
            ch.basic_ack(delivery_tag=method.delivery_tag)

    # Настраиваем prefetch - обрабатываем заданное количество задач
    channel.basic_qos(prefetch_count=settings.app.prefetch_count)

    # Начинаем слушать очередь
    channel.basic_consume(
        queue=settings.rabbitmq.tasks_queue,
        on_message_callback=callback
    )

    logger.info(
        f"⏳ Ожидание задач из '{settings.rabbitmq.tasks_queue}' (prefetch: {settings.app.prefetch_count})...")

    try:
        channel.start_consuming()
    except KeyboardInterrupt:
        logger.info("Остановка worker...")
        channel.stop_consuming()
    finally:
        connection.close()


if __name__ == '__main__':
    main()
