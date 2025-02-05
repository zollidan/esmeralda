from fastapi import APIRouter

from app.schemas import FileAddSchema
from app.services import FileService

postgresql_router = APIRouter(prefix='/api/postgresql_router', tags=['postgresql_router'])

@postgresql_router.get("/all")
async def get_all_files():
    return await FileService.find_all_in_db()

@postgresql_router.get("/file/{data_id}")
async def get_file_by_id(data_id: int):
    return await FileService.find_one_by_id_in_db(data_id)


@postgresql_router.post("/file/add")
async def add_file(file: FileAddSchema) -> dict:
    check = await FileService.add_to_db(**file.model_dump())
    if check:
        return {"message": "File added!", "file": check}
    else:
        return {"message": "Error", "file": check}