from fastapi import APIRouter, Depends, HTTPException, status, APIRouter, Response
from pydantic import BaseModel
from sqlalchemy.orm import Session

from app.models import User
from app.database import get_db
from app.parser.soccerway.index import soccerway
from . import schemas, models

router = APIRouter()

@router.get('/parser/soccerway')
def run_soccerway(date: str):
    
    result = soccerway(date)
 
    return result

