import requests
import sqlalchemy
from fastapi.encoders import jsonable_encoder
from sqlalchemy import select
from sqlalchemy.exc import SQLAlchemyError

from app.database import async_session_maker

class BaseDao:

    model = None

    @classmethod
    async def find_all_in_db(cls, **filter_by):
        async with async_session_maker() as session:
            query = select(cls.model).filter_by(**filter_by)
            result = await session.execute(query)
            return result.scalars().all()

    @classmethod
    async def find_one_by_id_in_db(cls, data_id):
        async with async_session_maker() as session:
            query = select(cls.model).filter_by(id=data_id)
            result = await session.execute(query)
            return result.scalar_one_or_none()

    @classmethod
    async def add_to_db(cls, **values):
        async with async_session_maker() as session:
            async with session.begin():
                new_instance = cls.model(**values)
                session.add(new_instance)
                try:
                    await session.commit()
                except SQLAlchemyError as e:
                    await session.rollback()
                    raise e

                return jsonable_encoder(new_instance)

    @classmethod
    async def delete_by_id(cls, data_id):
        async with async_session_maker() as session:
            async with session.begin():
                query = sqlalchemy.delete(cls.model).filter_by(id=data_id)
                result = await session.execute(query)
                try:
                    await session.commit()
                except SQLAlchemyError as e:
                    await session.rollback()
                    raise e
                return result.rowcount


def recored_and_upload_file(file_name):

    # URL вашего API
    url = "http://web:8000/api/files/upload/"

    # Файл, который нужно загрузить
    file_path = file_name

    # Заголовки
    headers = {
        "accept": "application/json",
    }

    # Открываем файл и отправляем его как multipart/form-data
    with open(file_path, "rb") as file:
        files = {
            "file_bytes": (file_path, file)  # Имя файла, файл, MIME-тип
        }
        response = requests.post(url, headers=headers, files=files)

    # Вывод результата
    print(response.status_code)
    print(response.json())

    return response.status_code