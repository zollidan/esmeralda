#!/bin/bash
# Скрипт настройки сервера для деплоя FastAPI через GitHub Actions

set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Проверка прав root
if [[ $EUID -ne 0 ]]; then
   print_error "Этот скрипт должен быть запущен с правами root"
   exit 1
fi

print_status "Начинаем настройку сервера для деплоя FastAPI приложения..."

# Обновление системы
print_status "Обновление системы..."
apt update && apt upgrade -y

# Установка необходимых пакетов
print_status "Установка необходимых пакетов..."
apt install -y \
    curl \
    wget \
    git \
    unzip \
    software-properties-common \
    apt-transport-https \
    ca-certificates \
    gnupg \
    lsb-release \
    fail2ban \
    ufw \
    htop

# Установка Docker
print_status "Установка Docker..."
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com -o get-docker.sh
    sh get-docker.sh
    usermod -aG docker ubuntu
    rm get-docker.sh
else
    print_warning "Docker уже установлен"
fi

# Установка Docker Compose
print_status "Установка Docker Compose V2..."
if ! command -v docker-compose &> /dev/null; then
    mkdir -p /usr/local/lib/docker/cli-plugins
    curl -SL https://github.com/docker/compose/releases/latest/download/docker-compose-linux-x86_64 -o /usr/local/lib/docker/cli-plugins/docker-compose
    chmod +x /usr/local/lib/docker/cli-plugins/docker-compose
    ln -sf /usr/local/lib/docker/cli-plugins/docker-compose /usr/local/bin/docker-compose
else
    print_warning "Docker Compose уже установлен"
fi

# Настройка firewall
print_status "Настройка firewall..."
ufw --force reset
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow 80
ufw allow 443
ufw --force enable

# Настройка fail2ban
print_status "Настройка fail2ban..."
cat > /etc/fail2ban/jail.local << EOF
[DEFAULT]
bantime = 3600
findtime = 600
maxretry = 5

[sshd]
enabled = true
port = ssh
logpath = /var/log/auth.log
maxretry = 3
EOF

systemctl enable fail2ban
systemctl start fail2ban

# Создание пользователя deploy
DEPLOY_USER="deploy"
print_status "Создание пользователя $DEPLOY_USER..."
if ! id "$DEPLOY_USER" &>/dev/null; then
    useradd -m -s /bin/bash $DEPLOY_USER
    usermod -aG docker $DEPLOY_USER
    
    # Создание SSH ключей
    sudo -u $DEPLOY_USER ssh-keygen -t ed25519 -f /home/$DEPLOY_USER/.ssh/id_ed25519 -N ""
    sudo -u $DEPLOY_USER touch /home/$DEPLOY_USER/.ssh/authorized_keys
    chmod 700 /home/$DEPLOY_USER/.ssh
    chmod 600 /home/$DEPLOY_USER/.ssh/authorized_keys
    chown -R $DEPLOY_USER:$DEPLOY_USER /home/$DEPLOY_USER/.ssh
else
    print_warning "Пользователь $DEPLOY_USER уже существует"
fi

# Создание директории проекта
print_status "Создание директории для FastAPI проекта..."
PROJECT_DIR="/opt/fastapi-app"
mkdir -p $PROJECT_DIR
chown $DEPLOY_USER:$DEPLOY_USER $PROJECT_DIR

# Создание директорий для Traefik
print_status "Настройка директорий для Traefik..."
mkdir -p $PROJECT_DIR/traefik/logs
touch $PROJECT_DIR/traefik/acme.json
chmod 600 $PROJECT_DIR/traefik/acme.json
chown -R $DEPLOY_USER:$DEPLOY_USER $PROJECT_DIR/traefik

# Создание Traefik конфигураций
print_status "Создание конфигураций Traefik..."

# traefik.yml
cat > $PROJECT_DIR/traefik/traefik.yml << EOF
api:
  dashboard: true
  insecure: false

entryPoints:
  http:
    address: ":80"
  https:
    address: ":443"

providers:
  docker:
    endpoint: "unix:///var/run/docker.sock"
    exposedByDefault: false
    network: proxy
  file:
    filename: /config.yml

certificatesResolvers:
  cloudflare:
    acme:
      email: \${CF_API_EMAIL}
      storage: /acme.json
      dnsChallenge:
        provider: cloudflare
        resolvers:
          - "1.1.1.1:53"
          - "1.0.0.1:53"

log:
  filePath: "/var/log/traefik/traefik.log"
  level: "INFO"

accessLog:
  filePath: "/var/log/traefik/access.log"
EOF

