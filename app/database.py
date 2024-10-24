from sqlalchemy import create_engine

from app.models import Base

sql_url = 'sqlite:///./my_database.db'

engine = create_engine(sql_url, echo=True)

def create_db():
    Base.metadata.create_all(engine)