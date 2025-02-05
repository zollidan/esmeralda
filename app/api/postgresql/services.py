from app.dao import BaseDao
from app.models import File


class FileService(BaseDao):
    model = File
