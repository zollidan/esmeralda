import os
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    POSTGRES_USER: str = "postgres"
    POSTGRES_HOST: str = "localhost"
    POSTGRES_PORT: str = "5432"
    POSTGRES_PASSWORD: str = "password"
    POSTGRES_DB: str = "esmeralda"
    AWS_ACCESS_KEY_ID: str = ""
    AWS_SECRET_ACCESS_KEY: str = ""
    AWS_BUCKET_NAME: str = "esmeralda-bucket"
    AWS_REGION: str = "ru-central1"
    model_config = SettingsConfigDict(env_file=os.path.join(os.path.dirname(os.path.dirname(__file__)), ".env"), env_file_encoding='utf-8')
    
    
settings = Settings()