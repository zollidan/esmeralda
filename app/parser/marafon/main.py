import os

import pandas as pd

from app.dao import recored_and_upload_file
from app.parser.marafon.parser_functions import *

headers = {
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36"
}


def run_marafon_parser(file_name):
    main_site_url = 'https://www.marathonbet.ru/su/betting/Football+-+11'

    main_data = []

    for page in tqdm(range(find_number_of_pages())):

        payload = {'page': page, 'pageAction': 'getPage'}
        req = requests.get(main_site_url, params=payload, headers=headers)

        soup = BeautifulSoup(req.json()[0]['content'], 'html.parser')

        league_block = soup.find_all("div", {"class": "category-container"})

        for league in league_block:
            league_name = league.find("a", {"class": "category-label-link"}).text
            games = league.find_all("div", {"class": "bg"})

            for game in games:
                game_data = find_game_data_and_write(game, league_name)
                main_data.append(game_data)

    columns = ['день', 'месяц', 'год', 'время','команда 1', 'команда 2', 'лига', 'П1', 'X', 'П2', '1X', '12', 'X2', 'фора 1', 'фора 1', 'фора 2', 'фора 2', 'тотал', 'меньше', 'тотал', 'больше', 'link']

    df = pd.DataFrame(main_data, columns=columns)

    df = find_correct_price(df)

    df.to_excel(file_name, index=False)

    recored_and_upload_file(file_name)

    os.remove(file_name)
