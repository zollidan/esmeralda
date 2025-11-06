import logging
import sys


class Logger:
    """Wrapper for logging with proper configuration"""

    def __init__(self, level: str = "INFO", filename: str = "worker.log"):
        """Configure logging with file and console output"""
        # Create logger
        self._logger = logging.getLogger("worker")
        self._logger.setLevel(getattr(logging, level.upper()))

        # Prevent duplicate logs
        if self._logger.handlers:
            return

        # File handler
        file_handler = logging.FileHandler(filename, encoding='utf-8')
        file_handler.setLevel(getattr(logging, level.upper()))
        file_formatter = logging.Formatter(
            '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
        )
        file_handler.setFormatter(file_formatter)

        # Console handler
        console_handler = logging.StreamHandler(sys.stdout)
        console_handler.setLevel(getattr(logging, level.upper()))
        console_formatter = logging.Formatter(
            '%(asctime)s - %(levelname)s - %(message)s'
        )
        console_handler.setFormatter(console_formatter)

        # Add handlers
        self._logger.addHandler(file_handler)
        self._logger.addHandler(console_handler)

    def info(self, message: str):
        self._logger.info(message)

    def error(self, message: str):
        self._logger.error(message)

    def warning(self, message: str):
        self._logger.warning(message)

    def debug(self, message: str):
        self._logger.debug(message)


logger = Logger()
