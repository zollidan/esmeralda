FROM python:3.12-slim

# Install uv.
COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/


# set env variables
ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

# Install the application dependencies.
WORKDIR /app

ADD . /app

RUN uv lock && uv sync --frozen --no-cache

CMD ["uv", "run", "uvicorn", "main:app", "--port", "80", "--host", "0.0.0.0"]