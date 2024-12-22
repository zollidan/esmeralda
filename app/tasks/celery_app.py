import os
from celery import Celery

from app.parser.soccerway.parser import run_soccerway
from app.parser.testparser.parser import parse_data_test

# Создание экземпляра Celery
celery_app = Celery(
    "tasks",
    # broker="redis://localhost:6379/0",  # Брокер
    # backend="redis://localhost:6379/1",
)

celery_app.conf.broker_url = os.environ.get("CELERY_BROKER_URL", "redis://localhost:6379")
celery_app.conf.result_backend = os.environ.get("CELERY_RESULT_BACKEND", "redis://localhost:6379")


celery_app.conf.update(
    task_serializer='json',
    result_serializer='json',
    accept_content=['json'],
)

# Автозагрузка задач из указанных модулей
celery_app.autodiscover_tasks(['app.tasks'])

@celery_app.task
def parse_data_task(date):
    return parse_data_test(date)

@celery_app.task
def soccerway_parser_task(user_date, my_file_name):
    return run_soccerway(user_date, my_file_name)