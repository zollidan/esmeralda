import os
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    """Application settings with environment variable support"""

    # RabbitMQ Configuration
    RABBITMQ_HOST: str = "localhost"
    RABBITMQ_PORT: int = 5672
    RABBITMQ_USER: str = "guest"
    RABBITMQ_PASSWORD: str = "guest"
    RABBITMQ_VHOST: str = "/"
    RABBITMQ_TASKS_QUEUE: str = "task_queue"
    RABBITMQ_RESULTS_QUEUE: str = "result_queue"

    # Worker Configuration
    WORKER_PREFETCH_COUNT: int = 1
    WORKER_HEARTBEAT: int = 600
    WORKER_BLOCKED_CONNECTION_TIMEOUT: int = 300

    # Storage Configuration (S3)
    S3_BUCKET: str = ""
    S3_REGION: str = "us-east-1"
    AWS_ACCESS_KEY_ID: str = ""
    AWS_SECRET_ACCESS_KEY: str = ""

    # Logging Configuration
    LOG_LEVEL: str = "INFO"
    LOG_FILE: str = "worker.log"

    model_config = SettingsConfigDict(
        env_file=os.path.join(os.path.dirname(
            os.path.abspath(__file__)), "..", ".env"),
        env_file_encoding='utf-8',
        case_sensitive=True
    )


settings: Settings = Settings()
