from sqlalchemy import Boolean, Column, DateTime, String
from sqlalchemy.orm import DeclarativeBase
from sqlalchemy.sql import func
import cuid 

class Base(DeclarativeBase):
    pass

class Files(Base):
    __tablename__ = 'Files'
    
    id = Column(String(25), primary_key=True, default=cuid.cuid)
    url = Column(String(255))
    time_created = Column(DateTime(timezone=True), server_default=func.now())


    def __repr__(self):
        return f"Files(id={self.id}, email={self.url})"