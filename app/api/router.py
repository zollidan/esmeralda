import uuid
from fastapi import APIRouter, Response
import requests
from celery.result import AsyncResult
from app.tasks.celery_app import celery_app, soccerway_parser_task
from app.config import settings, s3_client


router = APIRouter(prefix='/api', tags=['API'])

    
# Маршрут celery soccerway
@router.post('/parser/soccerway')
def run_soccerway_url_method(date: str):
    file_name = f'soccerway-{date}-{str(uuid.uuid4())}.xls'
    
    # send_file_record_s3(file_name)

    # os.remove(file_name)

    
    task = soccerway_parser_task.delay(user_date=date, my_file_name=file_name)  # Отправляю задачу в Celery # тестовая задача parse_data_test()
    
    
    """ тут отправляется сообещние в телеграм"""
    
    
    return {
        "message": "soccerway work started",
        "task_id": task.id,  # Возвращаем ID задачи для проверки статуса
        "status": 200
    }

# Маршрут для проверки статуса задачи
@router.get('/{task_id}')
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


@router.get('/s3/files')
def get_all_s3_files():
    
    return s3_client.list_objects(Bucket=settings.AWS_BUCKET)['Contents']


@router.get("/s3/files/{file_id}")
async def get_file_by_id(file_id:str):
    
    # Получить объект
    get_object_response = s3_client.get_object(Bucket=settings.AWS_BUCKET,Key=file_id)
    file_binary = get_object_response['Body'].read()
    
    # Получаем Content-Type из метаданных S3
    content_type = get_object_response['ContentType']
    
    return Response(
        content=file_binary,
        media_type=content_type
    )