import os
from celery import Celery

os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'esmeralda.settings')

app = Celery('esmeralda')
app.config_from_object('django.conf:settings', namespace='CELERY')
app.autodiscover_tasks()