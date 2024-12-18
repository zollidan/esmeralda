from app import celery_app
from app.parser.testparser.parser_funcs import parser_loop

from celery import Celery


# Функция для выполнения через Celery
@celery_app.task
def parse_data_test(self):
    print("[Test Parser start working]")
    parser_loop()  # Здесь вызывается ваша функция парсера