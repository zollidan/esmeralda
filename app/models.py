# создаем модель таблицы студентов
from sqlalchemy.orm import Mapped

from app.database import Base, int_pk, str_uniq


class File(Base):
    id: Mapped[int_pk]
    name: Mapped[str_uniq]
    file_url: Mapped[str_uniq]
