import datetime
import io
import pandas as pd
from celery import Celery
from config import settings
from database import SessionLocal
from models import ParsingTask
from boto3_utils import upload_to_s3
from sofascore.parser import main as sofascore_parser
celery_app = Celery('tasks', broker=settings.CELERY_BROKER_URL)

@celery_app.task(bind=True)
def run_parsing_task(self, task_id):
    with SessionLocal() as db:
        task = db.query(ParsingTask).filter(ParsingTask.id == task_id).first()
        
        try:
            task.status = 'PROCESSING'
            db.commit()
            excel_buffer = io.BytesIO()
            time_tommorow = datetime.datetime.now() + datetime.timedelta(days=1)
            df = sofascore_parser(time_tommorow.strftime("%Y-%m-%d"))        
            df.to_excel(excel_buffer, index=False)
            excel_buffer.seek(0)
            file_name = f"results/result_{task_id}.xlsx"
            file_url = upload_to_s3(excel_buffer.getvalue(), file_name)

            task.status = 'DONE'
            task.result_path = file_url
        except Exception as e:
            task.status = 'ERROR'
            task.error_log = str(e)
        finally:
            db.commit()
            db.close()