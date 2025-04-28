import logging
from contextlib import asynccontextmanager
from datetime import datetime
import os
from uuid import UUID, uuid4
from fastapi import UploadFile, File as FastAPIFile
from fastapi import FastAPI, HTTPException
from fastapi.responses import StreamingResponse
from minio import Minio, S3Error
from pydantic import BaseModel
from pydantic_settings import BaseSettings
from sqlmodel import Field, Session, SQLModel, create_engine, select
from celery.result import AsyncResult
from tasks import run_soccerway_1, celery_app
from passlib.context import CryptContext

BUCKET_NAME = os.environ.get("BUCKET_NAME")

# class Settings(BaseSettings):
#     pass

s3_client = Minio(
    str(os.environ.get("MINIO_ENDPOINT")),
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
    created_at: datetime = Field(default=datetime.now(), nullable=False)
    
class User(SQLModel, table=True):
    id: UUID = Field(default_factory=uuid4, primary_key=True, index=True)
    email: str = Field(index=True, unique=True, nullable=False)
    hashed_password: str = Field(nullable=False)
    created_at: datetime = Field(default=datetime.now(), nullable=False)

class UserCreateSchema(BaseModel):
    email: str
    password: str 

@asynccontextmanager
async def lifespan(app: FastAPI):
    create_db_and_tables()
    yield

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
app = FastAPI(lifespan=lifespan)

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")

@app.get("/hello/world")
def index():
    return {"message": "","version": "0.0.1"}

@app.post("/api/files/upload")
async def create_file(file: UploadFile = FastAPIFile(...)):
    file_id = uuid4()
    file_extension = os.path.splitext(file.filename)[1]
    file_name = f"{file_id}{file_extension}"

    try:
        s3_client.put_object(
            BUCKET_NAME,
            file_name,
            file.file,
            length=-1,  
            part_size=10 * 1024 * 1024,  
            content_type=file.content_type
        )

        file_url = f"/api/files/{file_name}"
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
    
    task = run_soccerway_1.delay()
    
    return {
        "message": "started",
        "task_id": task.id,
        "status": 200
    }
    
def verify_password(plain_password, hashed_password):
    return pwd_context.verify(plain_password, hashed_password)

def get_password_hash(password):
    return pwd_context.hash(password)
    
@app.post("/api/auth/register")
async def register_user(user: UserCreateSchema):
    
    hashed_password = get_password_hash(user.password)
    
    with Session(engine) as session:        
        existing_email = session.exec(select(User).where(User.email == user.email)).first()
        if existing_email:
            raise HTTPException(status_code=400, detail="Email already registered")
        
        user = User(email=user.email, hashed_password=hashed_password)
        session.add(user)
        session.commit()
        session.refresh(user)
    
    return {
        "message": "User registered successfully",
        "user": user,
    }

@app.get("/api/task/{task_id}")
async def get_task_status(task_id: str):
    result = AsyncResult(task_id, app=celery_app)
    return {
        "task_id": task_id,
        "status": result.status,
        "result": result.result if result.ready() else None
    }
