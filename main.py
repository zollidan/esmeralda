from fastapi import FastAPI, Depends, Request, Form, HTTPException
from fastapi.responses import RedirectResponse, StreamingResponse
from fastapi.templating import Jinja2Templates
from fastapi.staticfiles import StaticFiles
from sqlalchemy.orm import Session
from boto3_utils import get_s3_client
import models
from database import engine, get_db
from tasks import run_parsing_task
from config import settings
# Создаем таблицы
models.Base.metadata.create_all(bind=engine)

app = FastAPI()
templates = Jinja2Templates(directory="templates")

@app.get("/")
def index(request: Request, db: Session = Depends(get_db)):
    tasks = db.query(models.ParsingTask).order_by(models.ParsingTask.id.desc()).all()
    return templates.TemplateResponse("index.html", {"request": request, "tasks": tasks})

@app.post("/run")
def start_task(db: Session = Depends(get_db)):
    new_task = models.ParsingTask(status="PENDING")
    db.add(new_task)
    db.commit()
    db.refresh(new_task) 
    
    # Запускаем Celery
    run_parsing_task.delay(new_task.id)
    
    return RedirectResponse(url="/", status_code=303)

@app.get("/download/{task_id}")
def download_result(task_id: int, db: Session = Depends(get_db)):
    task = db.query(models.ParsingTask).filter(models.ParsingTask.id == task_id).first()
    
    if not task or not task.result_path:
        raise HTTPException(status_code=404, detail="Файл не найден или задача еще не завершена")

    bucket_name = settings.AWS_STORAGE_BUCKET_NAME
    s3_key = task.result_path.split(f"{bucket_name}/")[-1]

    s3 = get_s3_client()
    
    try:
        s3_object = s3.get_object(Bucket=bucket_name, Key=s3_key)
        
        return StreamingResponse(
            s3_object['Body'],
            media_type='text/csv',
            headers={
                "Content-Disposition": f"attachment; filename=result_{task_id}.csv"
            }
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Ошибка S3: {str(e)}")