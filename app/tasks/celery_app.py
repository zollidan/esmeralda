from celery import Celery

from app.parser.testparser.parser import parse_data_test

# Создание экземпляра Celery
celery_app = Celery(
    "tasks",
    broker="redis://localhost:6379/0",  # Брокер
    backend="redis://localhost:6379/1",
)

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