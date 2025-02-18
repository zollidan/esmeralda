from unittest.mock import AsyncMock, patch
from fastapi.testclient import TestClient
from httpx import AsyncClient, ASGITransport
from pydantic import BaseModel
import pytest
from app.main import app

# Определяем ожидаемую структуру данных с Pydantic
class FileSchema(BaseModel):
    name: str
    created_at: str
    updated_at: str
    file_url: str
    id: str


client = TestClient(app)

def test_run_soccerway_test_connection_success():
    response = client.get("/api/connection-test/soccerway")
    data = response.json()
    assert response.status_code == 200

    
def test_run_marafon_test_connection():
    response = client.get('/api/connection-test/marafon')
    assert response.status_code == 200
    assert response.json() == {"status": response.status_code}
    
@pytest.mark.asyncio
async def test_get_all_files():
    async with AsyncClient(
        transport=ASGITransport(app=app), base_url='http://test',
    ) as ac:
        
        response = await ac.get('/api/files/all')
    
    assert response.status_code == 200
    assert isinstance(response.json(), list)
    

# @pytest.mark.asyncio
# async def test_get_file_by_id():
#     async with AsyncClient(
#         transport=ASGITransport(app=app), base_url='http://test',
#     ) as ac:
        
#         response = await ac.get('/api/files/string')
    
#     assert response.status_code == 200
    
# @pytest.mark.asyncio
# async def test_create_upload_file():
#     async with AsyncClient(
#         transport=ASGITransport(app=app), base_url='http://test',
#     ) as ac:
        
#         response = await ac.get('/api/files/string')
    
#     assert response.status_code == 200
    
# @pytest.mark.asyncio
# async def delete_file_by_id():
#     async with AsyncClient(
#         transport=ASGITransport(app=app), base_url='http://test',
#     ) as ac:
        
#         response = await ac.get('/api/files/string')
    
#     assert response.status_code == 200


def test_get_all_s3_files():
    response = client.get('api/s3/files')
    assert response.status_code == 200
    assert isinstance(response.json(), list)

def test_get_file_by_id():
    response = client.get('api/s3/files/wrongkey')
    assert response.status_code == 500
    response = client.get('api/s3/files/marafon.xlsx')
    assert response.status_code == 200
    assert isinstance(response.content, bytes)
 
