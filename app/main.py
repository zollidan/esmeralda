from fastapi import FastAPI, Request
from contextlib import asynccontextmanager
import logging
from app.bot.handlers.admin_router import admin_router
from app.bot.handlers.user_router import user_router
from app.config import settings
from aiogram.types import Update
from app.bot.create_bot import bot, dp, start_bot, stop_bot
from art import tprint
from app.api.router import router
from app.pages.router import pages_router

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Код, выполняющийся при запуске приложения
    tprint("esmeralda")
    app.include_router(router)
    app.include_router(pages_router)
    dp.include_router(user_router)
    dp.include_router(admin_router)
    await start_bot()
    
    webhook_url = settings.get_webhook_url()  # Получаем URL вебхука
    await bot.set_webhook(
        url=webhook_url,
        allowed_updates=dp.resolve_used_update_types(),
        drop_pending_updates=True
    )
    logging.info(f"Webhook set to {webhook_url}")
    yield  # Приложение работает
    # Код, выполняющийся при завершении работы приложения
    await bot.delete_webhook()
    logging.info("Webhook removed")

    
# Инициализация FastAPI с методом жизненного цикла
app = FastAPI(lifespan=lifespan)


@app.post("/webhook")
async def webhook(request: Request) -> None:
    logging.info("Received webhook request")
    update = Update.model_validate(await request.json(), context={"bot": bot})
    await dp.feed_update(bot, update)
    logging.info("Update processed")


# origins = [
#     "https://aaf-bet.ru",
#     "https://www.aaf-bet.ru",
#     "http://localhost:3000",
# ]

# app.add_middleware(
#     CORSMiddleware,
#     allow_origins=origins,
#     allow_credentials=True,
#     allow_methods=["*"],
#     allow_headers=["*"],
# )