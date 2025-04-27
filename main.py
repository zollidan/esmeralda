import logging
from contextlib import asynccontextmanager
import os
from uuid import UUID, uuid4
from fastapi import UploadFile, File as FastAPIFile
from fastapi import FastAPI, HTTPException
from fastapi.responses import StreamingResponse
from minio import Minio, S3Error
from pydantic_settings import BaseSettings
from sqlmodel import Field, Session, SQLModel, create_engine, select
from celery.result import AsyncResult
from tasks import run_soccerway_1, celery_app

BUCKET_NAME = os.environ.get("BUCKET_NAME")

# class Settings(BaseSettings):
#     pass

s3_client = Minio(
    os.environ.get("MINIO_ENDPOINT"),
    access_key=os.environ.get("MINIO_ROOT_USER"),
    secret_key=os.environ.get("MINIO_ROOT_PASSWORD"),
    secure=False,
)

sqlite_file_name = "database.db"
sqlite_url = f"sqlite:///{sqlite_file_name}"

connect_args = {"check_same_thread": False}
engine = create_engine(sqlite_url, echo=True, connect_args=connect_args)

def create_db_and_tables():
    SQLModel.metadata.create_all(engine)

class File(SQLModel, table=True):
    id: UUID = Field(default_factory=uuid4, primary_key=True, index=True)
    name: str = Field(index=True)
    file_url: str = Field(index=True)
    created_at: datetime = Field(default=datetime.utcnow(), nullable=False)

@asynccontextmanager
async def lifespan(app: FastAPI):
    create_db_and_tables()
    yield

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
app = FastAPI(lifespan=lifespan)

@app.get("/hello/world")
def index():
    return {"message": "","version": "0.0.1"}

@app.post("/api/files/upload")
async def create_file(file: UploadFile = FastAPIFile(...)):
    file_id = uuid4()  # Генерация уникального UUID
    file_extension = os.path.splitext(file.filename)[1]
    file_name = f"{file_id}{file_extension}"  # Уникальное имя файла

    try:
        # Сохранение файла в MinIO
        s3_client.put_object(
            BUCKET_NAME,
            file_name,
            file.file,
            length=-1,  # Используется для потоков
            part_size=10 * 1024 * 1024,  # Размер частей (10MB)
            content_type=file.content_type
        )

        # Сохранение информации о файле в базе данных
        file_url = f"{BUCKET_NAME}/{file_name}"
        new_file = File(id=file_id, name=file.filename, file_url=file_url)

        with Session(engine) as session:
            session.add(new_file)
            session.commit()
            session.refresh(new_file)

        return {"id": str(file_id), "name": file.filename, "url": file_url}

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error uploading file: {e}")

@app.get("/api/files/")
def read_files():
    with Session(engine) as session:
        files = session.exec(select(File)).all()
        return files

@app.get("/api/files/{id}")
def get_file(id: str):
    try:
        stat = s3_client.stat_object(BUCKET_NAME, id)
        content_type = stat.content_type 
        
        obj = s3_client.get_object(BUCKET_NAME, id)

        return StreamingResponse(obj, media_type=content_type)

    except S3Error as e:
        raise HTTPException(status_code=404, detail=f"File not found: {e}")

@app.get("/api/files/{id}/info")
def read_file_info(id: UUID):
    with Session(engine) as session:
        file = session.get(File, id)
        if not file:
            raise HTTPException(status_code=404, detail="File not found")
        return file

@app.delete("/api/files/{id}")
def delete_file(id: UUID):
    with Session(engine) as session:
        file = session.get(File, id)
        if not file:
            raise HTTPException(status_code=404, detail="File not found")
        session.delete(file)
        session.commit()
        return {"ok": True}

@app.post("/api/run/soccerway1")
def run_soccerway():

    #file_name = f'soccerway-{date}-{str(uuid4())}.xls'
    
    task = run_soccerway_1.delay()
    
    return {
        "message": "started",
        "task_id": task.id,
        "status": 200
    }

@app.get("/api/task/{task_id}")
async def get_task_status(task_id: str):
    result = AsyncResult(task_id, app=celery_app)
    return {
        "task_id": task_id,
        "status": result.status,
        "result": result.result if result.ready() else None
    }
