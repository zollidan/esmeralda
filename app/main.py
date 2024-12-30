import logging
from contextlib import asynccontextmanager

from app.config import settings
from fastapi.middleware.cors import CORSMiddleware
from fastapi import FastAPI
from art import tprint
from app.api.router import router


logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

app = FastAPI()

app.include_router(router)

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

