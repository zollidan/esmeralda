"""
    бот для телеграмма, работает параллельно с api
    TODO: добавить обработку команд и взаимодействие с flower api tasks
"""

import asyncio
import logging
import sys
import re
from os import getenv

from aiogram import Bot, Dispatcher, types
from aiogram.filters import CommandStart, Command
from aiogram.types import Message
from aiogram.client.default import DefaultBotProperties
from aiogram.enums import ParseMode

import requests
from sqlmodel import Session, select
from main import APP_VERSION, settings, engine, File, s3_client
from tasks import run_soccerway_1, run_soccerway_2

TOKEN = settings.TELEGRAM_BOT_TOKEN

dp = Dispatcher()


def is_valid_date(date_str):
    """Проверка валидности даты в формате YYYY-MM-DD"""
    return bool(re.match(r'^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$', date_str))

# TELEGRAM BOT HANDLERS


@dp.message(CommandStart())
async def command_start_handler(message: Message) -> None:
    """Обработчик команды /start"""
    await message.answer(f"Привет, {message.from_user.full_name}! 🤖\n\n"
                         "Я бот для управления задачами Esmeralda.\n\n"
                         "Доступные команды:\n"
                         "/help - показать справку\n"
                         "/status - проверить статус системы\n"
                         "/files - показать список файлов")


@dp.message(Command("help"))
async def command_help_handler(message: Message) -> None:
    """Обработчик команды /help"""
    help_text = """
🤖 <b>Справка по боту aaf-bet.ru</b>

<b>Доступные команды:</b>
/start - запустить бота
/help - показать эту справку
/status - проверить статус системы
/files - показать список загруженных файлов
/run_soccerway1 YYYY-MM-DD YYYY-MM-DD - запустить задачу Soccerway 1
/run_soccerway2 YYYY-MM-DD YYYY-MM-DD - запустить задачу Soccerway 2
/parsers_status - показать статус парсеров

<b>Примеры:</b>
<code>/run_soccerway1 2024-01-01 2024-01-31</code>
<code>/run_soccerway2 2024-02-01 2024-02-28</code>
    """
    await message.answer(help_text, parse_mode="HTML")


@dp.message(Command("status"))
async def command_status_handler(message: Message) -> None:
    """Обработчик команды /status"""
    try:
        # Проверяем подключение к базе данных
        with Session(engine) as session:
            files_count = len(session.exec(select(File)).all())

        # Проверяем подключение к MinIO
        try:
            s3_client.list_buckets()
            minio_status = "✅ Подключен"
        except Exception:
            minio_status = "❌ Ошибка подключения"

        status_text = f"""
📊 <b>Статус системы aaf-bet.ru</b>

🗄️ <b>База данных:</b> ✅ Подключена
📁 <b>Файлов в системе:</b> {files_count}
🗂️ <b>MinIO:</b> {minio_status}
🤖 <b>Режим бота:</b> Standalone
⚙️ <b>Версия API:</b> {APP_VERSION}
        """
        await message.answer(status_text, parse_mode="HTML")
    except Exception as e:
        await message.answer(f"❌ Ошибка при получении статуса: {str(e)}")


@dp.message(Command("files"))
async def command_files_handler(message: Message) -> None:
    """Обработчик команды /files"""
    try:
        with Session(engine) as session:
            files = session.exec(select(File)).all()

        if not files:
            await message.answer("📁 Файлов в системе пока нет")
            return

        files_text = "📁 <b>Файлы в системе:</b>\n\n"
        for file in files[:10]:  # Показываем только первые 10 файлов
            files_text += f"📄 <b>{file.name}</b>\n"
            files_text += f"🔗 ID: <code>{file.id}</code>\n"
            files_text += f"📅 Создан: {file.created_at.strftime('%Y-%m-%d %H:%M')}\n\n"

        if len(files) > 10:
            files_text += f"... и еще {len(files) - 10} файлов"

        await message.answer(files_text, parse_mode="HTML")
    except Exception as e:
        await message.answer(f"❌ Ошибка при получении списка файлов: {str(e)}")


