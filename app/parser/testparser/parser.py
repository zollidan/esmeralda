from app.parser.testparser.parser_funcs import parser_loop

from celery import Celery

# Настройка Celery
test_parser_celery = Celery(
    'tasks',  # Имя приложения
    broker='redis://localhost:6379/0',  # URL вашего Redis-брокера
    backend='redis://localhost:6379/0'
)

# Функция для выполнения через Celery
@test_parser_celery.task
def parse_data_test(self):
    print("[Test Parser start working]")
    parser_loop()  # Здесь вызывается ваша функция парсера