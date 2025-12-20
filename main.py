from fastapi import FastAPI, Depends, Request, HTTPException
from fastapi.concurrency import asynccontextmanager
from fastapi.responses import RedirectResponse, StreamingResponse
from fastapi.templating import Jinja2Templates
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession
from boto3_utils import get_s3_client
import models
from database import engine, get_db, Base
from tasks import run_parsing_task
from config import settings

        
@asynccontextmanager
async def lifespan(app: FastAPI):
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
    yield

app = FastAPI(lifespan=lifespan)
templates = Jinja2Templates(directory="templates")


@app.get("/")
async def index(request: Request, db: AsyncSession = Depends(get_db)):
    stmt = select(models.ParsingTask).order_by(models.ParsingTask.id.desc())
    result = await db.execute(stmt)
    tasks = result.scalars().all()
    return templates.TemplateResponse("index.html", {"request": request, "tasks": tasks})


@app.post("/run")
async def start_task(db: AsyncSession = Depends(get_db)):
    new_task = models.ParsingTask(status="PENDING")
    db.add(new_task)
    await db.commit()
    await db.refresh(new_task)
    
    # Запуск Celery задачи (синхронный вызов — допустим, если Celery вне async loop)
    run_parsing_task.delay(new_task.id)
    
    return RedirectResponse(url="/", status_code=303)


@app.get("/download/{task_id}")
async def download_result(task_id: int, db: AsyncSession = Depends(get_db)):
    task = await db.get(models.ParsingTask, task_id)
    
    if not task or not task.result_path:
        raise HTTPException(status_code=404, detail="Файл не найден или задача еще не завершена")

    bucket_name = settings.AWS_STORAGE_BUCKET_NAME
    # Убираем префикс bucket_name, если он есть
    if task.result_path.startswith(f"{bucket_name}/"):
        s3_key = task.result_path[len(f"{bucket_name}/"):]
    else:
        s3_key = task.result_path

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