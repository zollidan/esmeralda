from datetime import datetime
from typing import List
from pydantic import BaseModel

class UserBody(BaseModel):
    name: str
    email: str
    password: str
