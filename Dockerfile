# Базовый образ
FROM python:3.10-slim

WORKDIR /code

COPY . . 

RUN pip install --no-cache-dir --upgrade -r requirements.txt

CMD ["uvicorn", "app.main:app", "--proxy-headers", "--host", "0.0.0.0", "--port", "80"]