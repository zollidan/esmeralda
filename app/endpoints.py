from fastapi import APIRouter

router = APIRouter()

@router.get('/users')
async def read_users():
    return {"users"}

@router.post('/users')
async def create_user():
    return {"create user POST request"}