import pandas as pd
import requests
from bs4 import BeautifulSoup
from parser_functions import find_game_data_and_write, find_number_of_pages, find_correct_price
from tqdm import tqdm
from art import tprint

main_site_url = 'https://www.marathonbet.ru/su/betting/Football+-+11'

main_data = []

tprint("MARAFON")
print("Работа начинается.")

for page in tqdm(range(find_number_of_pages())):

    payload = {'page': page, 'pageAction': 'getPage'}
    req = requests.get(main_site_url, params=payload)

    soup = BeautifulSoup(req.json()[0]['content'], 'html.parser')

    league_block = soup.find_all("div", {"class": "category-container"})

    for league in league_block:
        league_name = league.find("a", {"class": "category-label-link"}).text
        games = league.find_all("div", {"class": "bg"})

        for game in games:
            game_data = find_game_data_and_write(game, league_name)
            main_data.append(game_data)

#columns = ['день', 'месяц', 'год', 'команда 1', 'команда 2', 'лига', '1', '2', '3', '4', '5', '6', '7', '8', '9', '10', '11', '12', '13', '14', 'link', 'коэффициент', 'less', 'more']

df = pd.DataFrame(main_data)

df = find_correct_price(df)

df.to_excel('marafon.xlsx', index=False)