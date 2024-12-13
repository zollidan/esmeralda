from sqlalchemy import String, BigInteger, Integer, Date, Time, ForeignKey, Enum
from sqlalchemy.orm import Mapped, mapped_column, relationship
from app.database import Base
import enum

class File(Base):
    __tablename__ = 'files'
    
    file_id:Mapped[int] = mapped_column(BigInteger, primary_key=True)
    
    name: Mapped[str] = mapped_column(String, nullable=False)
    
    url: Mapped[str] = mapped_column(String, nullable=False)
    
class TelegramUser(Base):
    __tablename__ = 'telegram_user'
    
    telegram_id:Mapped[int] = mapped_column(BigInteger, primary_key=True)
    
    first_name: Mapped[str] = mapped_column(String, nullable=False)
    
    username: Mapped[str] = mapped_column(String, nullable=False)
    
    