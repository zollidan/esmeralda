import pandas as pd
import time, os, re, logging, sys
from datetime import datetime
from bs4 import BeautifulSoup
from playwright.sync_api import sync_playwright, TimeoutError as PlaywrightTimeoutError

def main(input_date:str) -> pd.DataFrame:

    # Отключаем логирование
    logging.getLogger('playwright').setLevel(logging.WARNING)
    logging.getLogger('urllib3').setLevel(logging.WARNING)

    # input_date = input('Введите дату в формате 2025-11-27:  ')

    # зарефакторить это кал
    url = f'https://m.sofascore.com/football/{input_date}'


    # зачем тут слайс
    date1 = url.split('/')[-1]
    date1_obj = datetime.strptime(date1, "%Y-%m-%d")

    def base_cleanup(page):
        """Базовая очистка для парсинга"""
        # print("🧹 Выполняю базовую очистку...")
        
        try:
            # 1. Очищаем cookies
            page.context.clear_cookies()
        except:
            pass
        
        try:
            # 2. Очищаем локальное хранилище
            page.evaluate("""
                () => {
                    try {
                        localStorage.clear();
                        sessionStorage.clear();
                    } catch(e) {}
                }
            """)
        except:
            pass
        
        # print("✓ Очистка завершена")
        time.sleep(0.5)




    def get_stealth_driver_chrome(opt):
        """Полный аналог вашей функции get_stealth_driver_chrome на Playwright"""
        playwright = sync_playwright().start()
        
        args = [
            "--no-sandbox",
            "--disable-dev-shm-usage",
            "--disable-gpu",
            "--ignore-certificate-errors",
            "--enable-unsafe-swiftshader",
            "--disable-popup-blocking",
            "--disable-notifications",
            "--disable-infobars",
            "--disable-extensions",
            "--disable-web-security",
            "--no-first-run",
            "--no-default-browser-check",
            "--disable-component-extensions-with-background-pages",
            "--log-level=3",
            "--disable-logging",
            "--disable-blink-features=AutomationControlled",
            "--disable-images",
            opt if opt else "--start-maximized"
        ]
        
        browser = playwright.chromium.launch(
            headless=False,
            args=args
        )
        
        # Создаем контекст с настройками
        context = browser.new_context(
            viewport={'width': 1920, 'height': 1080},
            user_agent='Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
            java_script_enabled=True,
            bypass_csp=True
        )
        
        # Блокируем ненужные ресурсы для ускорения
        def block_resources(route):
            if route.request.resource_type in ["image", "stylesheet", "font", "media"]:
                route.abort()
            else:
                route.continue_()
        
        # context.route("**/*", block_resources)
        
        page = context.new_page()
        
        # Добавляем stealth скрипты
        page.add_init_script("""
            Object.defineProperty(navigator, 'webdriver', { get: () => undefined });
            Object.defineProperty(navigator, 'plugins', { get: () => [1, 2, 3, 4, 5] });
            Object.defineProperty(navigator, 'languages', { get: () => ['en-US', 'en'] });
        """)
        
        

        # Возвращаем и playwright объект тоже, чтобы потом закрыть
        return page, browser, playwright

    def get_stealth_driver_firefox(opt=""):
        """Полный аналог вашей функции get_stealth_driver_firefox на Playwright"""
        playwright = sync_playwright().start()
        
        args = [
            "--no-sandbox",
            "--disable-dev-shm-usage",
            "--disable-gpu"
        ]
        
        if opt:
            args.append(opt)
        
        browser = playwright.firefox.launch(
            headless=False,
            args=args
        )
        
        context = browser.new_context(
            viewport={'width': 1920, 'height': 1080},
            user_agent='Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0',
            java_script_enabled=True,
            bypass_csp=True
        )
        
        # Блокируем ресурсы
        def block_resources(route):
            if route.request.resource_type in ["image", "stylesheet", "font", "media"]:
                route.abort()
            else:
                route.continue_()
        
        context.route("**/*", block_resources)
        
        page = context.new_page()
        
        # Добавляем stealth скрипты
        page.add_init_script("""
            Object.defineProperty(navigator, 'webdriver', { get: () => undefined });
            Object.defineProperty(navigator, 'plugins', { get: () => [1, 2, 3, 4, 5] });
        """)
        
        return page, browser, playwright

    def check_and_refresh_full_text(page):
        """Проверяет наличие текста про Favourites и обновляет страницу"""
        try:
            full_text = "Add to Favourites to keep track of upcoming events. You can adjust this and notifications later in the Favourites tab."
            
            # Ищем элемент с текстом
            element = page.query_selector(f"text={full_text}")
            
            if element:
                print("Найден span с текстом про Favourites, делаю refresh...")
                page.reload(wait_until="domcontentloaded", timeout=5000)
                return True
        except Exception as e:
            return False
        return False

    def get_file_path():
        """Получает путь к data.xlsx"""
        if getattr(sys, 'frozen', False):
            exe_dir = os.path.dirname(sys.executable)
            parent_dir = os.path.dirname(exe_dir)
            data_path = os.path.join(parent_dir, "data.xlsx")
        else:
            current_dir = os.path.abspath(".")
            data_path = os.path.join(current_dir, "data.xlsx")
        
        return data_path

    # ============= ОСНОВНОЙ КОД =============

    # Инициализация драйвера для firefox (как в вашем коде)
    opt = "--force-device-scale-factor=1"
    page, browser, playwright = get_stealth_driver_chrome(opt)

    try:
        # Переходим на страницу
        page.goto(url, wait_until="domcontentloaded", timeout=10000)
        
        # Проверяем попап "Help us improve"
        try:
            # Ищем div с текстом "Help us improve"
            div_xpath = '//div[.//span[contains(text(), "Help us improve")]]'
            div_element = page.wait_for_selector(f'xpath={div_xpath}', timeout=5000)
            
            if div_element:
                # Находим все кнопки внутри этого div
                buttons = div_element.query_selector_all('button')
                
                # Кликаем по первой кнопке (крестик)
                if buttons and len(buttons) > 0:
                    buttons[0].click()
                    print("Первая кнопка (крестик) нажата")
                    page.wait_for_timeout(1000)
        except Exception as e:
            pass  # Элемент не найден - это нормально
        
        # Ищем кнопку "Show all"
        try:
            button_xpath = '//div[contains(@class, "mdDown:pt_sm")]//button[contains(text(), "Show all")]'
            button = page.wait_for_selector(f'xpath={button_xpath}', timeout=5000)
            
            if button:
                button.click()
                print("Нажали Show all")
                page.evaluate("window.scrollTo(0, 0)")
                # page.evaluate("document.body.style.zoom = '0.3'")
                page.wait_for_timeout(10000)  # Ждем 10 секунд
        except Exception as e:
            print(f"Кнопка Show all не найдена: {e}")
        
        # Сохраняем исходный код страницы (аналог driver.page_source)
        sofascore_page = page.content()
        





        page.evaluate("document.body.style.zoom = '0.3'")
        sofascore_page = page.content()
        list_matches = set()

        # Скроллим небольшими шагами
        scroll_step = 500  # Всего 500 пикселей за раз (меньше чем было)
        current_scroll = 0

        list_matches = []  

        # Скроллим 20 раз
        for i in range(30):
            # Собираем матчи
            elements = page.query_selector_all('a[href*="/match/"]')
            for element in elements:
                href = element.get_attribute('href')
                if href and '/match/' in href:
                    if href.startswith('/'):
                        href = f"https://www.sofascore.com{href}"
                    # Проверяем, нет ли уже такой ссылки
                    if href not in list_matches:
                        list_matches.append(href)  # Добавляем в конец
            
            # Обновляем код страницы
            sofascore_page = page.content()
            
            print(f"Скролл {i+1}: собрано {len(list_matches)} матчей")
            
            # Скроллим небольшим шагом
            current_scroll += scroll_step
            page.evaluate(f"window.scrollTo(0, {current_scroll})")
            
            # Ждем немного
            page.wait_for_timeout(500)

        print(f"Итого после 20 скроллов: {len(list_matches)} матчей")








    except PlaywrightTimeoutError:
        print("Таймаут при загрузке страницы...")
        page.evaluate("window.stop();")
        page.wait_for_timeout(1000)
    except Exception as e:
        print(f"Ошибка: {e}")

    # Закрываем Firefox и переключаемся на Chrome (как в вашем коде)
    # browser.close()
    # playwright.stop()

    if len(list_matches) > 0:
        matches_results_list = []
        
        full_path = get_file_path()
        
        # Чтение существующего файла
        if os.path.exists(full_path):
            try:
                df = pd.read_excel(full_path)
                column_b_values = df['href'].tolist()
            except:
                column_b_values = []
        else:
            column_b_values = []
        
        count = 0
        
        # Инициализация драйвера для Chrome с новым масштабом
        # opt = "--force-device-scale-factor=1.0"
        # page, browser, playwright = get_stealth_driver_chrome(opt)
        
        for match_url in list_matches:
            match_url = 'https://www.sofascore.com/football/match/brentford-arsenal/Rsab#id:14025126'
            bad_connect = 0
            count += 1

            if 'football' not in match_url:
                print(f'ПРОПУСКАЕМ Матч {match_url} не относится к футболу')
                continue
            
            if count % 5 == 0:
                # print(f"\n⚠️  Достигнут цикл #{count} - выполняю очистку")
                base_cleanup(page)
            
            if match_url in column_b_values:
                print(f'Матч {match_url} уже обрабатывался')
                continue
            
            start_time = time.time()
            
            match_id = match_url.split('#id:')[-1] if '#id:' in match_url else ""
            matches_results_dict = {}
            
            try:
                page.goto(match_url, wait_until="domcontentloaded", timeout=5000)
            except PlaywrightTimeoutError:
                page.evaluate("window.stop();")
                page.wait_for_timeout(1000)
                print(f"Таймаут при загрузке {match_url}, останавливаем...")
            except Exception as e:
                print(f"Ошибка загрузки {match_url}: {e}")
                continue
            
            # Дополнительная проверка - ждем появления body
            try:
                page.wait_for_selector("body", timeout=5000)
            except:
                print("Основной контент не загрузился")
                continue
            
            # Поиск элемента с Favourite (ваша логика с 3 попытками)
            max_retries = 3
            retry_count = 0
            success = False
            data_list = []
            
            while retry_count < max_retries and not success:
                try:
                    # Ищем div с текстом Favourite
                    div_xpath = '//div[not(@class) and contains(., "Favourite")]'
                    div_element = page.query_selector(f'xpath={div_xpath}')
                    
                    if div_element:
                        # Находим все span внутри
                        spans = div_element.query_selector_all('span')
                        data_list = [span.text_content() for span in spans if span.text_content().strip()]
                        success = True
                    else:
                        retry_count += 1
                        if retry_count < max_retries:
                            print(f"Попытка {retry_count}/{max_retries} не удалась, перезагружаем страницу {match_url}")
                            page.reload(wait_until="domcontentloaded", timeout=5000)

                        else:
                            print(f"Не удалось найти элемент после {max_retries} попыток")
                except Exception as e:
                    retry_count += 1
                    if retry_count < max_retries:
                        print(f"Ошибка при поиске элемента, попытка {retry_count}/{max_retries}")
                        page.reload(wait_until="domcontentloaded", timeout=5000)

            
            if not success:
                print(f"Не удалось за {max_retries} попыток, пропускаем {match_url}")
                continue
            
            # Поиск даты и времени в найденных данных
            date_pattern = r'\b\d{2}/\d{2}/\d{4}\b'
            time_pattern = r'\b\d{2}:\d{2}\b'
            
            found_date = None
            found_time = None
            day = 0
            month = 0
            year = 0
            
            for text in data_list:
                text = text.strip()

                if not found_date and re.match(date_pattern, text):
                    found_date = text
                    try:
                        day, month, year = found_date.split('/')
                    except:
                        pass
                    
                if not found_time and re.match(time_pattern, text):
                    found_time = text
                
                if found_date and found_time:
                    break
            
            # Ищем и кликаем H2H кнопку
            try:
                h2h_button = page.wait_for_selector('button[data-testid="tab-matches"]', timeout=5000)
                if h2h_button:
                    h2h_button.click()
                    # print("Кликнут элемент H2H")
                    page.wait_for_timeout(1000)
            except Exception as e:
                print(f"Элемент H2H не найден: {e}")
            
            # Ищем кнопку "Показать больше" или "Show more"
            try:
                target_div = page.wait_for_selector('div[class="w_100%"]', timeout=5000)
                
                if target_div:
                    # Сначала ищем на русском
                    show_more_button = target_div.query_selector('button:has-text("Показать больше")')
                    
                    if not show_more_button:
                        # Если не нашли на русском, ищем на английском
                        show_more_button = target_div.query_selector('button:has-text("Show more")')
                    
                    if show_more_button:
                        show_more_button.click()
                        # print("Кликнули на кнопку 'Показать больше'")
                        page.wait_for_timeout(1000)
            except Exception as e:
                pass  # Кнопка может отсутствовать
            
            # Собираем все ссылки H2H
            try:
                main_div = page.query_selector('div[class="w_100%"]')
                
                if main_div:
                    match_links = main_div.query_selector_all('a[data-id]')
                    
                    # Инициализируем все переменные как в оригинальном коде
                    matches_data = []
                    games_h2h = 0
                    games_h2h_25 = 0
                    games_h2h_3 = 0
                    games_h2h_5 = 0
                    games_home_h2h = 0
                    h2h_matches_over_2_5 = 0
                    h2h_matches_under_2_5 = 0
                    balls_h2h = 0
                    balls_h2h_25 = 0
                    balls_h2h_3 = 0
                    balls_h2h_5 = 0
                    victories_on_the_home_field = 0
                    victories_on_the_guest_field = 0
                    draw = 0
                    draw_on_the_home_field = 0
                    loos_on_the_home_field = 0
                    сумма_мячей_в_очных_играх_на_поле_хозяев_25 = 0
                    сумма_мячей_в_очных_играх_на_поле_хозяев_5 = 0
                    сумма_мячей_в_очных_играх_на_поле_хозяев_3 = 0
                    количество_очных_игр_на_поле_хозяев_25 = 0
                    количество_очных_игр_на_поле_хозяев_5 = 0
                    количество_очных_игр_на_поле_хозяев_3 = 0
                    
                    # Получаем названия команд с основной страницы
                    all_bdis = page.query_selector_all('bdi')
                    command_name_1 = ""
                    command_name_2 = ""
                    
                    if len(all_bdis) >= 4:
                        command_name_1 = all_bdis[2].text_content().strip()
                        command_name_2 = all_bdis[3].text_content().strip()

                    # first_team_bdi = link.query_selector('xpath=.//div[1]/div[4]/div[1]/div[1]/div[1]//bdi')
                    # second_team_bdi = link.query_selector('xpath=.//div[1]/div[4]/div[1]/div[1]/div[2]//bdi')
                    
                    # if not first_team_bdi and not second_team_bdi:
                    # # if not first_team_bdi or not second_team_bdi:
                    #     continue
                        
                    # team1 = first_team_bdi.text_content().strip()
                    # team2 = second_team_bdi.text_content().strip()
                        
                    for link in match_links:
                        try:

                            # Получаем все bdi элементы внутри ссылки
                            link_bdis = link.query_selector_all('bdi')
                            commands = []
                            
                            for bdi in link_bdis:
                                text = bdi.text_content().strip()
                                if text and "Canceled" not in text and "Postponed" not in text:
                                    commands.append(text)
                            
                            if len(commands) < 4:
                                continue
                            
                            команда1 = commands[2]
                            команда2 = commands[3]

                            # Проверяем было ли более одного удаления у 1й команды
                            if len(commands) > 4:
                                if commands[3] == 'x2' or commands[3] == 'x3' or commands[3] == 'x4':
                                    commands[2] = str(commands[2]).replace(commands[3], '') # ['21/09/19', 'FT', 'Slaven', 'x2' 'Goricax2', 'x2']  то есть меняем Slavenx2 --> Slaven
                                    commands.pop(3) # удаляем 3й элемент списка
                                    

                            # Проверяем было ли более одного удаления у 2й команды
                            if len(commands) > 4:
                                if 'x' in commands[4]:
                                    commands[3] = str(commands[3]).replace(commands[4], '') # ['21/09/19', 'FT', 'Slaven', 'Goricax2', 'x2']  то есть меняем Goricax2 --> Gorica

                            
                            # Проверяем дату матча
                            try:
                                # print(commands)
                                date2 = commands[0]
                                date2_obj = datetime.strptime(date2, "%d/%m/%y")
                                if date2_obj >= date1_obj:
                                    # print('ПРОПУСКАЕМ', commands)
                                    continue
                            except:
                                continue
                            
                            # Получаем счет
                            score_elements = link.query_selector_all('[class*="currentScore"]')
                            filtered_scores = []
                            
                            for elem in score_elements:
                                text = elem.text_content().strip()
                                if text not in ['FT', 'AP', 'ПВ', 'ПП']:
                                    filtered_scores.append(text)




                            if len(filtered_scores) >= 2:

                                # Оставляем только элементы с цифрами
                                filtered_scores = [
                                                    elem.text_content().strip() 
                                                    for elem in score_elements 
                                                    if any(char.isdigit() for char in elem.text_content().strip())
                                                ]

                                home_score = filtered_scores[0]
                                away_score = filtered_scores[1]
                                
                                if '(' in home_score:
                                    home_score = int(home_score.split('(')[0])
                                else:
                                    try:
                                        home_score = int(home_score)
                                    except:
                                        home_score = 0
                                
                                if '(' in away_score:
                                    away_score = int(away_score.split('(')[0])
                                else:
                                    try:
                                        away_score = int(away_score)
                                    except:
                                        away_score = 0
                                
                                # Проверяем, играет ли команда дома
                                if команда1 in command_name_1 or command_name_1 in команда1:
                                    games_home_h2h += 1
                                    if home_score > away_score:
                                        victories_on_the_home_field += 1
                                    if home_score < away_score:
                                        loos_on_the_home_field += 1
                                    if away_score > home_score:
                                        victories_on_the_guest_field += 1
                                    if home_score == away_score:
                                        draw_on_the_home_field += 1
                                    
                                    if games_home_h2h <= 25:
                                        if команда1 in command_name_1 or command_name_1 in команда1:
                                            количество_очных_игр_на_поле_хозяев_25 += 1
                                            сумма_мячей_в_очных_играх_на_поле_хозяев_25 += (int(home_score) + int(away_score))
                                    
                                    if games_home_h2h <= 5:
                                        if команда1 in command_name_1 or command_name_1 in команда1:
                                            количество_очных_игр_на_поле_хозяев_5 += 1
                                            сумма_мячей_в_очных_играх_на_поле_хозяев_5 += (int(home_score) + int(away_score))
                                    
                                    if games_home_h2h <= 3:
                                        if команда1 in command_name_1 or command_name_1 in команда1:
                                            количество_очных_игр_на_поле_хозяев_3 += 1
                                            сумма_мячей_в_очных_играх_на_поле_хозяев_3 += (int(home_score) + int(away_score))
                                
                                games_h2h += 1
                                if games_h2h <= 25:
                                    games_h2h_25 += 1
                                    balls_h2h_25 += (int(home_score) + int(away_score))
                                    
                                    if команда1 in command_name_1 or command_name_1 in команда1:
                                        if int(home_score) + int(away_score) > 2.5:
                                            h2h_matches_over_2_5 += 1
                                        if int(home_score) + int(away_score) < 2.5:
                                            h2h_matches_under_2_5 += 1
                                
                                if games_h2h <= 5:
                                    games_h2h_5 += 1
                                    balls_h2h_5 += (int(home_score) + int(away_score))
                                
                                if games_h2h <= 3:
                                    games_h2h_3 += 1
                                    balls_h2h_3 += (int(home_score) + int(away_score))
                                
                                balls_h2h += (int(home_score) + int(away_score))
                                
                        except Exception as e:
                            continue  # Пропускаем матч с ошибкой
                            
            except Exception as e:
                print(f"Ошибка при сборе H2H данных: {e}")
                matches_data = []
            
            # Получаем ссылки на профили команд
            try:
                url_team1_element = page.query_selector('xpath=/html/body/div[1]/main/div/div[2]/div/div[1]/div[3]/div[1]/div/div[2]/div/div/div[1]/div/a')
                url_team2_element = page.query_selector('xpath=/html/body/div[1]/main/div/div[2]/div/div[1]/div[3]/div[1]/div/div[2]/div/div/div[3]/div/a')
                
                if url_team1_element and url_team2_element:
                    url_team1 = url_team1_element.get_attribute('href')
                    url_team2 = url_team2_element.get_attribute('href')
                    
                    # Получаем названия команд
                    all_bdis = page.query_selector_all('bdi')
                    if len(all_bdis) >= 4:
                        command_name_1 = all_bdis[2].text_content().strip()
                        command_name_2 = all_bdis[3].text_content().strip()
                    else:
                        print(f'команды не найдены {match_url}')
                        continue
                else:
                    print(f'ссылки на команды не найдены {match_url}')
                    continue
                    
            except Exception as e:
                print(f'Ошибка при получении ссылок на команды {match_url}: {e}')
                continue
            
            # Функция get_all_games (требует полной реализации)
            def get_all_games(page, n, command_name, command_name_1, command_name_2):
                """Полный аналог вашей функции get_all_games"""
                
                # Инициализация всех переменных как в оригинале
                balls_home = 0
                balls_away = 0
                balls_home_com1_com2 = 0
                
                result = {
                    'Название команды': '',
                    'Количество игр': 0,
                    'Количество мячей': 0,
                    'Количество мячей дома': 0,
                    'Количество мячей в гостях': 0,
                }
                
                # Функция обнуления переменных
                def zeroing():
                    nonlocal games, games_home, games_away, games_home_for_total, games_away_for_total
                    nonlocal scan_games, scan_games_home, scan_games_away
                    nonlocal win_home, loss_home, draw_home, win_away, loss_away, draw_away
                    nonlocal matches_over_2_5, matches_under_2_5, matches_over_2_5_guest, matches_under_2_5_guest
                    nonlocal matches_over_2_5_home, matches_under_2_5_home
                    nonlocal победы_команды_2_в_гостях, поражения_команды_2_в_гостях, игра_в_гостях_ничьи
                    nonlocal количество_игр_на_поле_хозяев, сумма_мячей_в_играх_на_поле_хозяев
                    nonlocal количество_игр_на_поле_гостей, сумма_мячей_в_играх_на_поле_гостей
                    nonlocal количество_попыток
                    
                    games = 0
                    games_home = 0
                    games_away = 0
                    games_home_for_total = 0
                    games_away_for_total = 0
                    scan_games = True
                    scan_games_home = True
                    scan_games_away = True
                    win_home = 0
                    loss_home = 0
                    draw_home = 0
                    win_away = 0
                    loss_away = 0
                    draw_away = 0
                    matches_over_2_5 = 0
                    matches_under_2_5 = 0
                    matches_over_2_5_guest = 0
                    matches_under_2_5_guest = 0
                    matches_over_2_5_home = 0
                    matches_under_2_5_home = 0
                    победы_команды_2_в_гостях = 0
                    поражения_команды_2_в_гостях = 0
                    игра_в_гостях_ничьи = 0
                    количество_игр_на_поле_хозяев = 0
                    сумма_мячей_в_играх_на_поле_хозяев = 0
                    количество_игр_на_поле_гостей = 0
                    сумма_мячей_в_играх_на_поле_гостей = 0
                    количество_попыток = 0
                
                # Инициализация переменных
                zeroing()
                
                while scan_games or scan_games_home or scan_games_away:
                    
                    try:
                        if "Error 503 backend read error" in page.content():
                            print('На странице ошибка: "Error 503 backend read error"')
                            page.reload(wait_until="domcontentloaded", timeout=5000)
                            print('Error 503 backend read error. Обновили страницу')
                            zeroing()
                            continue
                    except:
                        pass
                    
                    # Находим все карточки матчей
                    cards = page.query_selector_all('a[data-id]')
                    
                    for card in cards:
                        try:
                            # Ищем элементы с текстом
                            bdi_elements = card.query_selector_all('bdi')
                            commands = []
                            
                            for bdi in bdi_elements:
                                text = bdi.text_content().strip()
                                if text and "Canceled" not in text and "Postponed" not in text:
                                    commands.append(text)
                            
                            if len(commands) < 4:
                                continue
                            
                            команда1 = commands[2]
                            команда2 = commands[3]
                            
                            # Проверяем дату
                            try:
                                date2 = commands[0]
                                date2_obj = datetime.strptime(date2, "%d/%m/%y")
                                if date2_obj >= date1_obj:
                                    continue
                            except:
                                continue
                            
                            # Получаем счет
                            score_elements = card.query_selector_all('[class*="currentScore"]')
                            scores = []
                            
                            for elem in score_elements:
                                text = elem.text_content().strip()
                                if '(' in text:
                                    text = text.split('(')[0]
                                if text and text.isdigit() and ":" not in text:
                                    scores.append(text)
                            
                            if len(scores) >= 2:
                                счет_первой_команды = int(scores[0])
                                счет_второй_команды = int(scores[1])
                            else:
                                continue
                            
                            # Ваша оригинальная логика обработки матчей
                            # Я сохраняю ее полностью, но адаптирую синтаксис
                            
                            # Сканируем любые матчи
                            if games < n:
                                try:
                                    if команда1 in command_name or command_name in команда1:
                                        games += 1
                                        balls_home_com1_com2 += счет_первой_команды + счет_второй_команды
                                    
                                    if команда2 in command_name or command_name in команда2:
                                        games += 1
                                        balls_home_com1_com2 += счет_первой_команды + счет_второй_команды
                                except:
                                    pass
                            else:
                                scan_games = False
                            
                            # Сканируем домашние матчи
                            if games_home < n:
                                try:
                                    if команда1 in command_name or command_name in команда1:
                                        games_home += 1
                                        balls_home += счет_первой_команды
                                        
                                        if счет_первой_команды > счет_второй_команды:
                                            win_home += 1
                                        if счет_второй_команды > счет_первой_команды:
                                            loss_home += 1
                                        if счет_второй_команды == счет_первой_команды:
                                            draw_home += 1
                                        if счет_первой_команды + счет_второй_команды > 2.5:
                                            matches_over_2_5 += 1
                                        if счет_первой_команды + счет_второй_команды < 2.5:
                                            matches_under_2_5 += 1
                                    
                                    if команда1 in command_name_1 or command_name_1 in команда1:
                                        количество_игр_на_поле_хозяев += 1
                                        сумма_мячей_в_играх_на_поле_хозяев += счет_первой_команды + счет_второй_команды
                                except:
                                    pass
                            else:
                                scan_games_home = False
                            
                            # Сканируем гостевые матчи
                            if games_away < n:
                                try:
                                    if команда2 in command_name or command_name in команда2:
                                        games_away += 1
                                        balls_away += счет_второй_команды
                                        
                                        if счет_первой_команды < счет_второй_команды:
                                            win_away += 1
                                        if счет_первой_команды > счет_второй_команды:
                                            loss_away += 1
                                        if счет_первой_команды == счет_второй_команды:
                                            draw_away += 1
                                    
                                    if команда2 in command_name_2 or command_name_2 in команда2:
                                        if счет_первой_команды < счет_второй_команды:
                                            победы_команды_2_в_гостях += 1
                                        if счет_первой_команды > счет_второй_команды:
                                            поражения_команды_2_в_гостях += 1
                                        if счет_первой_команды == счет_второй_команды:
                                            игра_в_гостях_ничьи += 1
                                        
                                        if games_away <= n:
                                            количество_игр_на_поле_гостей += 1
                                            сумма_мячей_в_играх_на_поле_гостей += счет_первой_команды + счет_второй_команды
                                except:
                                    pass
                            else:
                                scan_games_away = False
                            
                            # Сканируем домашние матчи для подсчета тотала
                            if games_home_for_total < n:
                                try:
                                    if команда1 in command_name or command_name in команда1:
                                        games_home_for_total += 1
                                        if счет_первой_команды + счет_второй_команды > 2.5:
                                            matches_over_2_5_home += 1
                                        if счет_первой_команды + счет_второй_команды < 2.5:
                                            matches_under_2_5_home += 1
                                except:
                                    pass
                            
                            # Сканируем гостевые матчи для подсчета тотала
                            if games_away_for_total < n:
                                try:
                                    if команда2 in command_name or command_name in команда2:
                                        games_away_for_total += 1
                                        if счет_первой_команды + счет_второй_команды > 2.5:
                                            # print(f'{команда1} {счет_первой_команды} - {счет_второй_команды} {команда2}')
                                            matches_over_2_5_guest += 1
                                        if счет_первой_команды + счет_второй_команды < 2.5:
                                            matches_under_2_5_guest += 1
                                except:
                                    pass
                        
                        except Exception as e:
                            continue  # Пропускаем карточку с ошибкой
                    
                    # Поиск кнопки для загрузки предыдущих матчей
                    if scan_games or scan_games_home or scan_games_away:
                        try:
                            # Проверяем попап с текстом "Never miss a play"
                            try:
                                div_element = page.query_selector('xpath=//div[.//span[contains(text(), "Never miss a play")]]')
                                if div_element:
                                    favourite_button = div_element.query_selector('xpath=.//button[.//span[contains(text(), "Favourite")]]')
                                    if favourite_button:
                                        favourite_button.click()
                                        page.wait_for_timeout(1000)
                            except:
                                pass
                            
                            if количество_попыток > 10:
                                scan_games = False
                                scan_games_home = False
                                scan_games_away = False
                            
                            # Ищем кнопку со стрелочкой
                            try:
                                button = page.query_selector('xpath=//span[text()="Matches"]/ancestor::div[@class="d_none md:d_block"]/following-sibling::div//button[1]')
                                if button:
                                    if button.is_disabled():
                                        # print(f"Стрелочка <--- НЕАКТИВНА Количество нажатий: {количество_попыток} {match_url}")
                                        scan_games = False
                                        scan_games_home = False
                                        scan_games_away = False
                                    else:
                                        button.click()
                                        количество_попыток += 1
                                        page.wait_for_timeout(1000)
                                else:
                                    print("Стрелочка <--- не найдена")
                            except:
                                # check_and_refresh_full_text(page)
                                print("Ошибка сервера, перезагружаем")
                                page.reload(wait_until="domcontentloaded", timeout=5000)

                        
                        except Exception as e:
                            print(f"Ошибка при загрузке дополнительных матчей: {e}")
                            page.reload(wait_until="domcontentloaded", timeout=5000)

                
                # Заполняем результат (полная копия вашей функции)
                result['Название команды'] = command_name
                result['Количество игр'] = games
                result['Количество мячей'] = balls_home + balls_away
                result['сумма мячей хозяев'] = balls_home
                result['сумма мячей гостей'] = balls_away
                
                result['Домашние игры'] = games_home
                result['Гостевые игры'] = games_away
                
                result['Победы дома'] = win_home
                result['Ничьи'] = draw_home
                result['Поражения дома'] = loss_home
                
                result['Игра в гостях. Поражения 2й команды'] = поражения_команды_2_в_гостях
                result['Игра в гостях. Ничьи'] = игра_в_гостях_ничьи
                result['Игра в гостях. Победы 2й команды'] = победы_команды_2_в_гостях
                
                result['Все встречи. Свое поле. Тотал 2.5 Б/М'] = matches_over_2_5 + matches_under_2_5
                result['Все встречи. Свое поле. Тотал 2.5 М'] = matches_under_2_5
                result['Все встречи. Свое поле. Тотал 2.5 Б'] = matches_over_2_5
                
                result['Все встречи. Гостевое поле. Тотал 2.5 Б/М'] = matches_over_2_5_guest + matches_under_2_5_guest
                result['Все встречи. Гостевое поле. Тотал 2.5 М'] = matches_under_2_5_guest
                result['Все встречи. Гостевое поле. Тотал 2.5 Б'] = matches_over_2_5_guest
                
                result['Кол-во игр очн. (25) команда 1 и команда 2 на любом поле'] = games_h2h_25
                result['Сумма мячей очн. (25)'] = balls_h2h_25
                
                result['Кол-во игр очн. (5) команда 1 и команда 2 на любом поле'] = games_h2h_5
                result['Сумма мячей очн. (5)'] = balls_h2h_5
                
                result['Кол-во игр очн. (3) команда 1 и команда 2 на любом поле'] = games_h2h_3
                result['Сумма мячей очн. (3)'] = balls_h2h_3
                
                result['количество игр хозяев (25) на любом поле'] = games
                result['сумма мячей хозяев (25) на любом поле'] = balls_home_com1_com2
                result['количество игр хозяев (5) на любом поле'] = games
                result['сумма мячей хозяев (5) на любом поле'] = balls_home_com1_com2
                result['количество игр хозяев (3) на любом поле'] = games
                result['сумма мячей хозяев (3) на любом поле'] = balls_home_com1_com2
                
                result['количество игр гостей (25) на любом поле'] = games
                result['сумма мячей гостей (25) на любом поле'] = balls_home_com1_com2
                result['количество игр гостей (5) на любом поле'] = games
                result['сумма мячей гостей (5) на любом поле'] = balls_home_com1_com2
                result['количество игр гостей (3) на любом поле'] = games
                result['сумма мячей гостей (3) на любом поле'] = balls_home_com1_com2
                
                result['количество игр хозяев (25) на поле хозяев'] = количество_игр_на_поле_хозяев
                result['сумма мячей хозяев (25) на поле хозяев'] = сумма_мячей_в_играх_на_поле_хозяев
                
                result['количество игр хозяев (5) на поле хозяев'] = количество_игр_на_поле_хозяев
                result['сумма мячей хозяев (5) на поле хозяев'] = сумма_мячей_в_играх_на_поле_хозяев
                
                result['количество игр хозяев (3) на поле хозяев'] = количество_игр_на_поле_хозяев
                result['сумма мячей хозяев (3) на поле хозяев'] = сумма_мячей_в_играх_на_поле_хозяев
                
                result['количество игр гостей (25) на поле гостей'] = количество_игр_на_поле_гостей
                result['сумма мячей гостей (25)'] = сумма_мячей_в_играх_на_поле_гостей
                
                result['количество игр гостей (5) на поле гостей'] = количество_игр_на_поле_гостей
                result['сумма мячей гостей (5)'] = сумма_мячей_в_играх_на_поле_гостей
                
                result['количество игр гостей (3) на поле гостей'] = количество_игр_на_поле_гостей
                result['сумма мячей гостей (3)'] = сумма_мячей_в_играх_на_поле_гостей
                
                return result
            
            # Обрабатываем обе команды (полностью как в вашем коде)
            teams = [
                {'url': url_team1, 'name': command_name_1},
                {'url': url_team2, 'name': command_name_2}
            ]
            
            n_values = [25, 5, 3]
            results = {}
            
            for team in teams:
                team_results = {}
                
                for n in n_values:
                    try:
                        # print(f"https://www.sofascore.com{team['url']}")
                        page.goto(f"https://www.sofascore.com{team['url']}", wait_until="domcontentloaded", timeout=5000)
                    except PlaywrightTimeoutError:
                        page.evaluate("window.stop();")
                        page.wait_for_timeout(1000)

                    
                    team_results[n] = get_all_games(page, n, team['name'], command_name_1, command_name_2)
                    
                    if n != n_values[-1]:
                        try:
                            page.reload(wait_until="domcontentloaded", timeout=5000)
                            page.wait_for_timeout(1000)
                        except:
                            pass
                
                results[team['name']] = team_results
            
            # Заполняем финальный словарь (как в вашем коде)
            matches_results_dict['data_id'] = match_id
            matches_results_dict['href'] = match_url
            matches_results_dict['число'] = day
            matches_results_dict['месяц'] = month
            matches_results_dict['год'] = year
            matches_results_dict['время'] = found_time
            matches_results_dict['команда_1'] = command_name_1
            matches_results_dict['команда_2'] = command_name_2
            
            matches_results_dict['Сумма очных игр на поле команда 1'] = games_home_h2h
            matches_results_dict['Побед на своем поле в очных играх на поле команда 1'] = victories_on_the_home_field
            matches_results_dict['Ничьи в очных играх на поле команда 1'] = draw_on_the_home_field
            matches_results_dict['Поражений на своем поле в очных играх на поле команда 1'] = loos_on_the_home_field
            
            matches_results_dict['общее количество матчей дома первой команды'] = results[command_name_1][25]['Домашние игры']
            matches_results_dict['победа на своем поле'] = results[command_name_1][25]['Победы дома']
            matches_results_dict['ничья на своем поле'] = results[command_name_1][25]['Ничьи']
            matches_results_dict['поражение на своем поле'] = results[command_name_1][25]['Поражения дома']
            
            matches_results_dict['Общее количество матчей в гостях второй команды'] = results[command_name_2][25]['Гостевые игры']
            matches_results_dict['Поражения команды 2 в гостях'] = results[command_name_2][25]['Игра в гостях. Поражения 2й команды']
            matches_results_dict['Игра в гостях. Ничьи'] = results[command_name_2][25]['Игра в гостях. Ничьи']
            matches_results_dict['Победы команды 2 в гостях'] = results[command_name_2][25]['Игра в гостях. Победы 2й команды']
            
            matches_results_dict['Очные встречи. Свое поле. Тотал 2.5 Б/М'] = h2h_matches_over_2_5 + h2h_matches_under_2_5
            matches_results_dict['Очные встречи. Свое поле. Тотал 2.5 Б'] = h2h_matches_over_2_5
            matches_results_dict['Очные встречи. Свое поле. Тотал 2.5 М'] = h2h_matches_under_2_5
            
            matches_results_dict['Все встречи. Свое поле. Тотал 2.5 Б/М'] = results[command_name_1][25]['Все встречи. Свое поле. Тотал 2.5 Б/М']
            matches_results_dict['Все встречи. Свое поле. Тотал 2.5 Б'] = results[command_name_1][25]['Все встречи. Свое поле. Тотал 2.5 Б']
            matches_results_dict['Все встречи. Свое поле. Тотал 2.5 М'] = results[command_name_1][25]['Все встречи. Свое поле. Тотал 2.5 М']
            
            matches_results_dict['Все встречи. Гостевое поле. Тотал 2.5 Б/М'] = results[command_name_2][25]['Все встречи. Гостевое поле. Тотал 2.5 Б/М']
            matches_results_dict['Все встречи. Гостевое поле. Тотал 2.5 Б'] = results[command_name_2][25]['Все встречи. Гостевое поле. Тотал 2.5 Б']
            matches_results_dict['Все встречи. Гостевое поле. Тотал 2.5 М'] = results[command_name_2][25]['Все встречи. Гостевое поле. Тотал 2.5 М']
            
            matches_results_dict['Кол-во игр очн. (25) команда 1 и команда 2 на любом поле'] = results[command_name_1][25]['Кол-во игр очн. (25) команда 1 и команда 2 на любом поле']
            matches_results_dict['Сумма мячей очн. (25)'] = results[command_name_1][25]['Сумма мячей очн. (25)']
            matches_results_dict['Кол-во игр очн. (5) команда 1 и команда 2 на любом поле'] = results[command_name_1][5]['Кол-во игр очн. (5) команда 1 и команда 2 на любом поле']
            matches_results_dict['Сумма мячей очн. (5)'] = results[command_name_1][5]['Сумма мячей очн. (5)']
            matches_results_dict['Кол-во игр очн. (3) команда 1 и команда 2 на любом поле'] = results[command_name_1][3]['Кол-во игр очн. (3) команда 1 и команда 2 на любом поле']
            matches_results_dict['Сумма мячей очн. (3)'] = results[command_name_1][3]['Сумма мячей очн. (3)']
            
            matches_results_dict['количество игр хозяев (25) на любом поле'] = results[command_name_1][25]['количество игр хозяев (25) на любом поле']
            matches_results_dict['сумма мячей хозяев (25) на любом поле'] = results[command_name_1][25]['сумма мячей хозяев (25) на любом поле']
            matches_results_dict['количество игр хозяев (5) на любом поле'] = results[command_name_1][5]['количество игр хозяев (5) на любом поле']
            matches_results_dict['сумма мячей хозяев (5) на любом поле'] = results[command_name_1][5]['сумма мячей хозяев (5) на любом поле']
            matches_results_dict['количество игр хозяев (3) на любом поле'] = results[command_name_1][3]['количество игр хозяев (3) на любом поле']
            matches_results_dict['сумма мячей хозяев (3) на любом поле'] = results[command_name_1][3]['сумма мячей хозяев (3) на любом поле']
            
            matches_results_dict['количество игр гостей (25) на любом поле'] = results[command_name_2][25]['количество игр гостей (25) на любом поле']
            matches_results_dict['сумма мячей гостей (25) на любом поле'] = results[command_name_2][25]['сумма мячей гостей (25) на любом поле']
            matches_results_dict['количество игр гостей (5) на любом поле'] = results[command_name_2][5]['количество игр гостей (5) на любом поле']
            matches_results_dict['сумма мячей гостей (5) на любом поле'] = results[command_name_2][5]['сумма мячей гостей (5) на любом поле']
            matches_results_dict['количество игр гостей (3) на любом поле'] = results[command_name_2][3]['количество игр гостей (3) на любом поле']
            matches_results_dict['сумма мячей гостей (3) на любом поле'] = results[command_name_2][3]['сумма мячей гостей (3) на любом поле']
            
            matches_results_dict['Кол-во игр очн. (25) команда 1 и команда 2 на поле хозяев'] = количество_очных_игр_на_поле_хозяев_25
            matches_results_dict['Сумма мячей очн. (25) команда 1 и команда 2 на поле хозяев'] = сумма_мячей_в_очных_играх_на_поле_хозяев_25
            matches_results_dict['Кол-во игр очн. (5) команда 1 и команда 2 на поле хозяев'] = количество_очных_игр_на_поле_хозяев_5
            matches_results_dict['Сумма мячей очн. (5) команда 1 и команда 2 на поле хозяев'] = сумма_мячей_в_очных_играх_на_поле_хозяев_5
            matches_results_dict['Кол-во игр очн. (3) команда 1 и команда 2  на поле хозяев'] = количество_очных_игр_на_поле_хозяев_3
            matches_results_dict['Сумма мячей очн. (3) команда 1 и команда 2 на поле хозяев'] = сумма_мячей_в_очных_играх_на_поле_хозяев_3
            
            matches_results_dict['количество игр хозяев (25) на поле хозяев'] = results[command_name_1][25]['количество игр хозяев (25) на поле хозяев']
            matches_results_dict['сумма мячей хозяев (25) на поле хозяев'] = results[command_name_1][25]['сумма мячей хозяев (25) на поле хозяев']
            matches_results_dict['количество игр хозяев (5) на поле хозяев'] = results[command_name_1][5]['количество игр хозяев (5) на поле хозяев']
            matches_results_dict['сумма мячей хозяев (5) на поле хозяев'] = results[command_name_1][5]['сумма мячей хозяев (5) на поле хозяев']
            matches_results_dict['количество игр хозяев (3) на поле хозяев'] = results[command_name_1][3]['количество игр хозяев (3) на поле хозяев']
            matches_results_dict['сумма мячей хозяев (3) на поле хозяев'] = results[command_name_1][3]['сумма мячей хозяев (3) на поле хозяев']
            
            matches_results_dict['количество игр гостей (25) на поле гостей'] = results[command_name_2][25]['количество игр гостей (25) на поле гостей']
            matches_results_dict['сумма мячей гостей (25) на поле гостей'] = results[command_name_2][25]['сумма мячей гостей (25)']
            matches_results_dict['количество игр гостей (5) на поле гостей'] = results[command_name_2][5]['количество игр гостей (5) на поле гостей']
            matches_results_dict['сумма мячей гостей (5) на поле гостей'] = results[command_name_2][5]['сумма мячей гостей (5)']
            matches_results_dict['количество игр гостей (3) на поле гостей'] = results[command_name_2][3]['количество игр гостей (3) на поле гостей']
            matches_results_dict['сумма мячей гостей (3) на поле гостей'] = results[command_name_2][3]['сумма мячей гостей (3)']
            
            matches_results_list.append(matches_results_dict)
            
            # Сохраняем в Excel
            try:
                # full_path = get_file_path()
                
                # if os.path.exists(full_path):
                #     df = pd.read_excel(full_path)
                # else:
                #     df = pd.DataFrame()
                
                new_row = pd.DataFrame([matches_results_dict])
                df = pd.concat([df, new_row], ignore_index=True)
                # df.to_excel(full_path, index=False)
                
                return df
                
            except Exception as e:
                print(f"Ошибка при сохранении в Excel: {e}")
            
            end_time = time.time()
            execution_time = end_time - start_time
            print(f'{count} из {len(list_matches)}   Время выполнения: {execution_time:.4f} секунд  {command_name_1} против {command_name_2} {match_url}')

    else:
        print("Основной div с классом 'mdDown:pt_sm' не найден")

    # Дополнительные действия...
    page.wait_for_timeout(1000)

    # Закрываем драйвер
    browser.close()
    playwright.stop()
    
