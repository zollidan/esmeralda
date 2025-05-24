from datetime import datetime, timedelta, timezone
import logging
from contextlib import asynccontextmanager
import os
import re
from aiogram import Bot
from fastapi import Request, status
from typing import Annotated
from uuid import UUID, uuid4
from dotenv import load_dotenv
from fastapi import Depends, Response, UploadFile, File as FastAPIFile
from fastapi import FastAPI, HTTPException
from fastapi.responses import HTMLResponse, StreamingResponse
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from fastapi.middleware.cors import CORSMiddleware
from fastapi.staticfiles import StaticFiles
from fastapi.templating import Jinja2Templates
from starlette.middleware.base import BaseHTTPMiddleware
import jwt
from minio import Minio, S3Error
from pydantic_settings import BaseSettings
from passlib.context import CryptContext
from sqlmodel import Field, Session, SQLModel, create_engine, select
from tasks import run_soccerway_1, run_soccerway_2

load_dotenv(dotenv_path=".env")

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
    SECRET_KEY: str = os.environ.get("SECRET_KEY")
    ALGORITHM: str = os.environ.get("ALGORITHM")
    ACCESS_TOKEN_EXPIRE_MINUTES: int = os.environ.get("ACCESS_TOKEN_EXPIRE_MINUTES")
    POSTGRES_HOST: str = os.environ.get("POSTGRES_HOST")
    POSTGRES_PORT: str = os.environ.get("POSTGRES_PORT")
    POSTGRES_USER: str = os.environ.get("POSTGRES_USER")
    POSTGRES_PASSWORD: str = os.environ.get("POSTGRES_PASSWORD")
    POSTGRES_DB: str = os.environ.get("POSTGRES_DB")

settings = Settings()

s3_client = Minio(
    str(settings.MINIO_ENDPOINT),
    access_key=os.environ.get("MINIO_ROOT_USER"),
    secret_key=os.environ.get("MINIO_ROOT_PASSWORD"),
    secure=False,
)

bot = Bot(token=settings.TELEGRAM_BOT_TOKEN)

oauth2_scheme = OAuth2PasswordBearer(tokenUrl="auth/login")
pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")

postgres_url = f"postgresql://{settings.POSTGRES_USER}:{settings.POSTGRES_PASSWORD}@{settings.POSTGRES_HOST}:{settings.POSTGRES_PORT}/{settings.POSTGRES_DB}"
engine = create_engine(postgres_url, echo=False)

class User(SQLModel, table=True):
    id: UUID = Field(default_factory=uuid4, primary_key=True, index=True)
    username: str = Field(index=True)
    full_name: str = Field(index=True)
    disabled: bool = Field(index=True, default=False)
    hashed_password: str = Field(index=True)
    
class UserRegisterSchema(SQLModel):
    username: str
    full_name: str
    password: str
    
class UserLoginSchema(SQLModel):
    username: str
    password: str
    
class File(SQLModel, table=True):
    id: UUID = Field(default_factory=uuid4, primary_key=True, index=True)
    name: str = Field(index=True)
    file_url: str = Field(index=True)
    created_at: datetime = Field(default=datetime.now(), nullable=False)

def create_db_and_tables():
    SQLModel.metadata.create_all(engine)
    
def verify_password(plain_password, hashed_password):
    return pwd_context.verify(plain_password, hashed_password)

def get_password_hash(password):
    return pwd_context.hash(password)
    
def get_user(username: str):
    with Session(engine) as session:
        user = session.exec(select(User).where(User.username == username)).first()
        return user
    
def authenticate_user(username: str, password: str):
    user = get_user(username)
    if not user:
        return False
    if not verify_password(password, user.hashed_password):
        return False
    return user
    
def create_access_token(data: dict, expires_delta: timedelta | None = None):
    to_encode = data.copy()
    if expires_delta:
        expire = datetime.now(timezone.utc) + expires_delta
    else:
        expire = datetime.now(timezone.utc) + timedelta(minutes=15)
    to_encode.update({"exp": expire})
    encoded_jwt = jwt.encode(to_encode, settings.SECRET_KEY, algorithm=settings.ALGORITHM)
    return encoded_jwt

@asynccontextmanager
async def lifespan(app: FastAPI):
    create_db_and_tables()
    yield

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
app = FastAPI(lifespan=lifespan)

app.mount("/static", StaticFiles(directory="static"), name="static")
templates = Jinja2Templates(directory="templates")

# class JWTMiddleware(BaseHTTPMiddleware):
#     async def dispatch(self, request: Request, call_next):

#         # Skip authentication for public endpoints
#         if request.url.path in ["/auth/login", "/auth/registration", "/docs", "/openapi.json", "/api/files/upload", "/api/bot/send_report_message"]:
#             return await call_next(request)

#         # Extract token from Authorization header
#         auth_header = request.headers.get("Authorization")
#         if not auth_header or not auth_header.startswith("Bearer "):
#             return Response(
#                 status_code=status.HTTP_401_UNAUTHORIZED,
#                 content="Missing or invalid Authorization header",
#                 headers={"WWW-Authenticate": "Bearer"},
#             )

