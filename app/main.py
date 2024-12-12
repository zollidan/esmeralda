import logging
from contextlib import asynccontextmanager

from app.config import settings
from fastapi.staticfiles import StaticFiles
from fastapi.middleware.cors import CORSMiddleware
from aiogram.types import Update
from fastapi import FastAPI, Request
from art import tprint
from app.api.router import router
from app.pages.router import pages_router
from app.bot.create_bot import bot, dp


logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

@asynccontextmanager
async def lifespan(app: FastAPI):
    print("-"* 15, "LeFort ESMERALDA 2024", "-"* 15)
    tprint('ESMERALDA')

    # Код, выполняющийся при запуске приложения
    webhook_url = settings.get_webhook_url()  # Получаем URL вебхука
    await bot.set_webhook(
        url=webhook_url,
        allowed_updates=dp.resolve_used_update_types(),
        drop_pending_updates=True
    )
    logging.info(f"Webhook set to {webhook_url}")
    
    
    yield # Приложение работает
    # Код, выполняющийся при завершении работы приложения
    await bot.delete_webhook()
    logging.info("Webhook removed")
    
    logging.info("Shutting down appication...")
    


# Инициализация FastAPI с методом жизненного цикла
app = FastAPI(lifespan=lifespan)

app.include_router(router)
app.include_router(pages_router)

# Маршрут для обработки вебхуков
@app.post("/webhook")
async def webhook(request: Request) -> None:
    logging.info("Received webhook request")
    update = await request.json()  # Получаем данные из запроса
    # Обрабатываем обновление через диспетчер (dp) и передаем в бот
    await dp.feed_update(bot, update)
    logging.info("Update processed")

origins = [
    "https://aaf-bet.ru",
    "https://www.aaf-bet.ru",
    "http://localhost:3000",
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get('/')
def test():
    return {"msg": "test"}