# config.yml
cat > $PROJECT_DIR/traefik/config.yml << EOF
http:
  middlewares:
    https-redirect:
      redirectScheme:
        scheme: https
        permanent: true
    secure-headers:
      headers:
        sslRedirect: true
        forceSTSHeader: true
        stsIncludeSubdomains: true
        stsPreload: true
        stsSeconds: 31536000
EOF

# Пример .env файла
cat > $PROJECT_DIR/.env.example << EOF
# Cloudflare credentials for Traefik
CF_API_KEY=your_cloudflare_api_key
CF_API_EMAIL=your_email@example.com

# PostgreSQL configuration
POSTGRES_USER=postgres
POSTGRES_PASSWORD=secure_password_here
POSTGRES_DB=fastapi_db

# MinIO configuration
MINIO_ROOT_USER=minio
MINIO_ROOT_PASSWORD=minio_password

# Celery/Redis configuration
CELERY_BROKER_URL=redis://redis:6379/0
CELERY_RESULT_BACKEND=redis://redis:6379/0

# Application settings
DEBUG=False
SECRET_KEY=your_secret_key_here
ALLOWED_HOSTS=api.aaf-bet.ru,localhost,127.0.0.1
EOF

# Настройка логирования Docker
print_status "Настройка логирования Docker..."
mkdir -p /etc/docker
cat > /etc/docker/daemon.json << EOF
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  },
  "storage-driver": "overlay2"
}
EOF

# Создание скрипта для обновления
print_status "Создание скрипта для ручного обновления..."
cat > $PROJECT_DIR/update.sh << 'EOF'
#!/bin/bash
# Скрипт для ручного обновления FastAPI приложения

cd "$(dirname "$0")"

# Получаем последнюю версию из репозитория
git pull

# Перезапускаем сервисы
docker compose -f docker-compose.prod.yaml down
docker compose -f docker-compose.prod.yaml up -d

echo "✅ Приложение обновлено!"
EOF

chmod +x $PROJECT_DIR/update.sh
chown $DEPLOY_USER:$DEPLOY_USER $PROJECT_DIR/update.sh

# Создание скрипта для мониторинга
print_status "Создание скрипта для мониторинга..."
cat > $PROJECT_DIR/monitor.sh << 'EOF'
#!/bin/bash
# Скрипт для проверки состояния сервисов

echo "🔍 Проверка статуса контейнеров:"
docker ps

echo -e "\n🔄 Статус контейнеров docker-compose:"
docker compose -f docker-compose.prod.yaml ps

echo -e "\n📊 Использование ресурсов:"
docker stats --no-stream

echo -e "\n📜 Недавние логи приложения:"
docker compose -f docker-compose.prod.yaml logs --tail=50 web
EOF

chmod +x $PROJECT_DIR/monitor.sh
chown $DEPLOY_USER:$DEPLOY_USER $PROJECT_DIR/monitor.sh

# Перезапуск Docker
systemctl restart docker
systemctl enable docker

# Настройка cron-задания для очистки неиспользуемых Docker образов
print_status "Настройка автоматической очистки Docker..."
cat > /etc/cron.weekly/docker-cleanup << 'EOF'
#!/bin/bash
# Очистка неиспользуемых Docker образов и контейнеров
docker system prune -af --volumes
EOF

chmod +x /etc/cron.weekly/docker-cleanup

# Вывод информации о завершении
print_status "Настройка сервера завершена!"
echo ""
echo "🔑 SSH ключ для GitHub Actions (добавьте в секреты GitHub):"
echo "=========================================="
cat /home/$DEPLOY_USER/.ssh/id_ed25519
echo ""
echo "=========================================="
echo ""
echo "📋 Настройка GitHub Actions:"
echo "1. Добавьте SSH ключ выше в секреты GitHub Actions (SSH_PRIVATE_KEY)"
echo "2. Добавьте другие необходимые секреты:"
echo "   - HOST: $(curl -s ifconfig.me)"
echo "   - USERNAME: $DEPLOY_USER"
echo "   - PORT: 22"
echo ""
echo "3. Скопируйте пример .env файла и заполните его:"
echo "   cp $PROJECT_DIR/.env.example $PROJECT_DIR/.env"
echo ""
echo "🌐 Настройка домена:"
echo "1. Убедитесь, что ваш домен aaf-bet.ru указывает на IP: $(curl -s ifconfig.me)"
echo "2. Настройте DNS для поддоменов: api, flower, s3, traefik"
echo ""
echo "🚀 Инициализация проекта:"
echo "1. Клонируйте ваш репозиторий:"
echo "   cd $PROJECT_DIR && git clone https://github.com/ваш_репозиторий ."
echo "2. Запустите docker compose -f docker-compose.prod.yaml up -d"
echo ""
echo "🎉 Сервер готов!"