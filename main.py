import datetime
import logging
from contextlib import asynccontextmanager
import os
import re
from aiogram import Bot
from fastapi import status
from typing import Annotated
from uuid import UUID, uuid4
from dotenv import load_dotenv
from fastapi import Depends, Response, UploadFile, File as FastAPIFile
from fastapi import FastAPI, HTTPException
from fastapi.responses import StreamingResponse
from fastapi.security import OAuth2PasswordBearer
from minio import Minio, S3Error
from pydantic_settings import BaseSettings
from sqlmodel import Field, Session, SQLModel, create_engine, select
from tasks import run_soccerway_1

load_dotenv(dotenv_path=".env.dev")

class Settings(BaseSettings):
    MINIO_ENDPOINT: str = os.environ.get("MINIO_ENDPOINT")
    MINIO_ROOT_USER: str = os.environ.get("MINIO_ROOT_USER")
    MINIO_ROOT_PASSWORD: str = os.environ.get("MINIO_ROOT_PASSWORD")
    BUCKET_NAME: str = os.environ.get("BUCKET_NAME")
    CELERY_BROKER_URL: str = os.environ.get("CELERY_BROKER_URL")
    CELERY_RESULT_BACKEND: str = os.environ.get("CELERY_RESULT_BACKEND")
    APP_MODE: str = os.environ.get("APP_MODE")
    TELEGRAM_BOT_TOKEN: str = os.environ.get("TELEGRAM_BOT_TOKEN")
    TELEGRAM_CHAT_ID: str = os.environ.get("TELEGRAM_CHAT_ID")
    TELEGRAM_TOPIC_ID: int = os.environ.get("TELEGRAM_TOPIC_ID")

settings = Settings()

s3_client = Minio(
    str(settings.MINIO_ENDPOINT),
    access_key=os.environ.get("MINIO_ROOT_USER"),
    secret_key=os.environ.get("MINIO_ROOT_PASSWORD"),
    secure=False,
)

bot = Bot(token=settings.TELEGRAM_BOT_TOKEN)

sqlite_file_name = "database.db"
sqlite_url = f"sqlite:///{sqlite_file_name}"

connect_args = {"check_same_thread": False}
engine = create_engine(sqlite_url, echo=True, connect_args=connect_args)

def create_db_and_tables():
    SQLModel.metadata.create_all(engine)

class User(SQLModel, table=True):
    id: UUID = Field(default_factory=uuid4, primary_key=True, index=True)
    username: str = Field(index=True)
    email: str = Field(index=True)
    full_name: str = Field(index=True)
    disabled: bool = Field(index=True)

class File(SQLModel, table=True):
    id: UUID = Field(default_factory=uuid4, primary_key=True, index=True)
    name: str = Field(index=True)
    file_url: str = Field(index=True)
    created_at: datetime.datetime = Field(default=datetime.datetime.now(), nullable=False)

@asynccontextmanager
async def lifespan(app: FastAPI):
    create_db_and_tables()
    yield

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
app = FastAPI(lifespan=lifespan)

oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")

@app.get("/items/")
async def read_items(token: Annotated[str, Depends(oauth2_scheme)]):
    return {"token": token}

@app.post("/api/files/upload")
def create_file(file: UploadFile = FastAPIFile(...)):
    file_id = uuid4()
    file_extension = os.path.splitext(file.filename)[1]
    file_name = f"{file_id}{file_extension}"
    try:
        s3_client.put_object(
            settings.BUCKET_NAME,
            file_name,
            file.file,
            length=-1,
            part_size=10 * 1024 * 1024, 
            content_type=file.content_type
        )

        file_url = f"api/files/{file_name}"
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
        stat = s3_client.stat_object(settings.BUCKET_NAME, id)
        content_type = stat.content_type 
        
        obj = s3_client.get_object(settings.BUCKET_NAME, id)

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

def is_valid_date(date_str):
    return bool(re.match(r'^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$', date_str))

@app.post("/api/run/soccerway1")
def run_soccerway(start_date: str, end_date: str):
    
    if not (is_valid_date(start_date) and is_valid_date(end_date)):
            return Response(status_code=status.HTTP_400_BAD_REQUEST, content="Invalid date format. Use YYYY-MM-DD.")
        
    try:
        task = run_soccerway_1.delay()
    except Exception as e:
        logging.error(f"Error starting task: {e}")
        raise HTTPException(status_code=500, detail="Error starting task")
        
    return {
        "message": "started",
        "task_id": task.id,
        "status": 200
    }
    
@app.post("/api/run/soccerway2")
def run_soccerway(start_date: str, end_date: str):

    if not (is_valid_date(start_date) and is_valid_date(end_date)):
        return Response(status_code=status.HTTP_400_BAD_REQUEST, content="Invalid date format. Use YYYY-MM-DD.")
    try:
        print(start_date, end_date)
    except Exception as e:
        logging.error(f"Error starting task: {e}")
        raise HTTPException(status_code=500, detail="Error starting task")
        
    return {
        "message": "started",
        "task_id": "task.id",
    }
    
# @app.post("/api/run/marafon")
# def run_soccerway():

#     try:
#         task = run_marafon.delay()
#     except Exception as e:
#         logging.error(f"Error starting task: {e}")
#         raise HTTPException(status_code=500, detail="Error starting task")
        
#     return {
#         "message": "started",
#         "task_id": "task.id",
#     }

@app.get("/api/bot/send_report_message")
async def send_message(text: str):
    try:
        await bot.send_message(chat_id=settings.TELEGRAM_CHAT_ID, message_thread_id=settings.TELEGRAM_TOPIC_ID, text=text)
    except Exception as e:
        logging.error(f"Error sending message: {e}")
        raise HTTPException(status_code=500, detail="Error sending message")
        
    return {
        "message": "Message sent",
    }