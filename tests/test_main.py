from fastapi.testclient import TestClient

from app.api.connection_check.router import status_check_router

# Создаем тестовый клиент
client = TestClient(status_check_router)

# Тест для проверки соединения с Soccerway
def test_run_soccerway_test_connection():
    response = client.get("/api/connection-test/soccerway")
    assert response.status_code == 200
    assert "status" in response.json()
    assert response.json()["status"] == 200 or response.json()["status"] == "error"

# Тест для проверки соединения с Marafon
def test_run_marafon_test_connection():
    response = client.get("/api/connection-test/marafon")
    assert response.status_code == 200
    assert "status" in response.json()
    assert response.json()["status"] == 200 or response.json()["status"] == "error"

# Дополнительные тесты для проверки ошибок
def test_run_soccerway_test_connection_error():
    # Здесь можно имитировать ошибку, например, изменив URL или заголовки
    response = client.get("/api/connection-test/soccerway")
    assert response.status_code == 200
    assert "status" in response.json()
    if response.json()["status"] == "error":
        assert "status_code" in response.json()
        assert "message" in response.json()

def test_run_marafon_test_connection_error():
    # Здесь можно имитировать ошибку, например, изменив URL или заголовки
    response = client.get("/api/connection-test/marafon")
    assert response.status_code == 200
    assert "status" in response.json()
    if response.json()["status"] == "error":
        assert "status_code" in response.json()
        assert "message" in response.json()