from aiogram import Router, F
from aiogram.filters import CommandStart, Command
from aiogram.types import Message
from app.bot.keyboards.kbs import *
from app.parser.server_parser_funcions import get_files_s3

tg_user_router = Router()

welcome_text = (
        "👋 Привет! Добро пожаловать в мир QR-кодов! 🌟\n\n"
        "Я ваш помощник по работе с QR-кодами. Вот что я умею:\n\n"
        "📷 Сканировать QR-коды в реальном времени\n"
        "📤 Распознавать QR-коды на загруженных фото\n"
        "✨ Создавать новые QR-коды\n\n"
        "Выберите нужную функцию в меню ниже и начните работу! 🚀"
    )

@tg_user_router.message(CommandStart())
async def cmd_start(message: Message) -> None:
    
    await message.answer(welcome_text, reply_markup=main_keyboards())
    
@tg_user_router.message(F.text == 'Файлы')
async def cmd_back_home(message: Message) -> None:
    
    files = get_files_s3()
    
    await message.answer(files[0], reply_markup=files_keyboard())
    
# @tg_user_router.message(F.text == 'Сохранить последний файл')
# async def cmd_save_file(message: Message) -> None:
    
#     files = get_files_s3()
    
#     await message.answer(str(files), reply_markup=files_keyboard())
    
@tg_user_router.message(F.text == 'Назад')
async def cmd_back_home(message: Message) -> None:
    
    
    
    await message.answer(welcome_text, reply_markup=main_keyboards())
    
@tg_user_router.message(F.text == 'Парсеры')
async def cmd_parsers_home(message: Message) -> None:
    
    await message.answer("страница запуска парсеров",reply_markup=parsers_keyboard())
    
