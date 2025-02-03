import os
import boto3
from pydantic_settings import BaseSettings, SettingsConfigDict

class Settings(BaseSettings):
    AWS_ID: str
    AWS_SECRET_KEY: str
    AWS_BUCKET: str
    FLOWER_USER:str
    FLOWER_PASSWORD:str
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