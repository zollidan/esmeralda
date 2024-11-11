from aiogram.types import WebAppInfo, InlineKeyboardMarkup, ReplyKeyboardMarkup
from aiogram.utils.keyboard import InlineKeyboardBuilder, ReplyKeyboardBuilder
from app.config import settings

def main_keyboards() -> ReplyKeyboardMarkup:
    kb = ReplyKeyboardBuilder()
    generate_url = settings.BASE_SITE
    kb.button(text="Парсеры")
    kb.button(text="Файлы")
    kb.adjust(1)
    return kb.as_markup(resize_keyboard=True)

def files_keyboard() -> ReplyKeyboardMarkup:
    kb = ReplyKeyboardBuilder()
    kb.button(text="Сохранить последний файл")
    kb.button(text="Назад")
    kb.adjust(1)
    return kb.as_markup(resize_keyboard=True)

def parsers_keyboard() -> ReplyKeyboardMarkup:
    kb = ReplyKeyboardBuilder()
    kb.button(text="Запустить Soccerway")
    kb.adjust(1)
    return kb.as_markup(resize_keyboard=True)

