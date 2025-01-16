import os
from app.parser.soccerway.date_functions import *
from app.parser.soccerway.exel_functions import *
from app.parser.soccerway.parser_functions import *
from app.config import settings, s3_client


def run_soccerway(user_date, my_file_name):
    
    manager = ExcelManager(filename=my_file_name)

    # Получаем предел дат от пользователя
    # date_limit = get_date_limit()
    # date_start, date_end = date_limit
    
    """
    тут надо пофиксить и сделать чтобы можно было диапазон а не одну дату выбирать
    """
    
    
    date_start = user_date
    date_end = user_date


    # Список дат в указанном пределе, с разницой в день
    date_list = make_date_list(date_start, date_end)

    print(date_list, "\n[ Date List ]")
    
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

    try:
        print("[ Parse Matches ]")
        for match_info in tqdm(all_matches):
            try:

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Шапка
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                home_team_name = match_info.get('home_team_name')
                away_team_name = match_info.get('away_team_name')
                home_team_id = match_info.get('home_team_id')
                away_team_id = match_info.get('away_team_id')
                match_date = match_info.get('match_date')
                # # # Запись в Exel # # #
                manager.write(int(match_date.strftime('%d')))
                manager.write(int(match_date.strftime('%m')))
                manager.write(int(match_date.strftime('%Y')))
                manager.write(match_date.strftime('%H:%M:%S'))
                manager.write(home_team_name)
                manager.write(away_team_name)
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
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
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Получить 100 матчей друг против друга, для каждой категории 'home', 'away'
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                head_to_head = get_h2h_list(match_info, need=100)

                # Домашние h2h матчи команды А [100]
                max_len = min(len(head_to_head['home']), 100)
                home_h2h_25 = head_to_head['home'][0:max_len]

                # Получение статистики
                stats_home_h2h_25 = get_match_stats(home_h2h_25, home_team_id, max_len)
                goals_all = stats_home_h2h_25['goals_all']

                # Сумма всех голов за [25], [3], [5] матчей
                summ_goals_all_25 = get_summ_from_list(goals_all, 100)
                summ_goals_all_3 = get_summ_from_list(goals_all, 10)
                summ_goals_all_5 = get_summ_from_list(goals_all, 5)

                # # # Запись в Exel # # #
                # Голов за [25] матчей
                max_len = min(len(home_h2h_25), 100)
                manager.write(max_len)
                manager.write(summ_goals_all_25)
                # Голов за [3] матча
                max_len = min(len(home_h2h_25), 10)
                manager.write(max_len)
                manager.write(summ_goals_all_3)
                # Голов за [5] матчей
                max_len = min(len(home_h2h_25), 5)
                manager.write(max_len)
                manager.write(summ_goals_all_5)
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Последние домашние матчи Команды А [25]
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                home_last_25 = get_matches_with_filter(home_team_id, match_date, 100, filter="home")
                max_len = min(len(home_last_25), 100)

                # Получение статистики
                stats_home_last_25 = get_match_stats(home_last_25, home_team_id, max_len)
                goals_all = stats_home_last_25['goals_all']

                # Сумма всех голов за [25], [3], [5] матчей
                summ_goals_all_25 = get_summ_from_list(goals_all, 100)
                summ_goals_all_3 = get_summ_from_list(goals_all, 10)
                summ_goals_all_5 = get_summ_from_list(goals_all, 5)

                # # # Запись в Exel # # #
                # Голов за [25] матчей
                max_len = min(len(home_last_25), 100)
                manager.write(max_len)
                manager.write(summ_goals_all_25)
                # Голов за [3] матча
                max_len = min(len(home_last_25), 10)
                manager.write(max_len)
                manager.write(summ_goals_all_3)
                # Голов за [5] матчей
                max_len = min(len(home_last_25), 5)
                manager.write(max_len)
                manager.write(summ_goals_all_5)
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Последние выездные матчи Команды Б [25]
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                away_last_25 = get_matches_with_filter(away_team_id, match_date, 100, filter="away")
                max_len = min(len(away_last_25), 100)

                # Получение статистики
                stats_away_last_25 = get_match_stats(away_last_25, away_team_id, max_len)
                goals_all = stats_away_last_25['goals_all']

                # Сумма всех голов за [25], [3], [5] матчей
                summ_goals_all_25 = get_summ_from_list(goals_all, 100)
                summ_goals_all_3 = get_summ_from_list(goals_all, 10)
                summ_goals_all_5 = get_summ_from_list(goals_all, 5)

                # # # Запись в Exel # # #
                # Голов за [25] матчей
                max_len = min(len(away_last_25), 100)
                manager.write(max_len)
                manager.write(summ_goals_all_25)
                # Голов за [3] матча
                max_len = min(len(away_last_25), 10)
                manager.write(max_len)
                manager.write(summ_goals_all_3)
                # Голов за [5] матчей
                max_len = min(len(away_last_25), 5)
                manager.write(max_len)
                manager.write(summ_goals_all_5)
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Концовка
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                link = match_info.get('link')
                league_name = match_info.get('league_name')
                # # # Запись в Exel # # #
                manager.write(link)
                manager.write(league_name)
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Шапка № 1
                home_team_name = match_info.get('home_team_name')
                away_team_name = match_info.get('away_team_name')
                home_team_id = match_info.get('home_team_id')
                away_team_id = match_info.get('away_team_id')
                match_date = match_info.get('match_date')

                # Запись в Exel
                manager.write(int(match_date.strftime('%d')))
                manager.write(int(match_date.strftime('%m')))
                manager.write(int(match_date.strftime('%Y')))
                manager.write(match_date.strftime('%H:%M:%S'))
                manager.write(home_team_name)
                manager.write(away_team_name)
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # Получить 25 матчей друг против друга, для каждой категории 'home', 'away'
                head_to_head = get_h2h_list(match_info, need=100)
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
                manager.write(max_len)
                manager.write(stats_home_h2h_15['win'])
                manager.write(stats_home_h2h_15['draw'])
                manager.write(stats_home_h2h_15['lose'])
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
                manager.write(max_len)
                manager.write(stats_home_last_25['win'])
                manager.write(stats_home_last_25['draw'])
                manager.write(stats_home_last_25['lose'])
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Получить 25 последних ВЫЕЗДНЫХ матчей ПРОТИВНИКОВ
                away_last_25 = get_matches_with_filter(away_team_id, match_date, 100, filter="away")
                max_len = min(len(away_last_25), 100)

                # Получение статистики
                stats_away_last_25 = get_match_stats(away_last_25, away_team_id, max_len)

                # Запись в Exel
                manager.write(max_len)
                manager.write(stats_away_last_25['lose'])
                manager.write(stats_away_last_25['draw'])
                manager.write(stats_away_last_25['win'])
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Домашние h2h матчи [25]
                max_len = min(len(head_to_head['home']), 100)
                home_h2h_25 = head_to_head['home'][0:max_len]

                # Получение статистики
                stats_home_h2h_25 = get_match_stats(home_h2h_25, home_team_id, max_len)

                # Запись в Exel
                manager.write(max_len)
                manager.write(stats_home_h2h_15['less_2_5'])
                manager.write(stats_home_h2h_15['more_2_5'])
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Домашние матчи [25]
                max_len = min(len(home_last_25), 100)

                # Запись в Exel из готовой статистики
                manager.write(max_len)
                manager.write(stats_home_last_25['less_2_5'])
                manager.write(stats_home_last_25['more_2_5'])
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Выездные матчи ПРОТИВНИКА [25]
                max_len = min(len(away_last_25), 100)

                # Запись в Exel
                manager.write(max_len)
                manager.write(stats_away_last_25['less_2_5'])
                manager.write(stats_away_last_25['more_2_5'])
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Шапка № 2
                link = match_info.get('link')
                league_name = match_info.get('league_name')

                # Запись в Exel
                manager.write(link)
                manager.write(league_name)
                #
                manager.write(int(match_date.strftime('%d')))
                manager.write(int(match_date.strftime('%m')))
                manager.write(int(match_date.strftime('%Y')))
                manager.write(match_date.strftime('%H:%M:%S'))
                manager.write(home_team_name)
                manager.write(away_team_name)
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Все h2h матчи [25]
                max_len = min(len(head_to_head['all']), 100)
                all_h2h_25 = head_to_head['all'][0:max_len]

                # Получение статистики
                stats_all_h2h_25 = get_match_stats(all_h2h_25, home_team_id, max_len)
                goals_all = stats_all_h2h_25['goals_all']

                # Сумма голов [25]
                summ_goals_all_25 = get_summ_from_list(goals_all, 100)

                # Сумма голов [3]
                summ_goals_all_3 = get_summ_from_list(goals_all, 10)

                # Сумма голов [5]
                summ_goals_all_5 = get_summ_from_list(goals_all, 5)

                # Запись в Exel
                max_len = min(len(all_h2h_25), 100)
                manager.write(max_len)
                manager.write(summ_goals_all_25)
                #
                max_len = min(len(all_h2h_25), 10)
                manager.write(max_len)
                manager.write(summ_goals_all_3)
                #
                max_len = min(len(all_h2h_25), 5)
                manager.write(max_len)
                manager.write(summ_goals_all_5)
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Все матчи Команды А [25]
                all_last_25 = get_matches_with_filter(home_team_id, match_date, 100, filter="all")
                max_len = min(len(all_last_25), 100)

                # Получение статистики
                stats_all_last_25 = get_match_stats(all_last_25, home_team_id, max_len)
                goals_home = stats_all_last_25['goals_all']

                # Сумма голов команды А [25]
                summ_goals_home_25 = get_summ_from_list(goals_home, 100)

                # Сумма голов [3]
                summ_goals_home_3 = get_summ_from_list(goals_home, 10)

                # Сумма голов [5]
                summ_goals_home_5 = get_summ_from_list(goals_home, 5)

                # Запись в Exel
                max_len = min(len(all_last_25), 100)
                manager.write(max_len)
                manager.write(summ_goals_home_25)
                #
                max_len = min(len(all_last_25), 10)
                manager.write(max_len)
                manager.write(summ_goals_home_3)
                #
                max_len = min(len(all_last_25), 5)
                manager.write(max_len)
                manager.write(summ_goals_home_5)
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
                # Все матчи Команды Б [25]
                away_all_last_25 = get_matches_with_filter(away_team_id, match_date, 100, filter="all")
                max_len = min(len(away_all_last_25), 100)

                # Получение статистики
                stats_away_all_last_25 = get_match_stats(away_all_last_25, away_team_id, max_len)
                goals_home = stats_away_all_last_25['goals_all']

                # Сумма голов команды Б [25]
                summ_goals_home_25 = get_summ_from_list(goals_home, 100)

                # Сумма голов [3]
                summ_goals_home_3 = get_summ_from_list(goals_home, 10)

                # Сумма голов [5]
                summ_goals_home_5 = get_summ_from_list(goals_home, 5)

                # Запись в Exel
                max_len = min(len(away_all_last_25), 100)
                manager.write(max_len)
                manager.write(summ_goals_home_25)
                #
                max_len = min(len(away_all_last_25), 10)
                manager.write(max_len)
                manager.write(summ_goals_home_3)
                #
                max_len = min(len(away_all_last_25), 5)
                manager.write(max_len)
                manager.write(summ_goals_home_5)
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # ы# # # # # # # # #
                # Концовка | Запись в Exel
                manager.write(link)
                manager.write(league_name)
                # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

                # В случаи ошибки, итерация не выполнится и мы перезапишем тек. строку
                manager.next_row()


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

        # Сохранение файла
        manager.save()
        
        s3_client.upload_file(my_file_name, settings.AWS_BUCKET, my_file_name)

        os.remove(my_file_name)

        # input("Press Enter to continue...")
