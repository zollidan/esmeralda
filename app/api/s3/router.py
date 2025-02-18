from fastapi import APIRouter, HTTPException
from fastapi import Response
from app.config import s3_client, settings

s3_router = APIRouter(prefix='/api/s3', tags=['s3'])

@s3_router.get('/files')
def get_all_s3_files():
    return s3_client.list_objects(Bucket=settings.AWS_BUCKET)['Contents']


@s3_router.get("/files/{file_id}")
def get_file_by_id(file_id: str):
    try:
        get_object_response = s3_client.get_object(Bucket=settings.AWS_BUCKET, Key=file_id)
        file_binary = get_object_response['Body'].read()

        content_type = get_object_response['ContentType']

        return Response(
            content=file_binary,
            media_type=content_type
        )
    except Exception as e: 
        raise HTTPException(status_code=500, detail=str(e))