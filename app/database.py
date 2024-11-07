from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from app.models import Base, Files

sql_url = 'sqlite:///./my_database.db'

engine = create_engine(sql_url)

def create_db():
    Base.metadata.create_all(engine)


SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)



def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()