from contextlib import asynccontextmanager
import logging
from typing import AsyncGenerator, Dict
from fastapi import FastAPI
import redis.asyncio as redis
import json
import os
from pydantic import Field
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    redis_host: str = Field(default='localhost')
    redis_port: int = Field(default=6379)
    redis_db: int = Field(default=0)
    
    model_config = SettingsConfigDict(
        env_file=os.path.join(os.path.dirname(
            os.path.abspath(__file__)), ".env"),
        env_file_encoding='utf-8',
        case_sensitive=True
    )


class RedisClient:
    def __init__(self, settings: Settings) -> None:
        self.settings = settings
        self.client: redis.Redis | None = None
    
    async def connect(self) -> None:
        self.client = await redis.Redis(host=self.settings.redis_host, port=self.settings.redis_port, db=self.settings.redis_db)
        logging.info(f"Connected to Redis at {self.settings.redis_host}:{self.settings.redis_port}")
    
    async def disconnect(self) -> None:
        if self.client:
            await self.client.close()
            logging.info("Redis connection closed")
    
    async def save_parsers(self, parsers: Dict[str, str]) -> None:
        if not self.client:
            raise RuntimeError("Redis client is not connected")
        
        # Save as JSON string
        await self.client.set("available_parsers", json.dumps(parsers))
        
        logging.info(f"Loaded {len(parsers)} parsers to Redis")
    
    async def get_parsers(self) -> Dict[str, str]:
        if not self.client:
            raise RuntimeError("Redis client is not connected")
        
        data = await self.client.get("available_parsers")
        if data:
            return json.loads(data)
        return {}
    
    async def get_parser_names(self) -> set[str]:
        if not self.client:
            raise RuntimeError("Redis client is not connected")
        
        return await self.client.smembers("parser_names")


settings = Settings()

# Parsers dictionary
parsers: Dict[str, str] = {
    "soccerway": "SoccerwayParser",
    "sofacore": "SofascoreParser",
    "flashscore": "FlashscoreParser",
    "livescore": "LivescoreParser",
    "testParser": "TestParser",
}


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncGenerator[None, None]:
    # Startup
    redis_client = RedisClient(settings)
    await redis_client.connect()
    await redis_client.save_parsers(parsers)
    app.state.redis = redis_client
    
    yield
    
    # Shutdown
    await redis_client.disconnect()


app = FastAPI(lifespan=lifespan)

@app.get("/parsers")
async def get_parsers() -> Dict[str, str]:
    """Получить список всех парсеров"""
    return await app.state.redis.get_parsers()
