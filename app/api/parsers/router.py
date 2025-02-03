import uuid
from fastapi import APIRouter
from app.tasks.celery_app import celery_app, soccerway_parser_task, marafon_parser_task, soccerway_old_parser_task

parser_router = APIRouter(prefix='/api/parser', tags=['parsers'])

@parser_router.post('/soccerway')
def run_soccerway_url_method(date: str):
    file_name = f'soccerway-{date}-{str(uuid.uuid4())}.xls'
    
    task = soccerway_parser_task.delay(user_date=date, my_file_name=file_name)
    
    
    """ тут отправляется сообещние в телеграм"""
    
    
    return {
        "message": "soccerway work started",
        "task_id": task.id,
        "status": 200
    }

@parser_router.post('/soccerway_old')
def run_soccerway_old_url_method(date: str):
    file_name = f'soccerway_old-{date}-{str(uuid.uuid4())}.xls'
    task = soccerway_old_parser_task.delay(user_date=date, my_file_name=file_name)

    return {
        "message": "soccerway work started",
        "task_id": task.id,
        "status": 200
    }

@parser_router.post('/marafon')
def run_marafon_url_method():
    try:
        file_name = f'marafon-{str(uuid.uuid4())}.xlsx'

        task = marafon_parser_task.delay(file_name)

        return {
            "message": "marafon work started",
            "task_id": task.id,
            "status": 200
        }
    except Exception as e:
        return {"error", e}

