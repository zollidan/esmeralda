import os
from dotenv import load_dotenv
from pydantic_settings import BaseSettings


load_dotenv(dotenv_path=".env", override=True)


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
    ACCESS_TOKEN_EXPIRE_MINUTES: int = os.environ.get(
        "ACCESS_TOKEN_EXPIRE_MINUTES")
    POSTGRES_HOST: str = os.environ.get("POSTGRES_HOST")
    POSTGRES_PORT: str = os.environ.get("POSTGRES_PORT")
    POSTGRES_USER: str = os.environ.get("POSTGRES_USER")
    POSTGRES_PASSWORD: str = os.environ.get("POSTGRES_PASSWORD")
    POSTGRES_DB: str = os.environ.get("POSTGRES_DB")


settings = Settings()
