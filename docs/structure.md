project-root/
├── docker-compose.yml
├── .env
├── Makefile
├── README.md
│
├── deploy/ # devops/infra конфигурации
│ ├── Dockerfiles/
│ │ ├── go-api.Dockerfile
│ │ ├── python-worker.Dockerfile
│ │ └── nextjs.Dockerfile
│ ├── nginx/
│ │ └── nginx.conf
│ └── k8s/ # если решишь потом перейти на Kubernetes
│ ├── deployments/
│ ├── services/
│ └── secrets/
│
├── services/
│ ├── api/ # Go сервер (chi + gorm)
│ │ ├── cmd/
│ │ │ └── main.go # входная точка
│ │ ├── internal/
│ │ │ ├── app/ # маршруты, бизнес-логика
│ │ │ ├── models/ # gorm модели
│ │ │ ├── db/ # инициализация postgres
│ │ │ ├── mq/ # publisher для RabbitMQ
│ │ │ ├── storage/ # работа с S3 (minio)
│ │ │ ├── scheduler/ # cron/планировщик задач
│ │ │ └── config/ # env, настройки
│ │ ├── go.mod
│ │ ├── go.sum
│ │ └── Dockerfile -> ../../deploy/Dockerfiles/go-api.Dockerfile
│ │
│ ├── worker/ # Python-воркер (парсер)
│ │ ├── src/
│ │ │ ├── main.py
│ │ │ ├── parser/
│ │ │ │ ├── **init**.py
│ │ │ │ ├── selenium_parser.py
│ │ │ │ ├── utils.py
│ │ │ │ └── s3_uploader.py
│ │ │ └── mq_consumer.py
│ │ ├── requirements.txt
│ │ └── Dockerfile -> ../../deploy/Dockerfiles/python-worker.Dockerfile
│ │
│ └── frontend/ # Next.js (TypeScript)
│ ├── src/
│ │ ├── pages/
│ │ ├── components/
│ │ ├── lib/ # API SDK
│ │ └── styles/
│ ├── package.json
│ ├── next.config.js
│ └── Dockerfile -> ../../deploy/Dockerfiles/nextjs.Dockerfile
│
├── scripts/ # скрипты локального запуска / отладки
│ ├── init_db.sql
│ ├── reset_queues.sh
│ ├── run_api_local.sh
│ └── run_worker_local.sh
│
└── docs/
├── architecture.md # описание архитектуры
├── api_spec.md # документация REST API
├── worker_flow.md # как работает очередь/воркеры
└── deployment.md

---

Frontend (Next.js)
↓
Go API (chi + Gorm)
↓
RabbitMQ (очередь задач)
↓
Python workers (Docker)
↓
S3 + Postgres

---

Next.js
↓
Go (chi) + Gorm
├── POST /api/tasks → кладёт задание в RabbitMQ
├── cron → планирует задания
└── GET /tasks/status
Python worker(s)
↳ подписан на очередь
↳ обрабатывает задания (парсинг)
↳ результат в S3 + отчёт в Postgres через REST/gRPC
