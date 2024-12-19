from celery import Celery

# Создание экземпляра Celery
celery_app = Celery(
    "tasks",
    broker="redis://localhost:6379/0",  # Брокер
    backend="redis://localhost:6379/1",  # Backend
)

celery_app.conf.update(
    task_serializer='json',
    result_serializer='json',
    accept_content=['json'],
)

# Автозагрузка задач из указанных модулей
celery_app.autodiscover_tasks(['app.tasks'])