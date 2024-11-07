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
from app.repository import create_file_record
import uuid

def soccerway(time_date: str):
     
    
    file_name = f'soccerway-{time_date}-{str(uuid.uuid4())}.xlsx'
    
    run_soccerway(file_name)

    send_file_record_s3(file_name)

    os.remove(file_name)

    return f'soccerway worked with date {time_date}' 