from celery import Celery

# Создание экземпляра Celery
celery_app = Celery(
    "tasks",
    broker="redis://localhost:6379/0",  # Брокер
    backend="redis://localhost:6379/1",  # Backend
)

# Загрузка конфигурации (опционально)
celery_app.conf.update(
    task_serializer="json",
    result_serializer="json",
    accept_content=["json"],
    timezone="UTC",
    enable_utc=True,
)