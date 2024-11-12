import asyncio
import json
import time
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

async def soccerway_wrapper(time_date: str):
     
    
    file_name = f'soccerway-{time_date}-{str(uuid.uuid4())}.xlsx'
    
    run_soccerway(user_date=time_date, my_file_name=file_name)
    #run_soccerway_test(user_date=time_date, file_name=file_name)

    send_file_record_s3(file_name)

    os.remove(file_name)
    
    # for i in range(30):
    #     print(i)
    #     await asyncio.sleep(1) 
    
    
    """ тут отправляется сообещние в телеграм"""
    
    
    