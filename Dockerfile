FROM python:3.10-slim

WORKDIR /code


ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

# Установка системных зависимостей, wine и поддержка i386
# RUN dpkg --add-architecture i386 && apt-get update && apt-get install -y --no-install-recommends \
#     wine \
#     wine64 \
#     wine32:i386 \
#     libwine:i386 \
#     libwine \
#     && apt-get clean \
#     && rm -rf /var/lib/apt/lists/*

COPY requirements.txt .
RUN pip install -r requirements.txt

COPY . . 


# CMD ["uvicorn", "app.main:app", "--proxy-headers", "--host", "0.0.0.0", "--port", "8000"]