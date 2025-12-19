import io
import pandas as pd
from celery import Celery
from config import settings
from database import SessionLocal
from models import ParsingTask
from boto3_utils import upload_to_s3

celery_app = Celery('tasks', broker=settings.CELERY_BROKER_URL)

@celery_app.task(bind=True)
def run_parsing_task(self, task_id):
    db = SessionLocal()
    task = db.query(ParsingTask).filter(ParsingTask.id == task_id).first()
    
    try:
        task.status = 'PROCESSING'
        db.commit()

        df = pd.DataFrame([[1, 2, 3]], columns=['A', 'B', 'C'])
        
        csv_buffer = io.BytesIO()
        df.to_csv(csv_buffer, index=False)
        csv_buffer.seek(0)

        file_name = f"results/result_{task_id}.csv"
        file_url = upload_to_s3(csv_buffer.getvalue(), file_name)

        task.status = 'DONE'
        task.result_path = file_url
        
    except Exception as e:
        task.status = 'ERROR'
        task.error_log = str(e)
    finally:
        db.commit()
        db.close()