from uuid import UUID
from fastapi.testclient import TestClient
from httpx import AsyncClient
from main import app, settings
import io
from unittest.mock import MagicMock, patch
client = TestClient(app)


def test_healthcheck():
    response = client.get("/healthcheck")
    assert response.status_code == 200
    assert response.json() == {"status": "ok"}


def test_get_items():
    response = client.get("/api/files")
    assert response.status_code == 200
    data = response.json()

    assert isinstance(data, list)
    if not data:
        assert data == []
    else:
        first = data[0]
        assert isinstance(first, dict)
        assert "name" in first and isinstance(first["name"], str)
        assert "file_url" in first and isinstance(first["file_url"], str)
        assert "created_at" in first and isinstance(first["created_at"], str)
        assert "id" in first and isinstance(first["id"], str)


def test_upload_file_success():
    file_content = b"Hello, World!"
    test_file = ("test.txt", io.BytesIO(file_content), "text/plain")

    with patch("main.s3_client") as mock_put:
        mock_put.return_value = None

        response = client.post(
            "/api/files/upload",
            files={"file": test_file},
        )

        assert response.status_code == 200
        data = response.json()

        assert "id" in data and UUID(data["id"])
        assert data["name"] == "test.txt"
        assert data["url"].startswith("api/files/")


@patch("main.s3_client")
def test_get_image_file(mock_s3_client):
    mock_stat = MagicMock()
    mock_stat.content_type = "image/png"
    mock_s3_client.stat_object.return_value = mock_stat

    png_bytes = b"\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR"
    mock_s3_client.get_object.return_value = iter([png_bytes])

    response = client.get("/api/files/test-image.png")

    assert response.status_code == 200
    assert response.content.startswith(b"\x89PNG\r\n\x1a\n")
    assert response.headers["content-type"] == "image/png"


@patch("main.s3_client")
def test_get_file_not_found(mock_s3_client):
    from minio.error import S3Error

    mock_s3_client.stat_object.side_effect = S3Error(
        code="NoSuchKey",
        message="Object not found",
        resource="test-file.txt",
        request_id="xyz",
        host_id="abc",
        bucket_name="my-bucket",
        object_name="test-file.txt",
        response=None,
    )

    response = client.get("/api/files/test-file.txt")
    assert response.status_code == 404
    assert response.json()['detail'].startswith("File not found")


@patch("main.run_soccerway_1.delay")
def test_run_soccerway1_valid(mock_delay):
    mock_task = mock_delay.return_value
    mock_task.id = "1234abcd"

    response = client.post(
        "/api/run/soccerway1", params={"date_start": "2024-08-09", "date_end": "2024-12-23"})

    assert response.status_code == 200
    assert response.json() == {
        "message": "started",
        "task_id": "1234abcd"
    }
    mock_delay.assert_called_once_with("2024-08-09", "2024-12-23")


@patch("main.run_soccerway_2.delay")
def test_run_soccerway2_valid(mock_delay):
    mock_task = mock_delay.return_value
    mock_task.id = "1234abcd"

    response = client.post(
        "/api/run/soccerway2", params={"date_start": "2024-08-09", "date_end": "2024-12-23"})

    assert response.status_code == 200
    assert response.json() == {
        "message": "started",
        "task_id": "1234abcd"
    }
    mock_delay.assert_called_once_with("2024-08-09", "2024-12-23")
