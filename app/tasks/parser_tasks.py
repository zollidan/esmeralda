from app.parser.testparser.parser import parse_data_test
from .celery_app import celery_app


# Функция для выполнения через Celery
@celery_app.task
def parse_data_task(date):
    return parse_data_test(date)