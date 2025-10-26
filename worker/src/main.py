import pika
import json
from time import sleep


def parser(task_data):
    """Ваш парсер"""
    print(f"Начинаю парсинг для задачи: {task_data.get('task_name')}")

    for i in range(10):
        print(f"Парсинг... шаг {i+1}/10")
        sleep(3)

    print("Парсинг завершен!")


def callback(ch, method, properties, body):
    """Обработчик сообщений из RabbitMQ"""
    try:
        # Парсим данные из Go
        task = json.loads(body)
        print(f"\n[✓] Получена задача из RabbitMQ:")
        print(f"    Task ID: {task.get('task_id')}")

        # Запускаем парсер
        parser(task)

        # Подтверждаем успешную обработку
        ch.basic_ack(delivery_tag=method.delivery_tag)
        print("[✓] Задача успешно обработана\n")

    except Exception as e:
        print(f"[✗] Ошибка при обработке: {e}")
        # Отклоняем сообщение (не возвращаем в очередь)
        ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)


def main():
    print("=" * 50)
    print("Worker запущен и ожидает задачи...")
    print("=" * 50)

    # Подключение к RabbitMQ
    connection = pika.BlockingConnection(
        pika.ConnectionParameters(
            host='localhost',
            heartbeat=600,  # таймаут для длительных задач
            blocked_connection_timeout=300
        )
    )
    channel = connection.channel()

    # Декларируем очередь (имя должно совпадать с Go)
    channel.queue_declare(queue='tasks_queue', durable=True)

    # Обрабатываем по одному сообщению за раз
    channel.basic_qos(prefetch_count=1)

    # Подписываемся на очередь
    channel.basic_consume(
        queue='tasks_queue',
        on_message_callback=callback,
        auto_ack=False  # ручное подтверждение
    )

    print("\n[*] Ожидание сообщений. Для выхода: CTRL+C\n")

    try:
        channel.start_consuming()
    except KeyboardInterrupt:
        print("\n[!] Остановка worker...")
        channel.stop_consuming()

    connection.close()
    print("[✓] Worker остановлен")


if __name__ == "__main__":
    main()
