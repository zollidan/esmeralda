import asyncio
import datetime
import io
import concurrent.futures
import pandas as pd
from celery import Celery
from sqlalchemy import select
from asgiref.sync import async_to_sync

from config import settings
from database import SessionLocal 
from models import ParsingTask
from boto3_utils import upload_to_s3
from sofascore.parser import main as sofascore_parser

celery_app = Celery('tasks', broker=settings.CELERY_BROKER_URL)

async def _run_parsing_logic(task_id):
    """
    Внутренняя асинхронная логика для работы с БД и вызова парсера
    """
    async with SessionLocal() as db:
        result = await db.execute(select(ParsingTask).where(ParsingTask.id == task_id))
        task = result.scalar_one_or_none()

        if not task:
            print(f"Task {task_id} not found")
            return

        try:
            task.status = 'PROCESSING'
            await db.commit()

            time_tomorrow = datetime.datetime.now() + datetime.timedelta(days=1)
            date_str = time_tomorrow.strftime("%Y-%m-%d")

            loop = asyncio.get_running_loop()
            with concurrent.futures.ThreadPoolExecutor() as pool:
                df = await loop.run_in_executor(
                    pool, 
                    sofascore_parser, 
                    date_str
                )

            excel_buffer = io.BytesIO()
            df.to_excel(excel_buffer, index=False)
            excel_buffer.seek(0)
            
            file_name = f"results/result_{task_id}.xlsx"
            file_url = upload_to_s3(excel_buffer.getvalue(), file_name)

            task.status = 'DONE'
            task.result_path = file_url

        except Exception as e:
            await db.rollback()
            task.status = 'ERROR'
            task.error_log = str(e)
            print(f"Error in task {task_id}: {e}")
            
        finally:
            await db.commit()

@celery_app.task(bind=True)
def run_parsing_task(self, task_id):
    """
    Точка входа Celery.
    Так как Celery не поддерживает async def, используем async_to_sync.
    """
    return async_to_sync(_run_parsing_logic)(task_id)