import os
import boto3
from pydantic_settings import BaseSettings, SettingsConfigDict

class Settings(BaseSettings):
    AWS_ID: str
    AWS_SECRET_KEY: str
    AWS_BUCKET: str
    FLOWER_URL:str
    FLOWER_USER:str
    FLOWER_PASSWORD:str
    SUPABASE_USER: str
    SUPABASE_PASSWORD: str
    SUPABASE_HOST: str
    SUPABASE_PORT: str
    SUPABASE_DBNAME: str
    JWT_SECRET_KEY: str
    JWT_ALGORITHM:str
    JWT_USER:str
    JWT_PASSWORD:str
    JWT_ACCESS_TOKEN_EXPIRE_MINUTES:int
    model_config = SettingsConfigDict(
        env_file=os.path.join(os.path.dirname(os.path.abspath(__file__)), "..", ".env")
    )
settings = Settings()

s3_client = boto3.client(
    's3',
    endpoint_url='https://storage.yandexcloud.net',
    aws_access_key_id=settings.AWS_ID,
    aws_secret_access_key=settings.AWS_SECRET_KEY,
)

def get_db_url():
    return f"postgresql+asyncpg://{settings.SUPABASE_USER}:{settings.SUPABASE_PASSWORD}@{settings.SUPABASE_HOST}:{settings.SUPABASE_PORT}/{settings.SUPABASE_DBNAME}"


def get_auth_data():
    return {"secret_key": settings.JWT_SECRET_KEY, "algorithm": settings.JWT_ALGORITHM}