from fastapi import APIRouter, UploadFile, HTTPException
from sqlalchemy.exc import IntegrityError, SQLAlchemyError

from app.schemas import FileAddSchema
from app.api.postgresql.services import FileService
from app.config import s3_client, settings

files_router = APIRouter(prefix='/api/files', tags=['files'])

@files_router.get("/all")
async def get_all_files():
    return await FileService.find_all_in_db()

@files_router.get("/{data_id}")
async def get_file_by_id(data_id: int):
    return

@files_router.delete("/delete/{data_id}")
async def delete_file_by_id(data_id: str):

    file = await FileService.find_one_by_id_in_db(data_id)

    try:
        s3_client.delete_object(Bucket=settings.AWS_BUCKET, Key=file.name)
        await FileService.delete_by_id(data_id)
        return {"status": "deleted"}

    except SQLAlchemyError as e:
        return HTTPException(status_code=500, detail=f"Cannot delete file: {str(e)}")


@files_router.post("/upload/")
async def create_upload_file(file_bytes: UploadFile) -> dict:
    file_sql = FileAddSchema(
        name=file_bytes.filename,
        file_url=f"https://api.aaf-bet.ru/api/s3/files/{file_bytes.filename}"
    )

    try:
        file_content = await file_bytes.read()

        s3_client.put_object(
            Body=file_content,
            Bucket=settings.AWS_BUCKET,
            Key=file_bytes.filename
        )

        file = await FileService.add_to_db(**file_sql.model_dump())

        return {"file": file}

    except IntegrityError as e:
        raise HTTPException(status_code=400, detail=f"File already exists: {str(e)}")
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"An error occurred: {str(e)}")
    finally:
        await file_bytes.close()