import os
from pydantic import Field
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    redis_host: str = Field(default='localhost')
    redis_port: int = Field(default=6379)
    redis_db: int = Field(default=0)
    
    # Настройки для Arq
    arq_queue_name: str = Field(default='arq:queue')
    arq_max_jobs: int = Field(default=10)
    arq_job_timeout: int = Field(default=3600)  # 1 час
    
    model_config = SettingsConfigDict(
        env_file=os.path.join(os.path.dirname(
            os.path.abspath(__file__)), ".env"),
        env_file_encoding='utf-8',
        case_sensitive=True
    )


settings = Settings()