from aiogram import Router, F
from aiogram.filters import CommandStart
from aiogram.types import Message
from app.api.dao import TelegramUserDAO
from app.bot.utils.utils import greet_user, get_about_us_text

user_router = Router()

@user_router.message(CommandStart())
async def cmd_start(message: Message) -> None:
    """
    Обрабатывает команду /start.
    """
    user = await TelegramUserDAO.find_one_or_none(telegram_id=message.from_user.id)
    
    if not user:
        await TelegramUserDAO.add(
            telegram_id=message.from_user.id,
            first_name=message.from_user.first_name,
            username=message.from_user.username
        )
    
    await greet_user(message, is_new_user=not user)
    
@user_router.message(F.text == '🔙 Назад')
async def cmd_back_home(message: Message) -> None:
    """
    Обрабатывает нажатие кнопки "Назад".
    """
    await greet_user(message, is_new_user=False)
    
    
