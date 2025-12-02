Архитектура проекта
Поток данных будет выглядеть так:

User нажимает кнопку на Next.js.

Next.js шлет запрос в FastAPI.

FastAPI создает запись в Postgres (статус "PENDING") и кидает задачу в Redis.

Celery Worker видит задачу в Redis, запускает Selenium.

Selenium парсит, сохраняет Excel, грузит в S3.

Celery обновляет статус в Postgres на "COMPLETED" и сохраняет ссылку на S3.

User видит на фронте, что всё готово, и скачивает файл.
