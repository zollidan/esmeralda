from sqlalchemy.future import select
from sqlalchemy.exc import SQLAlchemyError
from sqlalchemy.orm import joinedload
from app.dao.base import BaseDAO
from app.models import TelegramUser, File
from app.database import async_session_maker

class TelegramUserDAO(BaseDAO):
    model = TelegramUser

class FileDAO(BaseDAO):
    model = File