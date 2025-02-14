# создаем модель таблицы студентов
from sqlalchemy.orm import Mapped

from app.database import Base, str_pk, str_uniq, str_null_true


class File(Base):
    id: Mapped[str_pk]
    name: Mapped[str_uniq]
    file_url: Mapped[str_uniq]

class User(Base):
    id: Mapped[str_pk]
    name: Mapped[str_null_true]
    login: Mapped[str_uniq]
    password: Mapped[str]