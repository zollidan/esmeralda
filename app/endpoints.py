from fastapi import APIRouter, Depends, HTTPException, status, APIRouter, Response
from pydantic import BaseModel
from sqlalchemy.orm import Session

from app.models import User
from app.database import get_db
from . import schemas, models, utils

router = APIRouter()

@router.get('/users')
def get_users(db: Session = Depends(get_db), limit: int = 10, page: int = 1, search: str = ''):
    skip = (page - 1) * limit

    users = db.query(User).all()
    return {'status': 'success', 'results': len(users), 'notes': users}

@router.post('/users', status_code=status.HTTP_201_CREATED)
async def create_user(payload: schemas.UserBody, db: Session = Depends(get_db)):
 
    hashed_password = utils.Hasher.get_password_hash(payload.password)
    
    new_user = models.User(name = payload.name, email = payload.email, hashed_password=hashed_password)
    
    print(new_user)
    
    # db.add(new_user)
    # db.commit()
    # db.refresh(new_user)
    
    return {"status": "success", "user": "new_user"}