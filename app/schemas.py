from pydantic import BaseModel, Field


class FileAddSchema(BaseModel):
    name: str
    file_url: str