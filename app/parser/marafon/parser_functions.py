import re
from datetime import datetime

import requests
from bs4 import BeautifulSoup
from tqdm import tqdm

headers = {
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36"
}

def write_prices_to_list(game_data, prices):
    for price_number in range(len(prices)):

        price = prices[price_number].text.replace(" ", "")

        matches = re.findall(r'\((.*?)\)|([0-9.]+)', price)

        if '—' in price:
            game_data.append(price)

        for match in matches:
            game_data.append(match[0] if match[0] else match[1])

    return game_data

months = {
    'янв': '01', 'фев': '02', 'мар': '03', 'апр': '04', 'май': '05', 'июн': '06',
    'июл': '07', 'авг': '08', 'сен': '09', 'окт': '10', 'ноя': '11', 'дек': '12'
}

def find_game_data_and_write(game, league_name):
    game_data = []
    team1 = game.find_all("a", {"class": "member-link"})[0].text.strip()
    team2 = game.find_all("a", {"class": "member-link"})[1].text.strip()
    date = game.find_all("div", {"class": "date"})[0].text.strip().split(" ")
    game_link = game.find_all("a", {"class": "member-link"})[0]['href']
    empty_prices = game.find_all("span", string=lambda text: '—' in text)
    prices = game.find_all("td", {"class": "price"})

    if empty_prices:
        prices = empty_prices + prices

    if len(date) == 3:
        game_data.append(date[0])
        game_data.append(months.get(date[1], date[1]))
        game_data.append(datetime.now().year)
        game_data.append(date[2])
    else:
        game_data.append(datetime.now().day)
        game_data.append(datetime.now().month)
        game_data.append(datetime.now().year)
        game_data.append(date[0])
    game_data.append(team1)
    game_data.append(team2)
    game_data.append(league_name)
    write_prices_to_list(game_data, prices)
    game_data.append("https://www.marathonbet.ru"+game_link)

    return game_data

def find_number_of_pages():
    main_site_url = 'https://www.marathonbet.ru/su/betting/Football+-+11'
    page_counter = 0
    next_page = True
    while next_page:
        payload = {'page': page_counter, 'pageAction': 'getPage'}
        req = requests.get(main_site_url, params=payload, headers=headers)
        print(req.status_code)
        print(req.text)
        if req.json()[1]['val']:
            page_counter += 1
        else:
            return page_counter

def find_correct_price(df):
    for index, row in tqdm(df.iterrows()):
        value = row[19]
        game_url = row[21]
        if value and float(value) != 2.5:
            new_prices_list = []

            req = requests.get(game_url, headers=headers)

            soup = BeautifulSoup(req.text, 'html.parser')

            total_block = soup.find("div", {"data-block-type-id": "3"}).find("table", {"class": "td-border"})
            new_prices = total_block.find_all("td", {"class": "price"})

            for new_price in new_prices:
                new_price_split = new_price.text.replace("\n", "").strip().split(" ")
                if new_price_split[0] == '(2.5)':
                    new_prices_list.append(new_price_split[1])

            # Обрабатываем только строки, для которых нашли значения
            if new_prices_list:
                df.at[index, 'коэффициент'] = '2.5'
                df.at[index, 'less'] = new_prices_list[0]
                df.at[index, 'more'] = new_prices_list[1]

        else:
            continue
    return df