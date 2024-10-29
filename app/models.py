from sqlalchemy import Boolean, Column, String
from sqlalchemy.orm import DeclarativeBase
import cuid 

class Base(DeclarativeBase):
    pass

class User(Base):
    __tablename__ = 'User'
    
    id = Column(String(25), primary_key=True, default=cuid.cuid)
    name = Column(String(30))
    email = Column(String, unique=True, index=True)
    hashed_password = Column(String)
    is_active = Column(Boolean, default=True)

    def __repr__(self):
        return f"User(id={self.id}, email={self.email})"