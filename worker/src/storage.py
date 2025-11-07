from worker.src.settings import settings


class S3Storage:

    def __init__(self):
        self.bucket_name = settings.S3_BUCKET_NAME
        self.aws_access_key_id = settings.AWS_ACCESS_KEY_ID
        self.aws_secret_access_key = settings.AWS_SECRET_ACCESS_KEY

    def _create_client(self):
        import boto3

        return boto3.client(
            "s3",
            aws_access_key_id=self.aws_access_key_id,
            aws_secret_access_key=self.aws_secret_access_key,
        )
