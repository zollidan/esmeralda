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


logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

@asynccontextmanager
async def lifespan(app: FastAPI):
    print("-"* 15, "LeFort ESMERALDA 2024", "-"* 15)
    tprint('ESMERALDA')

    app.include_router(router)
    app.include_router(pages_router)
    
    
    yield 

    
    logging.info("Shutting down appication...")
    


# Инициализация FastAPI с методом жизненного цикла
app = FastAPI(lifespan=lifespan)




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

