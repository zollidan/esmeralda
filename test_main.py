import pytest
from fastapi.testclient import TestClient
from uuid import uuid4
from main import app, create_db_and_tables, engine, File, User
from sqlmodel import Session

client = TestClient(app)

@pytest.fixture(autouse=True)
def setup_database():
    # Создаем таблицы перед каждым тестом
    create_db_and_tables()
    yield
    # Удаляем данные после каждого теста
    with Session(engine) as session:
        session.execute("DELETE FROM file")
        session.execute("DELETE FROM user")
        session.commit()

def test_hello_world():
    response = client.get("/hello/world")
    assert response.status_code == 200
    assert response.json() == {"message": "", "version": "0.0.1"}

def test_file_upload():
    file_content = b"test content"
    files = {"file": ("test.txt", file_content, "text/plain")}
    response = client.post("/api/files/upload", files=files)
    assert response.status_code == 200
    data = response.json()
    assert "id" in data
    assert data["name"] == "test.txt"
    assert "url" in data

def test_read_files():
    # Загружаем файл
    file_content = b"test content"
    files = {"file": ("test.txt", file_content, "text/plain")}
    client.post("/api/files/upload", files=files)

    # Проверяем список файлов
    response = client.get("/api/files/")
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 1
    assert data[0]["name"] == "test.txt"

def test_register_user():
    user_data = {"email": "test@example.com", "password": "password123"}
    response = client.post("/api/auth/register", json=user_data)
    assert response.status_code == 200
    data = response.json()
    assert data["message"] == "User registered successfully"
    assert data["user"]["email"] == "test@example.com"

def test_register_user_duplicate_email():
    user_data = {"email": "test@example.com", "password": "password123"}
    client.post("/api/auth/register", json=user_data)

    # Попытка зарегистрировать пользователя с тем же email
    response = client.post("/api/auth/register", json=user_data)
    assert response.status_code == 400
    assert response.json()["detail"] == "Email already registered"

def test_delete_file():
    # Загружаем файл
    file_content = b"test content"
    files = {"file": ("test.txt", file_content, "text/plain")}
    upload_response = client.post("/api/files/upload", files=files)
    file_id = upload_response.json()["id"]

    # Удаляем файл
    delete_response = client.delete(f"/api/files/{file_id}")
    assert delete_response.status_code == 200
    assert delete_response.json() == {"ok": True}

    # Проверяем, что файл удален
    get_response = client.get(f"/api/files/{file_id}")
    assert get_response.status_code == 404