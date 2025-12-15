import time
import io  # Библиотека для работы с файлами в памяти
import pandas as pd
from celery import shared_task
from django.core.files.base import ContentFile
from .models import ParsingTask

@shared_task
def run_parsing_task(task_id):
    task = ParsingTask.objects.get(id=task_id)
    task.status = 'PROCESSING'
    task.save()

    try:
        time.sleep(5) 
        
        data = [[1, 2, 3], [4, 5, 6], [7, 8, 9], [10, 11, 12]]
        df = pd.DataFrame(data=data, columns=['Col1', 'Col2', 'Col3'])
        
        csv_buffer = io.StringIO()
        
        df.to_csv(csv_buffer, index=False)
        
        file_content = csv_buffer.getvalue()

        file_name = f"result_{task_id}.csv"
        
        task.result_file.save(file_name, ContentFile(file_content.encode('utf-8')))
        
        task.status = 'DONE'
        
    except Exception as e:
        task.status = 'ERROR'
        task.error_log = str(e)
    finally:
        task.save()