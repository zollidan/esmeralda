# API Documentation - Microservices Architecture

## Архитектура системы

```
User → nginx → front (UI)
              ↓
         go api (main backend)
              ↓
    ┌─────────┼─────────┐
    ↓         ↓         ↓
RabbitMQ   postgres   redis   s3 yandex
    ↓
worker → parser1, parser2
    ↓              ↓
  redis      s3 yandex
```

## Технологический стек

- **Frontend**: React/Vue/Angular (через nginx)
- **Backend API**: Go (Fiber/Gin/Chi)
- **Message Broker**: RabbitMQ
- **Database**: PostgreSQL + postgresUS (расширение)
- **Cache**: Redis
- **Storage**: S3 Yandex Cloud
- **Worker**: Piko (Python/Go)
- **Parsers**: parser1, parser2

---

## 1. Authentication & Authorization API

### Base URL: `/api/v1/auth`

#### 1.1 Регистрация пользователя

```http
POST /api/v1/auth/register
```

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!",
  "username": "username",
  "first_name": "Ivan",
  "last_name": "Ivanov"
}
```

**Response (201 Created):**

```json
{
  "success": true,
  "data": {
    "user_id": "uuid",
    "email": "user@example.com",
    "username": "username",
    "created_at": "2025-11-12T10:30:00Z"
  },
  "message": "User registered successfully"
}
```

---

#### 1.2 Вход (Login)

```http
POST /api/v1/auth/login
```

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": {
      "user_id": "uuid",
      "email": "user@example.com",
      "username": "username"
    }
  }
}
```

---

#### 1.3 Обновление токена

```http
POST /api/v1/auth/refresh
```

**Request Body:**

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

---

#### 1.4 Выход (Logout)

```http
POST /api/v1/auth/logout
Authorization: Bearer {access_token}
```

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

---

#### 1.5 Проверка токена

```http
GET /api/v1/auth/verify
Authorization: Bearer {access_token}
```

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "valid": true,
    "user_id": "uuid",
    "expires_at": "2025-11-12T11:30:00Z"
  }
}
```

---

#### 1.6 Восстановление пароля (запрос)

```http
POST /api/v1/auth/password/reset-request
```

**Request Body:**

```json
{
  "email": "user@example.com"
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Password reset email sent"
}
```

---

#### 1.7 Восстановление пароля (подтверждение)

```http
POST /api/v1/auth/password/reset-confirm
```

**Request Body:**

```json
{
  "token": "reset_token_from_email",
  "new_password": "NewSecurePassword123!"
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Password reset successfully"
}
```

---

## 2. User Management API

### Base URL: `/api/v1/users`

#### 2.1 Получить профиль текущего пользователя

```http
GET /api/v1/users/me
Authorization: Bearer {access_token}
```

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "user_id": "uuid",
    "email": "user@example.com",
    "username": "username",
    "first_name": "Ivan",
    "last_name": "Ivanov",
    "created_at": "2025-11-12T10:30:00Z",
    "updated_at": "2025-11-12T10:30:00Z"
  }
}
```

---

#### 2.2 Обновить профиль

```http
PUT /api/v1/users/me
Authorization: Bearer {access_token}
```

**Request Body:**

```json
{
  "first_name": "Ivan",
  "last_name": "Petrov",
  "username": "new_username"
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "user_id": "uuid",
    "email": "user@example.com",
    "username": "new_username",
    "first_name": "Ivan",
    "last_name": "Petrov",
    "updated_at": "2025-11-12T12:00:00Z"
  }
}
```

---

#### 2.3 Удалить аккаунт

```http
DELETE /api/v1/users/me
Authorization: Bearer {access_token}
```

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Account deleted successfully"
}
```

---

## 3. Parser Management API

### Base URL: `/api/v1/parsers`

#### 3.1 Получить список доступных парсеров

```http
GET /api/v1/parsers
Authorization: Bearer {access_token}
```

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "parsers": [
      {
        "parser_id": "parser1",
        "name": "Parser 1",
        "description": "Description of parser 1",
        "status": "active",
        "supported_sources": ["source1", "source2"]
      },
      {
        "parser_id": "parser2",
        "name": "Parser 2",
        "description": "Description of parser 2",
        "status": "active",
        "supported_sources": ["source3", "source4"]
      }
    ]
  }
}
```

