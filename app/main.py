import logging
from contextlib import asynccontextmanager

import telegram
from app.bot.create_bot import bot, dp, stop_bot, start_bot
from app.bot.handlers.user_router import tg_user_router
from app.config import settings
from fastapi.staticfiles import StaticFiles
from aiogram.types import Update
from fastapi import FastAPI, Request
from art import tprint
from app.api.router import router
from app.pages.router import pages_router


logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

@asynccontextmanager
async def lifespan(app: FastAPI):
    print("-"* 15, "LeFort ESMERALDA 2024", "-"* 15)
    tprint('ESMERALDA')
    app.include_router(router)
    app.include_router(pages_router)
    yield
    logging.info("Shutting down appication...")
    


app = FastAPI(lifespan=lifespan)

