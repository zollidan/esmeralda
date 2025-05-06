from celery import Celery

import os
import sys
sys.path.append(os.path.dirname(os.path.abspath(__file__)))


celery_app = Celery(
    "tasks",
)

celery_app.conf.broker_url = os.environ.get("CELERY_BROKER_URL", "redis://redis:6379")
celery_app.conf.result_backend = os.environ.get("CELERY_RESULT_BACKEND", "redis://redis:6379")


celery_app.conf.update(
    task_serializer='json',
    result_serializer='json',
    accept_content=['json'],
)

@celery_app.task
def run_soccerway_1(date_start, date_end):
    from parsers.soccerway_1.parser import main
    main(date_start, date_end)

@celery_app.task
def run_soccerway_2(date_start, date_end):
    from parsers.soccerway_2.parser import main
    main(date_start, date_end)