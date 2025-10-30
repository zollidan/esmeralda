import boto3
from botocore.exceptions import ClientError
from io import BytesIO
import pandas as pd
from typing import Dict, Union
import logging

logger = logging.getLogger(__name__)


class S3Uploader:
    """Класс для загрузки файлов в S3/MinIO"""

    def __init__(self, endpoint_url: str, access_key: str, secret_key: str, bucket_name: str):
        self.bucket_name = bucket_name
        self.s3_client = boto3.client(
            's3',
            endpoint_url=endpoint_url,
            aws_access_key_id=access_key,
            aws_secret_access_key=secret_key
        )
        self._ensure_bucket_exists()

    def _ensure_bucket_exists(self):
        """Проверяет существование bucket и создает если нужно"""
        try:
            self.s3_client.head_bucket(Bucket=self.bucket_name)
            logger.info(f"Bucket '{self.bucket_name}' существует")
        except ClientError as e:
            error_code = e.response['Error']['Code']
            if error_code == '404':
                logger.info(f"Создаем bucket '{self.bucket_name}'")
                self.s3_client.create_bucket(Bucket=self.bucket_name)
            else:
                logger.error(f"Ошибка проверки bucket: {e}")
                raise

    def upload_dataframe(self, df: pd.DataFrame, filename: str) -> Dict[str, Union[bool, str]]:
        """
        Загружает DataFrame в S3 как Excel файл

        Args:
            df: DataFrame для сохранения
            filename: Имя файла в S3

        Returns:
            Dict с результатом: {'success': bool, 'filename': str, 'bucket': str, 'error': str}
        """
        try:
            buffer = BytesIO()
            df.to_excel(buffer, index=False, engine='openpyxl')
            buffer.seek(0)

            self.s3_client.upload_fileobj(
                buffer,
                self.bucket_name,
                filename,
                ExtraArgs={
                    'ContentType': 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
                }
            )

            logger.info(f"Файл '{filename}' успешно загружен в S3")

            return {
                'success': True,
                'filename': filename,
                'bucket': self.bucket_name,
                'path': f"{self.bucket_name}/{filename}"
            }

        except ClientError as e:
            error_msg = f"Ошибка загрузки в S3: {str(e)}"
            logger.error(error_msg)
            return {
                'success': False,
                'error': error_msg
            }
        except Exception as e:
            error_msg = f"Неожиданная ошибка при загрузке: {str(e)}"
            logger.error(error_msg)
            return {
                'success': False,
                'error': error_msg
            }

    def upload_file(self, file_path: str, s3_filename: str = None) -> Dict[str, Union[bool, str]]:
        """
        Загружает локальный файл в S3

        Args:
            file_path: Путь к локальному файлу
            s3_filename: Имя файла в S3 (если None, используется имя из file_path)

        Returns:
            Dict с результатом
        """
        if s3_filename is None:
            s3_filename = file_path.split('/')[-1]

        try:
            self.s3_client.upload_file(
                file_path, self.bucket_name, s3_filename)

            logger.info(f"Файл '{s3_filename}' успешно загружен в S3")

            return {
                'success': True,
                'filename': s3_filename,
                'bucket': self.bucket_name,
                'path': f"{self.bucket_name}/{s3_filename}"
            }

        except ClientError as e:
            error_msg = f"Ошибка загрузки в S3: {str(e)}"
            logger.error(error_msg)
            return {
                'success': False,
                'error': error_msg
            }
        except Exception as e:
            error_msg = f"Неожиданная ошибка: {str(e)}"
            logger.error(error_msg)
            return {
                'success': False,
                'error': error_msg
            }
