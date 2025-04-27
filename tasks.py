from celery import Celery

import os
import sys
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

from parsers.soccerway_1.parser import main
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
def run_soccerway_1():
    
    main()
