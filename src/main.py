from typing import Union
from .celery import add
from fastapi import FastAPI

app = FastAPI()


@app.post("/task")
def create_task():
    
    # task create logic
    
    pass
    
    
