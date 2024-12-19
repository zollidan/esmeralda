import asyncio
from fastapi import APIRouter, BackgroundTasks
from pydantic import BaseModel
import requests
from celery import Celery
from celery.result import AsyncResult
from app.celery_app import celery_app
from app.api.dao import FileDAO

from app.parser.parsers_wrapper import soccerway_wrapper
from app.parser.server_parser_funcions import get_files_s3
from app.parser.testparser.parser import parse_data_test


router = APIRouter(prefix='/api', tags=['API'])

class File(BaseModel):
    name:str
    url:str
    
# Маршрут FastAPI
@router.post('/parser/soccerway')
def run_soccerway_url_method(date: str):
    task = parse_data_test.delay(date)  # Отправляем задачу в Celery # тестовая задача parse_data_test()
    return {
        "message": "soccerway work started",
        "task_id": task.id,  # Возвращаем ID задачи для проверки статуса
        "status": 200
    }

# Маршрут для проверки статуса задачи
@router.get('/parser/soccerway/status/{task_id}')
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
    
@router.get("/tasks/")
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


@router.get('/parser/soccerway/connection-test')
def run_soccerway_test_connection():
    url = "https://ru.soccerway.com/a/block_h2h_matches?block_id=page_match_1_block_h2hsection_head2head_7_block_h2h_matches_1&action=changePage&callback_params=%7B%22page%22%3A+-1%2C+%22block_service_id%22%3A+%22match_h2h_comparison_block_h2hmatches%22%2C+%22team_A_id%22%3A+43000%2C+%22team_B_id%22%3A+43009%7D&params=%7B%22page%22%3A+0%7D"

    response = requests.get(url, headers={'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:102.0) Gecko/20100101 Firefox/131.0'})
    
    # Проверка на успешный ответ
    if response.status_code == 200:
        return {"status": "success", "data": response.json()}
    else:
        return {"status": "error", "status_code": response.status_code, "message": response.text}


@router.get('/files')
def get_files():
    
    files = get_files_s3()
    
    return files
    
    
    
@router.get('/files/sql')
async def get_files_sql():
    
    files = await FileDAO.find_all()
    
    return files

@router.post('/files/sql', status_code=201)
async def create_file_sql(file: File):
    
    await FileDAO.add(
        name=file.name,
        url=file.url
    )
    
    return {"file": file}