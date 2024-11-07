import json
from fastapi import Depends
from art import tprint
import pandas as pd
import os
from sqlalchemy.orm import Session
import requests

from app import models
from app.database import get_db
from app.repository import create_file_record

def soccerway(time_date: str):
    
    file_name = f'soccerway-{time_date}.xlsx'
    
    data = [[1,2,3], [4,5,6], [time_date]]
    
    df = pd.DataFrame(data=data)
    
    create_file_record_href = f'https://yandex.bucker.com/myfile/{file_name}'

    create_file_body = json.dumps({
        "url": create_file_record_href
    })
    
    create_file_record = requests.post(f"http://localhost:8000/files/add", data=create_file_body)
    
    df.to_excel(file_name, index=False)

    

    # os.remove(file_name)

    return f'soccerway worked with date {time_date}' 