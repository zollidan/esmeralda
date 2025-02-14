from fastapi.testclient import TestClient

from app.api.connection_check.router import status_check_router

# Создаем тестовый клиент
client = TestClient(status_check_router)

# Тест для проверки соединения с Soccerway
def test_run_soccerway_test_connection():
    response = client.get("/api/connection-test/soccerway")
    assert response.status_code == 200
    assert "status" in response.json()
    # assert response.json()["status"] == int

# Тест для проверки соединения с Marafon
def test_run_marafon_test_connection():
    response = client.get("/api/connection-test/marafon")
    assert response.status_code == 200
    assert "status" in response.json()
    # assert response.json()["status"] == int
