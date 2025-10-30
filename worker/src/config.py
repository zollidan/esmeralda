from pydantic import Field
from pydantic_settings import BaseSettings, SettingsConfigDict


class S3Config(BaseSettings):
    """S3 Storage Configuration"""

    model_config = SettingsConfigDict(
        env_file='.env',
        env_file_encoding='utf-8',
        case_sensitive=False,
        extra='ignore'
    )

    endpoint_url: str = Field(
        default='https://storage.yandexcloud.net',
        alias='S3URL',
        description='S3 endpoint URL'
    )
    region: str = Field(
        default='ru-central1',
        alias='S3_REGION',
        description='S3 region'
    )
    access_key: str = Field(
        ...,
        alias='S3_ACCESS_KEY',
        description='S3 access key'
    )
    secret_key: str = Field(
        ...,
        alias='S3_SECRET_KEY',
        description='S3 secret key'
    )
    bucket_name: str = Field(
        default='esmeralda',
        alias='BUCKET_NAME',
        description='S3 bucket name'
    )


class RabbitMQConfig(BaseSettings):
    """RabbitMQ Configuration"""

    model_config = SettingsConfigDict(
        env_file='.env',
        env_file_encoding='utf-8',
        case_sensitive=False,
        extra='ignore'
    )

    host: str = Field(
        default='localhost',
        alias='RABBITMQ_HOST',
        description='RabbitMQ host'
    )
    port: int = Field(
        default=5672,
        alias='RABBITMQ_PORT',
        description='RabbitMQ port'
    )
    username: str = Field(
        default='guest',
        alias='RABBITMQ_USER',
        description='RabbitMQ username'
    )
    password: str = Field(
        default='guest',
        alias='RABBITMQ_PASSWORD',
        description='RabbitMQ password'
    )
    tasks_queue: str = Field(
        default='parser_tasks',
        alias='RABBITMQ_TASKS_QUEUE',
        description='Queue name for incoming tasks'
    )
    results_queue: str = Field(
        default='parser_results',
        alias='RABBITMQ_RESULTS_QUEUE',
        description='Queue name for results'
    )

    @property
    def connection_url(self) -> str:
        """Build RabbitMQ connection URL"""
        return f"amqp://{self.username}:{self.password}@{self.host}:{self.port}/"


class AppConfig(BaseSettings):
    """Application Configuration"""

    model_config = SettingsConfigDict(
        env_file='.env',
        env_file_encoding='utf-8',
        case_sensitive=False,
        extra='ignore'
    )

    log_level: str = Field(
        default='INFO',
        alias='LOG_LEVEL',
        description='Logging level'
    )
    worker_name: str = Field(
        default='worker-1',
        alias='WORKER_NAME',
        description='Worker instance name'
    )
    prefetch_count: int = Field(
        default=1,
        alias='PREFETCH_COUNT',
        description='Number of tasks to prefetch from queue'
    )


class Settings(BaseSettings):
    """Main settings container"""

    model_config = SettingsConfigDict(
        env_file='.env',
        env_file_encoding='utf-8',
        case_sensitive=False,
        extra='ignore'
    )

    s3: S3Config = Field(default_factory=S3Config)
    rabbitmq: RabbitMQConfig = Field(default_factory=RabbitMQConfig)
    app: AppConfig = Field(default_factory=AppConfig)


# Singleton instance
_settings: Settings | None = None


def get_settings() -> Settings:
    """Get settings singleton instance"""
    global _settings
    if _settings is None:
        _settings = Settings()
    return _settings