**Примечание:** Список парсеров хранится в Redis и обновляется worker'ом при старте.

---

#### 3.2 Получить информацию о конкретном парсере

```http
GET /api/v1/parsers/{parser_id}
Authorization: Bearer {access_token}
```

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "parser_id": "parser1",
    "name": "Parser 1",
    "description": "Detailed description",
    "status": "active",
    "supported_sources": ["source1", "source2"],
    "configuration": {
      "timeout": 30,
      "retry_attempts": 3
    },
    "last_health_check": "2025-11-12T12:00:00Z"
  }
}
```

---

## 4. Parsing Tasks API

### Base URL: `/api/v1/tasks`

#### 4.1 Создать задачу на парсинг

```http
POST /api/v1/tasks
Authorization: Bearer {access_token}
```

**Request Body:**

```json
{
  "parser_id": "parser1",
  "source_url": "https://example.com/data",
  "parameters": {
    "depth": 2,
    "format": "json"
  },
  "priority": "normal"
}
```

**Response (201 Created):**

```json
{
  "success": true,
  "data": {
    "task_id": "uuid",
    "parser_id": "parser1",
    "status": "queued",
    "created_at": "2025-11-12T12:00:00Z",
    "estimated_completion": "2025-11-12T12:05:00Z"
  }
}
```

**Процесс:**

1. Go API создает запись в PostgreSQL
2. Отправляет сообщение в RabbitMQ
3. Сохраняет задачу в Redis для быстрого доступа
4. Worker получает задачу из RabbitMQ
5. Worker обновляет статус в Redis

---

#### 4.2 Получить список задач пользователя

```http
GET /api/v1/tasks
Authorization: Bearer {access_token}
```

**Query Parameters:**

- `status` (optional): `queued`, `processing`, `completed`, `failed`
- `parser_id` (optional): filter by parser
- `limit` (optional, default: 20): number of results
- `offset` (optional, default: 0): pagination offset
- `sort` (optional, default: `created_at`): sort field
- `order` (optional, default: `desc`): `asc` or `desc`

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "tasks": [
      {
        "task_id": "uuid",
        "parser_id": "parser1",
        "status": "completed",
        "created_at": "2025-11-12T12:00:00Z",
        "started_at": "2025-11-12T12:01:00Z",
        "completed_at": "2025-11-12T12:05:00Z",
        "result_url": "https://s3.yandexcloud.net/bucket/results/uuid.json"
      }
    ],
    "pagination": {
      "total": 100,
      "limit": 20,
      "offset": 0,
      "has_more": true
    }
  }
}
```

---

#### 4.3 Получить информацию о задаче

```http
GET /api/v1/tasks/{task_id}
Authorization: Bearer {access_token}
```

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "task_id": "uuid",
    "user_id": "uuid",
    "parser_id": "parser1",
    "source_url": "https://example.com/data",
    "status": "processing",
    "progress": 45,
    "created_at": "2025-11-12T12:00:00Z",
    "started_at": "2025-11-12T12:01:00Z",
    "estimated_completion": "2025-11-12T12:05:00Z",
    "parameters": {
      "depth": 2,
      "format": "json"
    },
    "logs": [
      {
        "timestamp": "2025-11-12T12:01:30Z",
        "level": "info",
        "message": "Started parsing"
      }
    ]
  }
}
```

**Примечание:** Статус и прогресс читаются из Redis для актуальности.

---

#### 4.4 Отменить задачу

```http
DELETE /api/v1/tasks/{task_id}
Authorization: Bearer {access_token}
```

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Task cancelled successfully"
}
```

---

#### 4.5 Получить результат задачи

