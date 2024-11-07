from datetime import datetime
from typing import List
from pydantic import BaseModel

class FileBody(BaseModel):
    url: str
