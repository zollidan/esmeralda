from celery.result import AsyncResult
from fastapi import APIRouter

from app.tasks.celery_app import celery_app

tasks_router = APIRouter(prefix='/api/tasks', tags=['tasks'])

@tasks_router.get("/all")
async def get_task_statuses():
    """
    Возвращает список задач Celery для фронтенда.

    тут временные меры, надо будет лучше переделать

    """
    # Получаем все задачи из backend Redis
    task_ids = celery_app.backend.client.keys("celery-task-meta-*")
    tasks = []

    for task_id in task_ids:
        # Извлекаем ID задачи из ключа Redis
        clean_id = task_id.decode().replace("celery-task-meta-", "")
        async_result = AsyncResult(clean_id, app=celery_app)

        # Собираем информацию о задаче
        tasks.append({
            "task_id": clean_id,
            "status": async_result.status,
            "result": async_result.result if async_result.ready() else None,
        })

    return {"tasks": tasks}

@tasks_router.get('/{task_id}')
def check_task_status(task_id: str):
    task_result = AsyncResult(task_id, app=celery_app)
    if task_result.state == 'PENDING':
        return {"task_id": task_id, "status": "PENDING", "result": None}
    elif task_result.state == 'SUCCESS':
        return {"task_id": task_id, "status": "SUCCESS", "result": task_result.result}
    elif task_result.state == 'FAILURE':
        return {"task_id": task_id, "status": "FAILURE", "result": str(task_result.info)}
    else:
        return {"task_id": task_id, "status": task_result.state, "result": None}