```http
GET /api/v1/tasks/{task_id}/result
Authorization: Bearer {access_token}
```

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "task_id": "uuid",
    "status": "completed",
    "result_url": "https://s3.yandexcloud.net/bucket/results/uuid.json",
    "result_size": 1024000,
    "format": "json",
    "completed_at": "2025-11-12T12:05:00Z",
    "download_url": "/api/v1/tasks/{task_id}/download"
  }
}
```

---

#### 4.6 Скачать результат задачи

```http
GET /api/v1/tasks/{task_id}/download
Authorization: Bearer {access_token}
```

**Response (200 OK):**

- Content-Type: application/json | text/csv | application/xml
- Файл будет загружен из S3 Yandex Cloud

---

## 5. Files/Storage API

### Base URL: `/api/v1/files`

#### 5.1 Загрузить файл

```http
POST /api/v1/files/upload
Authorization: Bearer {access_token}
Content-Type: multipart/form-data
```

**Request Body:**

```
file: <binary data>
description: "Optional description"
```

**Response (201 Created):**

```json
{
  "success": true,
  "data": {
    "file_id": "uuid",
    "filename": "document.pdf",
    "size": 1024000,
    "url": "https://s3.yandexcloud.net/bucket/uploads/uuid/document.pdf",
    "uploaded_at": "2025-11-12T12:00:00Z"
  }
}
```

---

#### 5.2 Получить список файлов

```http
GET /api/v1/files
Authorization: Bearer {access_token}
```

**Query Parameters:**

- `limit` (optional, default: 20)
- `offset` (optional, default: 0)

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "files": [
      {
        "file_id": "uuid",
        "filename": "document.pdf",
        "size": 1024000,
        "url": "https://s3.yandexcloud.net/bucket/uploads/uuid/document.pdf",
        "uploaded_at": "2025-11-12T12:00:00Z"
      }
    ],
    "pagination": {
      "total": 50,
      "limit": 20,
      "offset": 0
    }
  }
}
```

---

#### 5.3 Удалить файл

```http
DELETE /api/v1/files/{file_id}
Authorization: Bearer {access_token}
```

**Response (200 OK):**

```json
{
  "success": true,
  "message": "File deleted successfully"
}
```

---

## 6. Health Check & Monitoring API

### Base URL: `/api/v1/health`

#### 6.1 Проверка здоровья API

```http
GET /api/v1/health
```

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "timestamp": "2025-11-12T12:00:00Z",
    "services": {
      "database": "healthy",
      "redis": "healthy",
      "rabbitmq": "healthy",
      "s3": "healthy"
    }
  }
}
```

---

#### 6.2 Статистика системы

```http
GET /api/v1/health/stats
Authorization: Bearer {access_token}
```

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "tasks": {
      "queued": 5,
      "processing": 12,
      "completed": 1543,
      "failed": 23
    },
    "parsers": {
      "active": 2,
      "inactive": 0
    },
    "workers": {
      "active": 3,
      "idle": 2
    }
  }
}
```

---

## 7. WebSocket API (Real-time Updates)

### Base URL: `ws://api.example.com/api/v1/ws`

#### 7.1 Подключение к WebSocket

```
ws://api.example.com/api/v1/ws?token={access_token}
```

#### 7.2 Подписка на обновления задачи

**Client → Server:**

```json
{
  "action": "subscribe",
  "resource": "task",
  "task_id": "uuid"
}
```

**Server → Client (updates):**

```json
{
  "event": "task_update",
  "data": {
    "task_id": "uuid",
    "status": "processing",
    "progress": 67,
    "timestamp": "2025-11-12T12:03:00Z"
  }
}
```

---

## 8. Admin API (опционально)

### Base URL: `/api/v1/admin`

#### 8.1 Управление парсерами

```http
POST /api/v1/admin/parsers/{parser_id}/restart
Authorization: Bearer {admin_token}
```

#### 8.2 Просмотр логов

```http
GET /api/v1/admin/logs
Authorization: Bearer {admin_token}
```

---

## Redis Schema

### Ключи для парсеров (обновляется worker'ом)

```
parsers:list → Set ["parser1", "parser2"]
parser:{parser_id}:info → Hash {name, status, description, ...}
parser:{parser_id}:health → String (timestamp последней проверки)
```

### Ключи для задач

