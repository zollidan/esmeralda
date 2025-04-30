import urllib
import xlwt
from bs4 import BeautifulSoup

import re
import json
from art import tprint
import os
import sys
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

from date_functions import *
from request_functions import *

from exel_functions import *
from parser_functions import *
from constants import *
import requests



def main(date_start, date_end):

    tprint("soccerway-1")

    # Открываем файл Exel на запись
    wb = xlwt.Workbook()
    ws = wb.add_sheet('Big Data')

    # Создаём шапку в Exel
    make_header(ws)

    # # Получаем предел дат от пользователя
    # date_limit = get_date_limit()
    # date_start = date_limit[0]
    # date_end = date_limit[1]

    # Список дат в указанном пределе, с разницой в день
    date_list = make_date_list(date_start, date_end)

    # Получить список матчей за все даты сразу
    all_matches = []
    for date_str in date_list:
        try:
            response_url = BASE_LINK + date_str.replace('-', '/')
            response = get_response(response_url)
            matches_from_page = get_matches(response, date_str)
            all_matches.extend(matches_from_page)
        except:
            if DEBUG_MODE:
                print("[!] Main Request Loop Error [!]\n[args]:", BASE_LINK, date_str)

    if DEBUG_MODE:
        print(all_matches, "\n[ All Matches Info ]")

    # Старт парсинга
    try:
        match_number = 1
        for match_info in all_matches:
            try:

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Шапка № 1
                home_team_name = match_info.get('home_team_name')
                away_team_name = match_info.get('away_team_name')
                home_team_id = match_info.get('home_team_id')
                away_team_id = match_info.get('away_team_id')
                match_date = match_info.get('match_date')

                # Запись в Exel
                ws.write(match_number, 0, int(match_date.strftime('%d')))
                ws.write(match_number, 1, int(match_date.strftime('%m')))
                ws.write(match_number, 2, int(match_date.strftime('%Y')))
                ws.write(match_number, 3, match_date.strftime('%H:%M:%S'))
                ws.write(match_number, 4, home_team_name)
                ws.write(match_number, 5, away_team_name)
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # Получить 25 матчей друг против друга, для каждой категории 'home', 'away'
                head_to_head = get_h2h_list(match_info, need=25)
                # head_to_head = {
                #     'home_team_name': home_team_name,
                #     'away_team_name': away_team_name,
                #     'home_team_id': home_team_id,
                #     'away_team_id': away_team_name,
                #     'date': match_date, # datetime object
                #     'home': [],   # Домашние матчи
                #     'away': [],   # Выездные
                #     'all':  [],   # Все общие матчи
                # }

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Домашние h2h матчи [15]
                max_len = min(len(head_to_head['home']), 15)
                home_h2h_15 = head_to_head['home'][0:max_len]

                # Получение статистики
                stats_home_h2h_15 = get_match_stats(home_h2h_15, home_team_id, max_len)

                # Запись в Exel
                ws.write(match_number, 6, max_len)
                ws.write(match_number, 7, stats_home_h2h_15['win'])
                ws.write(match_number, 8, stats_home_h2h_15['draw'])
                ws.write(match_number, 9, stats_home_h2h_15['lose'])
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                #
                #
                #
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Получить 25 последних ДОМАШНИХ матчей
                home_last_25 = get_matches_with_filter(home_team_id, match_date, 25, filter="home")
                max_len = min(len(home_last_25), 25)
                print("------------")
                print(home_last_25)
                # Получение статистики
                stats_home_last_25 = get_match_stats(home_last_25, home_team_id, max_len)

                # Запись в Exel
                ws.write(match_number, 10, max_len)
                ws.write(match_number, 11, stats_home_last_25['win'])
                ws.write(match_number, 12, stats_home_last_25['draw'])
                ws.write(match_number, 13, stats_home_last_25['lose'])
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Получить 25 последних ВЫЕЗДНЫХ матчей ПРОТИВНИКОВ
                away_last_25 = get_matches_with_filter(away_team_id, match_date, 25, filter="away")
                max_len = min(len(away_last_25), 25)

                # Получение статистики
                stats_away_last_25 = get_match_stats(away_last_25, away_team_id, max_len)

                # Запись в Exel
                ws.write(match_number, 14, max_len)
                ws.write(match_number, 15, stats_away_last_25['lose'])
                ws.write(match_number, 16, stats_away_last_25['draw'])
                ws.write(match_number, 17, stats_away_last_25['win'])
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Домашние h2h матчи [25]
                max_len = min(len(head_to_head['home']), 25)
                home_h2h_25 = head_to_head['home'][0:max_len]

                # Получение статистики
                stats_home_h2h_25 = get_match_stats(home_h2h_25, home_team_id, max_len)

                # Запись в Exel
                ws.write(match_number, 18, max_len)
                ws.write(match_number, 19, stats_home_h2h_15['less_2_5'])
                ws.write(match_number, 20, stats_home_h2h_15['more_2_5'])
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Домашние матчи [25]
                max_len = min(len(home_last_25), 25)

                # Запись в Exel из готовой статистики
                ws.write(match_number, 21, max_len)
                ws.write(match_number, 22, stats_home_last_25['less_2_5'])
                ws.write(match_number, 23, stats_home_last_25['more_2_5'])
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Выездные матчи ПРОТИВНИКА [25]
                max_len = min(len(away_last_25), 25)

                # Запись в Exel
                ws.write(match_number, 24, max_len)
                ws.write(match_number, 25, stats_away_last_25['less_2_5'])
                ws.write(match_number, 26, stats_away_last_25['more_2_5'])
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Шапка № 2
                link = match_info.get('link')
                league_name = match_info.get('league_name')

                # Запись в Exel
                ws.write(match_number, 27, link)
                ws.write(match_number, 28, league_name)
                #
                ws.write(match_number, 29, int(match_date.strftime('%d')))
                ws.write(match_number, 30, int(match_date.strftime('%m')))
                ws.write(match_number, 31, int(match_date.strftime('%Y')))
                ws.write(match_number, 32, match_date.strftime('%H:%M:%S'))
                ws.write(match_number, 33, home_team_name)
                ws.write(match_number, 34, away_team_name)
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Все h2h матчи [25]
                max_len = min(len(head_to_head['all']), 25)
                all_h2h_25 = head_to_head['all'][0:max_len]

                # Получение статистики
                stats_all_h2h_25 = get_match_stats(all_h2h_25, home_team_id, max_len)
                goals_all = stats_all_h2h_25['goals_all']

                # Сумма голов [25]
                summ_goals_all_25 = get_summ_from_list(goals_all, 25)

                # Сумма голов [3]
                summ_goals_all_3 = get_summ_from_list(goals_all, 3)

                # Сумма голов [5]
                summ_goals_all_5 = get_summ_from_list(goals_all, 5)

                # Запись в Exel
                max_len = min(len(all_h2h_25), 25)
                ws.write(match_number, 35, max_len)
                ws.write(match_number, 36, summ_goals_all_25)
                #
                max_len = min(len(all_h2h_25), 3)
                ws.write(match_number, 37, max_len)
                ws.write(match_number, 38, summ_goals_all_3)
                #
                max_len = min(len(all_h2h_25), 5)
                ws.write(match_number, 39, max_len)
                ws.write(match_number, 40, summ_goals_all_5)
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Все матчи Команды А [25]
                all_last_25 = get_matches_with_filter(home_team_id, match_date, 25, filter="all")
                max_len = min(len(all_last_25), 25)

                # Получение статистики
                stats_all_last_25 = get_match_stats(all_last_25, home_team_id, max_len)
                goals_home = stats_all_last_25['goals_all']

                # Сумма голов команды А [25]
                summ_goals_home_25 = get_summ_from_list(goals_home, 25)

                # Сумма голов [3]
                summ_goals_home_3 = get_summ_from_list(goals_home, 3)

                # Сумма голов [5]
                summ_goals_home_5 = get_summ_from_list(goals_home, 5)

                # Запись в Exel
                max_len = min(len(all_last_25), 25)
                ws.write(match_number, 41, max_len)
                ws.write(match_number, 42, summ_goals_home_25)
                #
                max_len = min(len(all_last_25), 3)
                ws.write(match_number, 43, max_len)
                ws.write(match_number, 44, summ_goals_home_3)
                #
                max_len = min(len(all_last_25), 5)
                ws.write(match_number, 45, max_len)
                ws.write(match_number, 46, summ_goals_home_5)
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Все матчи Команды Б [25]
                away_all_last_25 = get_matches_with_filter(away_team_id, match_date, 25, filter="all")
                max_len = min(len(away_all_last_25), 25)

                # Получение статистики
                stats_away_all_last_25 = get_match_stats(away_all_last_25, away_team_id, max_len)
                goals_home = stats_away_all_last_25['goals_all']

                # Сумма голов команды Б [25]
                summ_goals_home_25 = get_summ_from_list(goals_home, 25)

                # Сумма голов [3]
                summ_goals_home_3 = get_summ_from_list(goals_home, 3)

                # Сумма голов [5]
                summ_goals_home_5 = get_summ_from_list(goals_home, 5)

                # Запись в Exel
                max_len = min(len(away_all_last_25), 25)
                ws.write(match_number, 47, max_len)
                ws.write(match_number, 48, summ_goals_home_25)
                #
                max_len = min(len(away_all_last_25), 3)
                ws.write(match_number, 49, max_len)
                ws.write(match_number, 50, summ_goals_home_3)
                #
                max_len = min(len(away_all_last_25), 5)
                ws.write(match_number, 51, max_len)
                ws.write(match_number, 52, summ_goals_home_5)
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Концовка | Запись в Exel
                ws.write(match_number, 53, link)
                ws.write(match_number, 54, league_name)
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                if DEBUG_MODE:
                    print(" ************************* ")
                    print(" Home Team: ", home_team_name)
                    print(" Away Team: ", away_team_name)
                    print(" ************************* ")

                    print("Link:", match_info['link'])

                    print("Home Last (25)", len(away_last_25))
                    print(away_last_25)
                    print("Stats:", stats_away_last_25)

                    # for debug_tmp_elem in away_last_25:
                    #     print(debug_tmp_elem['score_home'], ":", debug_tmp_elem['score_away'], "\t", debug_tmp_elem['match_date'].strftime("%d/%m/%y"))

                # В случаи ошибки, итерация не выполнится
                match_number += 1
            except:
                if DEBUG_MODE:
                    print('Match Parsing Error:', match_info)
    except:
        if DEBUG_MODE:
            print("*** Main Parsing Loop Error ***")
    finally:

        # Надпись о завершении работы
        if QUITE_MODE:
            print(WORK_IS_DONE_MESSAGE)

        # Сохраняем файл
        name = EXEL_NAME
        k = 1
        while True:
            try:
                wb.save(name + '.xls')
                print("[ Файл сохранён:", name + ".xls ]")
                break
            except:
                if k >= MAX_COUNT_REPEAT_EXEL:
                    print("!!! Файл не сохранён !!!")
                    break
                name += str(k)
                k += 1

        try:

            url = "http://web:8000/api/files/upload"
            files = {'file': (name + '.xls', open(name + '.xls', 'rb'), 'application/vnd.ms-excel')}
            response = requests.post(url, files=files)
                
            if response.status_code == 200:
                print("[ Файл успешно отправлен на сервер FastAPI ]")
            else:
                print("[ Ошибка при отправке файла на сервер FastAPI ]", response.status_code, response.text)
        except Exception as e:
            print("[ Ошибка при отправке файла на сервер FastAPI ]", str(e))
    