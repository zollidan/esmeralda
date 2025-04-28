import pytest
from fastapi.testclient import TestClient
from uuid import uuid4
from main import app, create_db_and_tables, engine, File, User
from sqlmodel import Session

client = TestClient(app)

def test_hello_world():
    response = client.get("/hello/world")
    assert response.status_code == 200
    assert response.json() == {"message": "", "version": "0.0.1"}
