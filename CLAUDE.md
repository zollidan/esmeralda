# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Project Does

**Esmeralda** fetches and processes football match statistics from [api.api-sport.ru](https://api.api-sport.ru). Users submit a date, the system fetches all matches for that date, calculates H2H / home / away / goal-total statistics per match (up to 20 parallel goroutines), stores results in PostgreSQL, and serves them via a Vue 3 frontend.

## Architecture

Event-driven two-service architecture:

1. **API service** (`cmd/api/`) — Gin HTTP server on `:8080`. Accepts `POST /api/tasks` (date input), publishes a `ParseTask` to the `tasks:parse` Redis Stream, and listens on `tasks:results` to update task status in PostgreSQL. Serves the embedded Vue frontend via catch-all route.

2. **Parser service** (`cmd/parser/`) — Redis Stream consumer (consumer group). Picks up tasks, calls api.api-sport.ru for that date's matches, processes them in parallel, writes game rows to PostgreSQL, then publishes `TaskResult` back to `tasks:results`.

**Message flow:** `POST /api/tasks` → Redis Stream `tasks:parse` → Parser → PostgreSQL + Redis Stream `tasks:results` → API updates Task status.

**Redis messages:**

- `ParseTask`: `{id, date: "YYYY-MM-DD", created_at}`
- `TaskResult`: `{task_id, status: pending|processing|done|error|failed|cancelled, error}`

**Database models:** `Task` (UUID, date, status, created_at) and `Game` (49 statistics columns per match).

## Commands

### Backend (Go)

```bash
# Run services
go run ./cmd/api/main.go
go run ./cmd/parser/main.go

# Build
go build ./cmd/api
go build ./cmd/parser

# Test all
go test ./...

# Test a single package
go test ./internal/stats/...
go test ./internal/processor/...
```

### Frontend

```bash
cd front
bun run dev       # Vite dev server
bun run build     # TypeScript check + Vite build
bun run preview   # Preview production build
```

### Infrastructure

```bash
# Start PostgreSQL + Redis
docker compose up -d

# Production-like setup
docker compose -f docker-compose.infrastructure.yml up -d
```

## Configuration

Copy `.env.infrastructure.example` to `.env`. Required variable:

| Variable          | Description                    |
| ----------------- | ------------------------------ |
| `SPORT_API_TOKEN` | API token for api.api-sport.ru |
| `SERVER_PORT`     | HTTP port (default: `8080`)    |
| `POSTGRES_*`      | PostgreSQL connection settings |
| `REDIS_ADDR`      | Redis address                  |

## Key Internal Packages

- `internal/api/` — HTTP client for api.api-sport.ru and response models
- `internal/processor/` — parallel match processing (semaphore, up to 20 goroutines; currently hardcoded to process first 20 matches in `processor.go:48`)
- `internal/stats/` — H2H, home/away, goal-total calculations; `stats.Row` maps to 49 Excel columns
- `internal/queue/` — Redis Streams producer/consumer with consumer groups
- `internal/export/` — Excel `.xlsx` export via `excelize` StreamWriter
- `internal/server/` — Gin route handlers
- `internal/static/` — embedded frontend assets (Go `embed`)

## Tech Stack

- **Go 1.25** with Gin, GORM, go-redis
- **PostgreSQL 16**, **Redis** (Streams with consumer groups)
- **Vue 3 + TypeScript**, built with Vite
- **Docker Compose** for local infrastructure
