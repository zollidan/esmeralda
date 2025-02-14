import requests
from fastapi import APIRouter

from app.config import settings

tasks_router = APIRouter(prefix='/api/tasks', tags=['tasks'])

@tasks_router.get("/info/{task_id}")
def get_task_status(task_id: str):
    response = requests.get(f"{settings.FLOWER_URL}/api/task/info/{task_id}") #, auth=HTTPBasicAuth(settings.FLOWER_USER, settings.FLOWER_PASSWORD)
    return response.json()

@tasks_router.get("/all")
def get_flower_tasks():
    response = requests.get(f"{settings.FLOWER_URL}/api/tasks") #, auth=HTTPBasicAuth(settings.FLOWER_USER, settings.FLOWER_PASSWORD)
    return response.json()
