import asyncio
from fastapi import APIRouter, BackgroundTasks
import requests

from app.parser.parsers_wrapper import soccerway_wrapper
from app.parser.server_parser_funcions import get_files_s3

router = APIRouter(prefix='/api', tags=['API'])

# @router.get('/')
# def index_page():
    
    
    
#     return INDEX_MESSAGE

# Создаем глобальный объект цикла событий
loop = asyncio.get_event_loop()

@router.get('/parser/soccerway')
def run_soccerway_url_method(background_tasks: BackgroundTasks, date: str):
    
    background_tasks.add_task(soccerway_wrapper, date)
 
    return {"message": "soccerway work started"}

@router.get('/parser/soccerway/connection-test')
def run_soccerway_test_connection():
    
    url = "https://ru.soccerway.com/a/block_h2h_matches?block_id=page_match_1_block_h2hsection_head2head_7_block_h2h_matches_1&action=changePage&callback_params=%7B%22page%22%3A+-1%2C+%22block_service_id%22%3A+%22match_h2h_comparison_block_h2hmatches%22%2C+%22team_A_id%22%3A+43000%2C+%22team_B_id%22%3A+43009%7D&params=%7B%22page%22%3A+0%7D"

    response = requests.get(url, headers={'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:102.0) Gecko/20100101 Firefox/131.0'})
 
    return response.status_code

@router.get('/files')
def get_files():
    
    files = get_files_s3()
    
    
    return files
    

# @router.post('/files/add', status_code=status.HTTP_201_CREATED)
# async def add_file(payload: schemas.FileBody, db: Session = Depends(get_db)):
    
#     file = models.Files(url=payload.url)
    
#     db.add(file)
#     db.commit()
#     db.refresh(file)
  
#     return file