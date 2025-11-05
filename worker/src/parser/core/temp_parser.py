import pandas as pd
import uuid
import logging
from typing import Dict, Any

logger = logging.getLogger(__name__)


def parse_soccerway(task_data: Dict[str, Any]) -> pd.DataFrame:
    """
    Ваш парсер - заменишь на свою логику

    Args:
        task_data: Данные задачи из RabbitMQ

    Returns:
        pd.DataFrame с результатами парсинга
    """
    logger.info(
        f"Начало парсинга Soccerway для задачи {task_data.get('task_id')}")

    # ========================================
    # ЗДЕСЬ ТВОЯ ЛОГИКА ПАРСИНГА
    # ========================================

    # Пример:
    d = {'col1': [1, 2], 'col2': [3, 4]}
    df = pd.DataFrame(data=d)

    logger.info(f"Парсинг завершен. Получено {len(df)} строк")

    return df


def parse_another_site(task_data: Dict[str, Any]) -> pd.DataFrame:
    """
    Еще один парсер - можешь добавить сколько нужно
    """
    logger.info(
        f"Начало парсинга другого сайта для задачи {task_data.get('task_id')}")

    # Твоя логика
    df = pd.DataFrame({'data': [1, 2, 3]})

    return df
