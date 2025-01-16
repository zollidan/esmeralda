from fastapi.testclient import TestClient

from app.api.parsers.router import router



client = TestClient(router)

