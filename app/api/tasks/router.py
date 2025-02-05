from celery.result import AsyncResult
from fastapi import APIRouter

from app.tasks.celery_app import celery_app

tasks_router = APIRouter(prefix='/api/tasks', tags=['tasks'])

@tasks_router.get("/{task_id}")
def get_task_status(task_id: str):
    """Получение статуса и результата конкретной задачи"""
    task_result = AsyncResult(task_id, app=celery_app)
    return {
        "task_id": task_id,
        "status": task_result.status,
        "result": task_result.result
    }


