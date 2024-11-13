#!/bin/bash

sudo apt update

sudo apt install python3.12-venv

python3 -m venv venv

# Указываем путь к виртуальному окружению
source venv/bin/activate

pip install -r requirements.txt

# # Запускаем приложение через uvicorn
# uvicorn app.main:app --reload --host 0.0.0.0 --port 80