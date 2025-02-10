from pydantic import BaseModel, Field


class FileAddSchema(BaseModel):
    name: str
    file_url: str

class UserAuthSchema(BaseModel):
    email:str = Field(..., description="Электронная почта")
    password: str = Field(..., min_length=5, max_length=50, description="Пароль, от 5 до 50 знаков")
