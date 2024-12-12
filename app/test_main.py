from fastapi.testclient import TestClient
from .api.router import router


client = TestClient(router)

def test_index_route():
    response = client.get("/api/")
    assert response.status_code == 200
    assert response.json() == {"version": "0.0.1"}
    
def test_run_soccerway_test_connection():
    response = client.get("/api/parser/soccerway/connection-test")
    assert response.status_code == 200
    
    
def test_files_router():
    response = client.get("/api/files")
    assert response.status_code == 200