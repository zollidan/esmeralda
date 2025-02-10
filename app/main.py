import logging

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.api.connection_check.router import status_check_router
from app.api.parsers.router import parser_router
from app.api.postgresql.router import files_router
from app.api.s3.router import s3_router
from app.api.tasks.router import tasks_router
from app.api.users.router import user_router

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

app = FastAPI()

app.include_router(status_check_router)
app.include_router(parser_router)
app.include_router(s3_router)
app.include_router(tasks_router)
app.include_router(files_router)
app.include_router(user_router)

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

