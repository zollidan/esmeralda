from fastapi import APIRouter, Depends, HTTPException, status, APIRouter, Response
from pydantic import BaseModel
from sqlalchemy.orm import Session
import requests
from app.models import User
from app.database import get_db
from app.parser.soccerway.index import soccerway
from . import schemas, models

router = APIRouter()

@router.get('/parser/soccerway')
def run_soccerway(date: str):
    
    result = soccerway(date)
 
    return result

@router.get('/parser/soccerway/test')
def run_soccerway_test():
    
    url = "https://ru.soccerway.com/a/block_h2h_matches?block_id=page_match_1_block_h2hsection_head2head_7_block_h2h_matches_1&action=changePage&callback_params=%7B%22page%22%3A+-1%2C+%22block_service_id%22%3A+%22match_h2h_comparison_block_h2hmatches%22%2C+%22team_A_id%22%3A+43000%2C+%22team_B_id%22%3A+43009%7D&params=%7B%22page%22%3A+0%7D"

    response = requests.get(url, headers={'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:102.0) Gecko/20100101 Firefox/131.0'})

    print(response.status_code)
 
    return response