#         token = auth_header.split("Bearer ")[1]
#         try:
#             payload = jwt.decode(token, settings.SECRET_KEY, algorithms=[settings.ALGORITHM])
#             username = payload.get("sub")
#             if username is None:
#                 return Response(
#                     status_code=status.HTTP_401_UNAUTHORIZED,
#                     content="Could not validate credentials",
#                     headers={"WWW-Authenticate": "Bearer"},
#                 )
#         except jwt.InvalidTokenError:
#             return Response(
#                 status_code=status.HTTP_401_UNAUTHORIZED,
#                 content="Invalid token",
#                 headers={"WWW-Authenticate": "Bearer"},
#             )

#         user = get_user(username)
#         if user is None:
#             return Response(
#                 status_code=status.HTTP_401_UNAUTHORIZED,
#                 content="User not found",
#                 headers={"WWW-Authenticate": "Bearer"},
#             )

#         if user.disabled:
#             return Response(
#                 status_code=status.HTTP_400_BAD_REQUEST,
#                 content="Inactive user",
#             )

#         # Store user in request.state
#         request.state.user = user
#         response = await call_next(request)
#         return response

# app.add_middleware(JWTMiddleware)

origins = [
    "https://esmeralda-frontend.vercel.app"
    "http://localhost:3000",
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/healthcheck")
def healthcheck():
    return {"status": "ok"}

@app.get("/info")
def info():
    return {
        "version": "1.0.1",
        "description": "API for Esmeralda project",
        "author": "LeFort Tech Team @2025",
        "email": ""
    }

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
def run_soccerway1(date_start: str, date_end: str):
    
    if not (is_valid_date(date_start) and is_valid_date(date_end)):
            return Response(status_code=status.HTTP_400_BAD_REQUEST, content="Invalid date format. Use YYYY-MM-DD.")
        
    try:
        task = run_soccerway_1.delay(date_start, date_end)
    except Exception as e:
        logging.error(f"Error starting task: {e}")
        raise HTTPException(status_code=500, detail="Error starting task")
        
    return {
        "message": "started",
        "task_id": task.id,
    }
    
@app.post("/api/run/soccerway2")
def run_soccerway2(date_start: str, date_end: str):

    if not (is_valid_date(date_start) and is_valid_date(date_end)):
        return Response(status_code=status.HTTP_400_BAD_REQUEST, content="Invalid date format. Use YYYY-MM-DD.")
    try:
        task = run_soccerway_2.delay(date_start, date_end)
    except Exception as e:
        logging.error(f"Error starting task: {e}")
        raise HTTPException(status_code=500, detail="Error starting task")
        
    return {
        "message": "started",
        "task_id": task.id,
    }
    
# @app.post("/api/run/marafon")
# def run_marafon():

#     try:
#         task = run_marafon.delay()
#     except Exception as e:
#         logging.error(f"Error starting task: {e}")
#         raise HTTPException(status_code=500, detail="Error starting task")
        
#     return {
#         "message": "started",
#         "task_id": "task.id",
#     }

@app.post("/api/bot/send_report_message")
async def send_message(text: str):
    try:
        await bot.send_message(chat_id=settings.TELEGRAM_CHAT_ID, message_thread_id=settings.TELEGRAM_TOPIC_ID, text=text)
    except Exception as e:
        logging.error(f"Error sending message: {e}")
        raise HTTPException(status_code=500, detail="Error sending message")
        
    return {
        "message": "Message sent",
    }

# @app.post("/auth/login")
# async def login_for_access_token(
#     form_data: Annotated[OAuth2PasswordRequestForm, Depends()],
# ):
#     user = authenticate_user(form_data.username, form_data.password)
    
#     if not user:
#         raise HTTPException(
#             status_code=status.HTTP_401_UNAUTHORIZED,
#             detail="Incorrect username or password",
#             headers={"WWW-Authenticate": "Bearer"},
#         )

#     access_token = create_access_token(
#         data={"sub": user.username, "userId": str(user.id)}, expires_delta=timedelta(minutes=settings.ACCESS_TOKEN_EXPIRE_MINUTES)
#     )
#     return {"access_token": access_token, "token_type": "bearer"}

            
# @app.post("/auth/registration")
# async def registration(form_data: UserRegisterSchema):
#     with Session(engine) as session:
#         existing_user = session.exec(select(User).where(User.username == form_data.username)).first()
#         if existing_user:
#             raise HTTPException(status_code=400, detail="Username already exists")

#         hashed_password = pwd_context.hash(form_data.password)
#         user = User(username=form_data.username, full_name=form_data.full_name, hashed_password=hashed_password)
        
#         session.add(user)
#         session.commit()
#         session.refresh(user)
#         return {
#             "id": str(user.id),
#             "username": user.username,
#             "full_name": user.full_name,
#         }
        
# @app.get("/users/me/", response_model=User)
# async def read_users_me(request: Request):
#     return request.state.user

@app.get("/admin", response_class=HTMLResponse)
async def read_item(request: Request, id: str):
    return templates.TemplateResponse(
        request=request, name="index.html"
    )
