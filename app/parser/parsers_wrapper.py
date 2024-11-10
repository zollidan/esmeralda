import json
from fastapi import Depends
from art import tprint
import pandas as pd
import os
from sqlalchemy.orm import Session
import requests

from app import models
from app.database import get_db
from app.parser.server_parser_funcions import send_file_record_s3
from app.parser.soccerway.parser import run_soccerway
from app.repository import create_file_record
import uuid

def soccerway_wrapper(time_date: str):
     
    
    file_name = f'soccerway-{time_date}-{str(uuid.uuid4())}.xlsx'
    
    run_soccerway(user_date=time_date, my_file_name=file_name)
    # run_soccerway_test(user_date=time_date, file_name=file_name)

    send_file_record_s3(file_name)

    os.remove(file_name)

    return f'soccerway succsess {time_date}' 

def run_soccerway_test(user_date: str, file_name: str):
    
    print(user_date)
    
    df = pd.DataFrame({'a': [1, 2, 3], 'b': [4, 5, 6]})
    df.to_excel(file_name, index=False)
    
    
    print(f'run_soccerway_test {user_date}')