```
task:{task_id}:status → String ("queued", "processing", "completed", "failed")
task:{task_id}:progress → Number (0-100)
task:{task_id}:logs → List [log entries]
user:{user_id}:tasks → Set [task_ids]
```

---

## RabbitMQ Queues

### Очереди

1. **`parsing_tasks`** - основная очередь задач

   - Routing key: `task.parse`
   - Message format:

   ```json
   {
     "task_id": "uuid",
     "parser_id": "parser1",
     "source_url": "https://example.com",
     "parameters": {...},
     "user_id": "uuid"
   }
   ```

2. **`parsing_results`** - результаты парсинга

   - Routing key: `task.result`
   - Message format:

   ```json
   {
     "task_id": "uuid",
     "status": "completed",
     "result_url": "s3://...",
     "completed_at": "2025-11-12T12:05:00Z"
   }
   ```

3. **`parser_health`** - health checks парсеров
   - Routing key: `parser.health`

---

## PostgreSQL Schema

### Таблица: users

```sql
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE
);
```

### Таблица: tasks

```sql
CREATE TABLE tasks (
    task_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(user_id),
    parser_id VARCHAR(50) NOT NULL,
    source_url TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'queued',
    parameters JSONB,
    result_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    error_message TEXT
);

CREATE INDEX idx_tasks_user_id ON tasks(user_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_created_at ON tasks(created_at DESC);
```

### Таблица: files

```sql
CREATE TABLE files (
    file_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(user_id),
    filename VARCHAR(255) NOT NULL,
    s3_key TEXT NOT NULL,
    size BIGINT,
    mime_type VARCHAR(100),
    uploaded_at TIMESTAMP DEFAULT NOW()
);
```

---

## Error Responses

### Стандартный формат ошибок

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {
      "field": "Additional context"
    }
  }
}
```

### HTTP Status Codes

- `200 OK` - успешный запрос
- `201 Created` - ресурс создан
- `400 Bad Request` - неверный запрос
- `401 Unauthorized` - не авторизован
- `403 Forbidden` - доступ запрещен
- `404 Not Found` - ресурс не найден
- `409 Conflict` - конфликт (дубликат)
- `422 Unprocessable Entity` - ошибка валидации
- `429 Too Many Requests` - превышен rate limit
- `500 Internal Server Error` - внутренняя ошибка сервера
- `503 Service Unavailable` - сервис недоступен

---

## Rate Limiting

- **Authentication endpoints**: 5 requests per minute
- **Standard API endpoints**: 100 requests per minute
- **File upload**: 10 requests per minute
- **WebSocket connections**: 5 concurrent connections per user

---

## Security

### JWT Token

- **Access Token**: 1 hour expiry
- **Refresh Token**: 7 days expiry
- Algorithm: HS256 или RS256
- Claims: `user_id`, `email`, `exp`, `iat`

### API Key (опционально)

```http
X-API-Key: your_api_key_here
```

---

## Environment Variables

```env
# Server
PORT=8080
ENV=production

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=myapp
DB_USER=postgres
DB_PASSWORD=password

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/

# S3 Yandex Cloud
S3_ENDPOINT=https://storage.yandexcloud.net
S3_BUCKET=my-bucket
S3_ACCESS_KEY=your_access_key
S3_SECRET_KEY=your_secret_key

# JWT
JWT_SECRET=your_secret_key
JWT_EXPIRY=3600

# CORS
CORS_ORIGINS=http://localhost:3000,https://yourdomain.com
```

---

## Deployment Notes

1. **Docker Compose** для локальной разработки
2. **Kubernetes** для production
3. **CI/CD** через GitHub Actions или GitLab CI
4. **Мониторинг**: Prometheus + Grafana
5. **Логирование**: ELK Stack или Loki

---

## Next Steps

1. Реализовать базовую структуру Go API
2. Настроить подключения к PostgreSQL, Redis, RabbitMQ
3. Реализовать authentication middleware
4. Создать worker для обработки задач
5. Интегрировать парсеры
6. Настроить S3 storage
7. Добавить WebSocket для real-time updates
8. Написать тесты
9. Подготовить Docker окружение
10. Настроить CI/CD
