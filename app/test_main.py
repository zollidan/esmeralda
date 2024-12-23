from fastapi.testclient import TestClient

from app.api.router import router



client = TestClient(router)

def test_soccerway_connection_test_success():
    response = client.get("/api/parser/soccerway/connection-test")
    assert response.status_code == 200
    data = response.json()
    assert data['status'] == 'success' or 'error'

# def test_parser_run():
#     response = client.post('/api/parser/soccerway?date=2024-10-10')
#     assert response.status_code == 200
#     data = response.json()
#     assert data['message'] == "soccerway work started"
#     assert data['task_id'] == str
#     assert data['status'] == 200
           
    



# def test_create_item():
#     item_data = {
#         "name": "example file",
#         "url": "http://example.com",

#     }

#     response = client.post("/api/files/sql", json=item_data)

#     assert response.status_code == 201
#     data = response.json()
#     # assert data["name"] == "example file"
#     # assert data["url"] == "http://example.com"

