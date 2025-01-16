import requests
from fastapi import APIRouter

status_check_router = APIRouter(prefix='/api/connection-test', tags=['connection-test'])

@status_check_router.get('/soccerway')
def run_soccerway_test_connection():
    url = "https://ru.soccerway.com/a/block_h2h_matches?block_id=page_match_1_block_h2hsection_head2head_7_block_h2h_matches_1&action=changePage&callback_params=%7B%22page%22%3A+-1%2C+%22block_service_id%22%3A+%22match_h2h_comparison_block_h2hmatches%22%2C+%22team_A_id%22%3A+43000%2C+%22team_B_id%22%3A+43009%7D&params=%7B%22page%22%3A+0%7D"

    response = requests.get(url, headers={
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:102.0) Gecko/20100101 Firefox/131.0'})

    # Проверка на успешный ответ
    if response.status_code == 200:
        return {"status": response.status_code, }
    else:
        return {"status": "error", "status_code": response.status_code, "message": response.text}


@status_check_router.get('/marafon')
def run_marafon_test_connection():
    url = "https://www.marathonbet.ru/su/betting/Football+-+11?page=0&pageAction=getPage"

    response = requests.get(url, headers={
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36'})

    if response.status_code == 200:
        return {"status": response.status_code, }
    else:
        return {"status": "error", "status_code": response.status_code, "message": response.text}


