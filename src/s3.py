import logging
import boto3
from botocore.config import Config
from botocore.exceptions import ClientError
from settings import settings

class S3Storage:
    def __init__(self) -> None:
        _config = Config(
            region_name=settings.AWS_REGION,
            retries={
                'max_attempts': 3,
                'mode': 'standard'
            },
            connect_timeout=10,
            read_timeout=30
        )
        
        self.s3_client = boto3.client(
            's3',
            aws_access_key_id=settings.AWS_ACCESS_KEY_ID,
            aws_secret_access_key=settings.AWS_SECRET_ACCESS_KEY,
            config=_config
        )
        
        self.__bucket_name = settings.AWS_BUCKET_NAME
    
    def upload_file(self, *, file_name: str, object_name: str) -> bool:
        try:
            self.s3_client.upload_file(file_name, self.__bucket_name, object_name)
        except ClientError as e:
            logging.error(e)
            return False
        return True
    
s3 = S3Storage()


