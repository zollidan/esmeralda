import os
from celery import Celery

from app.parser.marafon.main import run_marafon_parser
from app.parser.soccerway.parser import run_soccerway

# Создание экземпляра Celery
celery_app = Celery(
    "tasks",
)

celery_app.conf.broker_url = os.environ.get("CELERY_BROKER_URL", "redis://localhost:6379")
celery_app.conf.result_backend = os.environ.get("CELERY_RESULT_BACKEND", "redis://localhost:6379")


celery_app.conf.update(
    task_serializer='json',
    result_serializer='json',
    accept_content=['json'],
)

celery_app.autodiscover_tasks(['app.tasks'])

@celery_app.task
def soccerway_parser_task(user_date, my_file_name):
    return run_soccerway(user_date, my_file_name)

@celery_app.task
def marafon_parser_task(file_name):
    return run_marafon_parser(file_name)