import boto3
import os
from typing import Optional
from botocore.exceptions import ClientError, NoCredentialsError

from worker.src.logger import logger
from worker.src.settings import settings


class S3Storage:
    """S3 Storage handler for file uploads"""

    def __init__(self):
        self.bucket_name = settings.S3_BUCKET
        self.region = settings.S3_REGION
        
        # Initialize S3 client
        try:
            self.s3_client = boto3.client(
                's3',
                region_name=self.region,
                aws_access_key_id=settings.AWS_ACCESS_KEY_ID,
                aws_secret_access_key=settings.AWS_SECRET_ACCESS_KEY
            )
            self._validate_bucket()
            logger.info(f"S3 storage initialized successfully for bucket: {self.bucket_name}")
        except NoCredentialsError:
            logger.error("AWS credentials not found. Please check your environment variables.")
            raise
        except Exception as e:
            logger.error(f"Failed to initialize S3 storage: {e}")
            raise

    def _validate_bucket(self):
        """Validate that the bucket exists and is accessible"""
        try:
            self.s3_client.head_bucket(Bucket=self.bucket_name)
            logger.info(f"Bucket {self.bucket_name} exists and is accessible")
        except ClientError as e:
            error_code = e.response['Error']['Code']
            if error_code == '404':
                logger.error(f"Bucket {self.bucket_name} does not exist")
                raise
            elif error_code == '403':
                logger.error(f"Access denied to bucket {self.bucket_name}")
                raise
            else:
                logger.error(f"Error accessing bucket {self.bucket_name}: {e}")
                raise

    def upload(self, file_path: str, task_id: str) -> str:
        """
        Upload file to S3 storage
        
        Args:
            file_path: Local path to the file to upload
            task_id: Task ID to use in S3 key
            
        Returns:
            str: S3 key where the file was uploaded
        """
        if not os.path.exists(file_path):
            raise FileNotFoundError(f"File not found: {file_path}")

        # Create S3 key using task_id and original filename
        filename = os.path.basename(file_path)
        s3_key = f"tasks/{task_id}/{filename}"
        
        try:
            logger.info(f"Uploading {file_path} to s3://{self.bucket_name}/{s3_key}")
            
            self.s3_client.upload_file(
                file_path,
                self.bucket_name,
                s3_key,
                ExtraArgs={
                    'ServerSideEncryption': 'AES256',  # Enable server-side encryption
                    'ContentType': self._get_content_type(filename)
                }
            )
            
            logger.info(f"Successfully uploaded {file_path} to s3://{self.bucket_name}/{s3_key}")
            return s3_key
            
        except ClientError as e:
            logger.error(f"Failed to upload {file_path} to S3: {e}")
            raise
        except Exception as e:
            logger.error(f"Unexpected error uploading {file_path} to S3: {e}")
            raise

    def download(self, s3_key: str, local_path: str) -> bool:
        """
        Download file from S3 storage
        
        Args:
            s3_key: S3 key of the file to download
            local_path: Local path where to save the file
            
        Returns:
            bool: True if download was successful, False otherwise
        """
        try:
            logger.info(f"Downloading s3://{self.bucket_name}/{s3_key} to {local_path}")
            
            # Create directory if it doesn't exist
            os.makedirs(os.path.dirname(local_path), exist_ok=True)
            
            self.s3_client.download_file(self.bucket_name, s3_key, local_path)
            logger.info(f"Successfully downloaded s3://{self.bucket_name}/{s3_key} to {local_path}")
            return True
            
        except ClientError as e:
            logger.error(f"Failed to download {s3_key} from S3: {e}")
            return False
        except Exception as e:
            logger.error(f"Unexpected error downloading {s3_key} from S3: {e}")
            return False

    def delete(self, s3_key: str) -> bool:
        """
        Delete file from S3 storage
        
        Args:
            s3_key: S3 key of the file to delete
            
        Returns:
            bool: True if deletion was successful, False otherwise
        """
        try:
            logger.info(f"Deleting s3://{self.bucket_name}/{s3_key}")
            
            self.s3_client.delete_object(Bucket=self.bucket_name, Key=s3_key)
            logger.info(f"Successfully deleted s3://{self.bucket_name}/{s3_key}")
            return True
            
        except ClientError as e:
            logger.error(f"Failed to delete {s3_key} from S3: {e}")
            return False
        except Exception as e:
            logger.error(f"Unexpected error deleting {s3_key} from S3: {e}")
            return False

    def file_exists(self, s3_key: str) -> bool:
        """
        Check if file exists in S3 storage
        
        Args:
            s3_key: S3 key to check
            
        Returns:
            bool: True if file exists, False otherwise
        """
        try:
            self.s3_client.head_object(Bucket=self.bucket_name, Key=s3_key)
            return True
        except ClientError as e:
            if e.response['Error']['Code'] == '404':
                return False
            else:
                logger.error(f"Error checking if {s3_key} exists in S3: {e}")
                raise

    def get_file_url(self, s3_key: str, expiration: int = 3600) -> Optional[str]:
        """
        Generate a presigned URL for the S3 object
        
        Args:
            s3_key: S3 key of the file
            expiration: URL expiration time in seconds (default: 1 hour)
            
        Returns:
            str: Presigned URL or None if error
        """
        try:
            url = self.s3_client.generate_presigned_url(
                'get_object',
                Params={'Bucket': self.bucket_name, 'Key': s3_key},
                ExpiresIn=expiration
            )
            return url
        except ClientError as e:
            logger.error(f"Failed to generate presigned URL for {s3_key}: {e}")
            return None

    def _get_content_type(self, filename: str) -> str:
        """
        Determine content type based on file extension
        
        Args:
            filename: Name of the file
            
        Returns:
            str: MIME content type
        """
        ext = os.path.splitext(filename)[1].lower()
        content_types = {
            '.xlsx': 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
            '.xls': 'application/vnd.ms-excel',
            '.csv': 'text/csv',
            '.json': 'application/json',
            '.txt': 'text/plain',
            '.pdf': 'application/pdf',
            '.png': 'image/png',
            '.jpg': 'image/jpeg',
            '.jpeg': 'image/jpeg',
            '.gif': 'image/gif',
        }
        return content_types.get(ext, 'application/octet-stream')

    def cleanup_local_file(self, file_path: str) -> bool:
        """
        Delete local file after successful upload to S3
        
        Args:
            file_path: Path to local file to delete
            
        Returns:
            bool: True if deletion was successful, False otherwise
        """
        try:
            if os.path.exists(file_path):
                os.remove(file_path)
                logger.info(f"Successfully deleted local file: {file_path}")
                return True
            return False
        except Exception as e:
            logger.error(f"Failed to delete local file {file_path}: {e}")
            return False