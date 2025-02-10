from fastapi import APIRouter, HTTPException, Response
from starlette import status

from app.api.users.auth import authenticate_user, create_access_token, decode_jwt
from app.schemas import UserAuthSchema

user_router = APIRouter(prefix='/api/users', tags=['users'])

@user_router.post("/login")
async def auth_user(response: Response,user_data: UserAuthSchema):
    user = authenticate_user(email=user_data.email, password=user_data.password)
    if user is False:
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED,
                            detail='Неверная почта или пароль')
    access_token = create_access_token({"sub": str(user_data.email)})
    # response.set_cookie(key="users_access_token", value=access_token, httponly=True)
    return {'access_token': access_token, 'refresh_token': None}
