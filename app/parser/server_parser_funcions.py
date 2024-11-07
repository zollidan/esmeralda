import json
import boto3
import requests
import os
from dotenv import load_dotenv

load_dotenv()

def send_file_record_s3(file_name):

    session = boto3.session.Session()
    s3 = session.client(
        service_name='s3',
        endpoint_url='https://storage.yandexcloud.net',
        aws_access_key_id=os.getenv("AWS_ID"),
        aws_secret_access_key=os.getenv("AWS_SECRET_KEY"),
    )
    
    ## Из файла
    s3.upload_file(file_name, 'esmeralda', file_name)


    # create_file_body = json.dumps({
    #     "url": create_file_record_href
    # })
    
    # create_file_record = requests.post(f"http://localhost:8000/files/add", data=create_file_body)
    
# def get_files_list():
    
#     session = boto3.session.Session()
    
#     s3 = session.client(
#         service_name='s3',
#         endpoint_url='https://storage.yandexcloud.net',
#         aws_access_key_id=aws_access_key_id,
#         aws_secret_access_key=aws_secret_access_key,
#     )
    
#     obj_list = []
    
#     for key in s3.list_objects(Bucket='bucket-name')['Contents']:
     
#         obj_list.append(key['Key'])
#     return obj_list
        