@dp.message(Command("run_soccerway1"))
async def command_run_soccerway1_handler(message: Message) -> None:
    """Обработчик команды /run_soccerway1"""
    try:
        # Парсим аргументы команды
        args = message.text.split()[1:]

        if len(args) != 2:
            await message.answer("❌ Неверный формат команды!\n\n"
                                 "Использование: <code>/run_soccerway1 YYYY-MM-DD YYYY-MM-DD</code>\n"
                                 "Пример: <code>/run_soccerway1 2024-01-01 2024-01-31</code>",
                                 parse_mode="HTML")
            return

        date_start, date_end = args

        if not (is_valid_date(date_start) and is_valid_date(date_end)):
            await message.answer("❌ Неверный формат дат! Используйте формат YYYY-MM-DD")
            return

        # Запускаем задачу
        task = run_soccerway_1.delay(date_start, date_end)

        await message.answer(f"✅ Задача Soccerway 1 запущена!\n\n"
                             f"📅 Период: {date_start} - {date_end}\n"
                             f"🆔 ID задачи: <code>{task.id}</code>",
                             parse_mode="HTML")

    except Exception as e:
        logging.error(f"Error in run_soccerway1 command: {e}")
        await message.answer(f"❌ Ошибка при запуске задачи: {str(e)}")


@dp.message(Command("run_soccerway2"))
async def command_run_soccerway2_handler(message: Message) -> None:
    """Обработчик команды /run_soccerway2"""
    try:
        # Парсим аргументы команды
        args = message.text.split()[1:]  # Убираем саму команду

        if len(args) != 2:
            await message.answer("❌ Неверный формат команды!\n\n"
                                 "Использование: <code>/run_soccerway2 YYYY-MM-DD YYYY-MM-DD</code>\n"
                                 "Пример: <code>/run_soccerway2 2024-01-01 2024-01-31</code>",
                                 parse_mode="HTML")
            return

        date_start, date_end = args

        if not (is_valid_date(date_start) and is_valid_date(date_end)):
            await message.answer("❌ Неверный формат дат! Используйте формат YYYY-MM-DD")
            return

        # Запускаем задачу
        task = run_soccerway_2.delay(date_start, date_end)

        await message.answer(f"✅ Задача Soccerway 2 запущена!\n\n"
                             f"📅 Период: {date_start} - {date_end}\n"
                             f"🆔 ID задачи: <code>{task.id}</code>",
                             parse_mode="HTML")

    except Exception as e:
        logging.error(f"Error in run_soccerway2 command: {e}")
        await message.answer(f"❌ Ошибка при запуске задачи: {str(e)}")


@dp.message(Command("tasks"))
async def command_tasks_handler(message: Message) -> None:
    """Обработчик команды /tasks - показать статус задач"""
    try:
        # Здесь можно добавить логику для получения статуса задач из Celery
        # Например, через Flower API или Redis
        await message.answer("📋 <b>Статус задач:</b>\n\n"
                             "Для просмотра детального статуса задач используйте Flower:\n"
                             "🌸 <code>flower.aaf-bet.ru</code>",
                             parse_mode="HTML")
    except Exception as e:
        logging.error(f"Error in tasks command: {e}")
        await message.answer(f"❌ Ошибка при получении статуса задач: {str(e)}")


@dp.message()
async def echo_handler(message: Message) -> None:
    """Обработчик всех остальных сообщений"""
    await message.answer("🤖 Я не понимаю эту команду.\n\n"
                         "Используйте /help чтобы увидеть список доступных команд.")


async def main() -> None:
    """Основная функция запуска бота"""
    # Инициализация бота с настройками по умолчанию
    bot = Bot(token=TOKEN, default=DefaultBotProperties(
        parse_mode=ParseMode.HTML))

    logging.info("🤖 Запуск Telegram бота...")

    try:
        # Запуск polling
        await dp.start_polling(bot)
    except KeyboardInterrupt:
        logging.info("🤖 Получен сигнал остановки")
    except Exception as e:
        logging.error(f"❌ Ошибка при работе бота: {e}")
    finally:
        # Закрытие сессии бота
        await bot.session.close()
        logging.info("🤖 Telegram бот остановлен")

if __name__ == "__main__":
    # Настройка логирования
    logging.basicConfig(
        level=logging.INFO,
        format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
        stream=sys.stdout
    )

    # Запуск бота
    asyncio.run(main())
