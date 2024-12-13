from sqlalchemy import String, BigInteger, Integer, Date, Time, ForeignKey, Enum
from sqlalchemy.orm import Mapped, mapped_column, relationship
from app.database import Base
import enum

class File(Base):
    __tablename__ = 'files'
    
    file_id: Mapped[int] = mapped_column(Integer, primary_key=True,
                                           autoincrement=True) 
    
    name: Mapped[str] = mapped_column(String, nullable=False)
    
    url: Mapped[str] = mapped_column(String, nullable=False)
    
class User(Base):
    __tablename__ = 'user'
    
    file_id:Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    
    login: Mapped[str] = mapped_column(String, nullable=False)
    
    password: Mapped[str] = mapped_column(String, nullable=False)
    
    