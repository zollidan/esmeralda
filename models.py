from sqlalchemy import Column, Integer, String, DateTime, Text
from datetime import datetime
from database import Base

class ParsingTask(Base):
    __tablename__ = "parsing_tasks"

    id = Column(Integer, primary_key=True, index=True)
    status = Column(String, default="PENDING")  # PENDING, PROCESSING, DONE, ERROR
    result_path = Column(String, nullable=True)
    error_log = Column(Text, nullable=True)
    created_at = Column(DateTime, default=datetime.utcnow)