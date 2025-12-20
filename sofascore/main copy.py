import pandas as pd
import time, os, re, logging, sys
from selenium import webdriver
from datetime import datetime
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.common.exceptions import TimeoutException
from bs4 import BeautifulSoup
from selenium.webdriver.chrome.options import Options
from selenium.common.exceptions import NoSuchElementException


# Отключаем логи Selenium
os.environ['WDM_LOG_LEVEL'] = '0'
os.environ['WDM_PRINT_FIRST_LINE'] = 'False'

# Отключаем логирование
logging.getLogger('selenium').setLevel(logging.WARNING)
logging.getLogger('urllib3').setLevel(logging.WARNING)

input_date = input('Введите дату в формате 2025-11-27:  ')

url = f'https://m.sofascore.com/football/{input_date}'

date1 = url.split('/')[-1]
date1_obj = datetime.strptime(date1, "%Y-%m-%d")



def base_cleanup(driver):
    """Базовая очистка для парсинга"""
    print("🧹 Выполняю базовую очистку...")
    
    # 1. Очищаем cookies (самое важное)
    driver.delete_all_cookies()
    
    # 2. Очищаем локальное хранилище
    driver.execute_script("""
        try {
            localStorage.clear();
            sessionStorage.clear();
        } catch(e) {}
    """)
    
    # 3. Закрываем лишние вкладки (оставляем только одну)
    if len(driver.window_handles) > 1:
        main_window = driver.current_window_handle
        for handle in driver.window_handles:
            if handle != main_window:
                driver.switch_to.window(handle)
                driver.close()
        driver.switch_to.window(main_window)
    
    print("✓ Очистка завершена")
    time.sleep(0.5)  # Даем время на завершение операций

# Ваш основной цикл с очисткой каждые 5 итераций
list_matches = [...]  # ваш список матчей

driver = webdriver.Chrome()
cleaned_count = 0





def get_stealth_driver_chrome(opt):
    
    option = Options()

    option.add_experimental_option("prefs", {
        "profile.managed_default_content_settings.images": 2,
        "profile.default_content_setting_values.images": 2,
        "profile.default_content_setting_values.video": 2,  # Блокировать видео
        "profile.default_content_setting_values.audio": 2,  # Блокировать аудио
        "profile.default_content_setting_values.flash": 2,  # Блокировать флеш
    })

    # Базовые настройки
    option.add_argument("--no-sandbox")
    option.add_argument("--disable-dev-shm-usage")
    option.add_argument("--disable-gpu")
    option.add_argument("--ignore-certificate-errors")
    option.add_argument("--enable-unsafe-swiftshader")


    # Отключение поп-ап
    option.add_argument("--disable-popup-blocking")
    option.add_argument("--disable-notifications")
    option.add_argument("--disable-infobars")
    option.add_argument("--disable-extensions")
    option.add_argument("--disable-web-security")
    option.add_argument("--no-first-run")
    option.add_argument("--no-default-browser-check")
    option.add_argument("--disable-component-extensions-with-background-pages")

    # Предотвращение обнаружения автоматизации
    option.add_experimental_option("prefs", {
        "profile.default_content_setting_values.notifications": 2,  # Блокировать уведомления
        "profile.default_content_setting_values.popups": 0,  # Блокировать pop-up'ы
        "profile.managed_default_content_settings.images": 2,  # 2 = Блокировать
        "profile.default_content_setting_values.images": 2,
        "credentials_enable_service": False,
        "profile.password_manager_enabled": False
    })


    # Отключение логов
    option.add_argument("--log-level=3")
    option.add_argument("--disable-logging")
    option.add_experimental_option('excludeSwitches', ['enable-logging'])

    # Скрытый режим
    # option.add_argument("--headless")

    # User-Agent
    option.add_argument("user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
    
    # Настройки против детекта
    option.add_argument("--disable-blink-features=AutomationControlled")
    option.add_experimental_option("excludeSwitches", ["enable-automation"])
    option.add_experimental_option('useAutomationExtension', False)

    # Дополнительные stealth настройки
    option.add_argument("--disable-extensions")
    option.add_argument("--disable-plugins")
    option.add_argument("--disable-images")  # Ускорит загрузку
    # option.add_argument("--disable-javascript")  # Осторожно: может сломать функциональность

    option.add_argument(opt)
    
    # Размер окна
    # option.add_argument("--window-size=1920,1080")
    option.add_argument("--start-maximized")
    
    driver = webdriver.Chrome(options=option)
    
    # Скрываем WebDriver
    driver.execute_script("Object.defineProperty(navigator, 'webdriver', {get: () => undefined})")
    driver.set_page_load_timeout(10)

    return driver

def get_stealth_driver_firefox(opt=""):
    from selenium.webdriver.firefox.options import Options
    from selenium.webdriver.firefox.service import Service
    
    option = Options()
    
    # Базовые настройки
    option.add_argument("--no-sandbox")
    option.add_argument("--disable-dev-shm-usage")
    option.add_argument("--disable-gpu")
    if opt:
        option.add_argument(opt)
    
    # Отключение поп-ап
    option.add_argument("--disable-popup-blocking")
    
    # Отключение логов
    option.add_argument("--log-level=3")
    
    # User-Agent
    option.add_argument("user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0")
    
    # Настройки через preferences (аналог prefs в Chrome)
    option.set_preference("dom.webdriver.enabled", False)
    option.set_preference("useAutomationExtension", False)
    
    # Блокировка медиа-контента
    option.set_preference("permissions.default.image", 2)  # 2 = Блокировать изображения
    option.set_preference("media.autoplay.default", 5)     # 5 = Блокировать аудио/видео
    option.set_preference("media.autoplay.enabled", False)
    option.set_preference("media.volume_scale", "0.0")     # Отключить звук
    
    # Блокировка уведомлений и всплывающих окон
    option.set_preference("permissions.default.desktop-notification", 2)
    option.set_preference("dom.popup_maximum", 0)
    option.set_preference("privacy.popups.showBrowserMessage", False)
    
    # Отключение расширений и плагинов
    option.set_preference("extensions.enabledScopes", 0)
    option.set_preference("plugin.state.flash", 0)
    
    # Отключение сервисов
    option.set_preference("signon.rememberSignons", False)
    option.set_preference("network.cookie.cookieBehavior", 1)  # Блокировать куки от сторонних сайтов
    
    # Ускорение загрузки
    option.set_preference("permissions.default.stylesheet", 2)  # Блокировать CSS (осторожно)
    
    # Против детекта автоматизации
    option.set_preference("dom.webnotifications.enabled", False)
    option.set_preference("app.update.auto", False)
    option.set_preference("app.update.enabled", False)
    option.set_preference("browser.safebrowsing.enabled", False)
    option.set_preference("browser.safebrowsing.malware.enabled", False)
    
    # Размер окна
    option.add_argument("--width=1920")
    option.add_argument("--height=1080")
    # или для максимального размера:
    # option.add_argument("--start-maximized")
    
    # Дополнительные stealth настройки
    option.set_preference("privacy.resistFingerprinting", True)
    option.set_preference("privacy.trackingprotection.enabled", True)
    option.set_preference("webgl.disabled", True)  # Отключить WebGL для снижения отпечатка
    
    # Отключение автоматических обновлений и отчетов
    option.set_preference("app.update.auto", False)
    option.set_preference("app.update.enabled", False)
    option.set_preference("browser.send_pings", False)
    option.set_preference("browser.search.update", False)
    
    driver = webdriver.Firefox(options=option)
    
    # Скрываем WebDriver (аналогично Chrome)
    driver.execute_script("Object.defineProperty(navigator, 'webdriver', {get: () => undefined})")
    
    # Устанавливаем таймаут
    driver.set_page_load_timeout(10)
    
    # Дополнительные скрипты для маскировки
    driver.execute_script("""
        Object.defineProperty(navigator, 'plugins', {
            get: () => [1, 2, 3, 4, 5]
        });
        Object.defineProperty(navigator, 'languages', {
            get: () => ['en-US', 'en']
        });
    """)
    
    return driver

# Инициализация драйвера для chrome
opt = "--force-device-scale-factor=0.05"
# driver = get_stealth_driver_chrome(opt)

# Инициализация драйвера для firefox
driver = get_stealth_driver_firefox()



try:
    driver.get(url)
    driver.set_page_load_timeout(3)

    try:
        # Найти div с текстом "Help us improve"
        div_element = WebDriverWait(driver, 10).until(
            EC.presence_of_element_located((By.XPATH, '//div[.//span[contains(text(), "Help us improve")]]'))
        )
        
        # Найти все кнопки внутри этого div
        buttons = div_element.find_elements(By.TAG_NAME, 'button')
        
        # Кликнуть по первой кнопке (это будет кнопка закрытия с крестиком)
        if buttons:
            buttons[0].click()
            print("Первая кнопка (крестик) нажата")
        else:
            print("Кнопки не найдены")

    except Exception as e:
        print(f"Ошибка: {e}")



    button = driver.find_element(By.XPATH, '//div[contains(@class, "mdDown:pt_sm")]//button[contains(text(), "Show all")]')
    if button:
        button.click()
        print("Нажали Show all")
        time.sleep(10)

    # print(f"Страница загружена {url}")
except TimeoutException:
    print("Прерываем долгую загрузку...")
    driver.execute_script("window.stop();")



sofascore_page = driver.page_source


def check_and_refresh_full_text(driver):
    try:
        # Полный текст для поиска
        full_text = "Add to Favourites to keep track of upcoming events. You can adjust this and notifications later in the Favourites tab."
        
        # Ищем span с точным текстом
        element = driver.find_element((By.XPATH, f"//span[text()='{full_text}']"))
        
        if element:
            print("Найден span с текстом про Favourites, делаю refresh...")
            driver.refresh()
            # time.sleep(3)
            return True
            
    except Exception as e:
        # Элемент не найден - это нормально
        return False



def get_file_path():
    """Получает путь к data.xlsx, который лежит на уровень выше exe"""
    
    if getattr(sys, 'frozen', False):
        # Для exe - папка где находится exe файл
        exe_dir = os.path.dirname(sys.executable)
        parent_dir = os.path.dirname(exe_dir)  # поднимаемся на уровень выше
        data_path = os.path.join(parent_dir, "data.xlsx")
    else:
        # Для скрипта - текущая папка скрипта
        current_dir = os.path.abspath(".")
        data_path = os.path.join(current_dir, "data.xlsx")
    
    return data_path


try:
    # Ждем появления хотя бы одного матча
    WebDriverWait(driver, 20).until(EC.presence_of_element_located((By.CSS_SELECTOR, 'a[data-id][href*="/match/"]')))
    
    # Используем более специфичный селектор
    match_links = driver.find_elements(By.CSS_SELECTOR, 'div[class*="mdDown:pt_sm"] a[data-id][href*="/match/"]')
    
    list_matches = []
    
    for link in match_links:
        href = link.get_attribute('href')
        data_id = link.get_attribute('data-id')
        if href:
            if href.startswith('/'):
                href = f"https://www.sofascore.com{href}"
            list_matches.append(href)
    
    print(f"Найдено матчей: {len(list_matches)}")
except Exception as e:
    print(f"Ошибка: {e}")



if len(list_matches) > 0:
        
        matches_results_list = []

        full_path = get_file_path()

        # Чтение файла
        if os.path.exists(full_path):
            df = pd.read_excel(full_path)
            # Получить все значения из столбца B
            column_b_values = df['href'].tolist()  # или df.iloc[:, 1].tolist()
        else:
            column_b_values = []
        
        count = 0

        driver.quit()
        opt = "--force-device-scale-factor=1.0"
        driver = get_stealth_driver_chrome(opt)

        
        for match_url in list_matches:
            bad_connect = 0
            count += 1

            if count % 5 == 0:
                print(f"\n⚠️  Достигнут цикл #{count} - выполняю очистку")
                base_cleanup(driver)

            if match_url in column_b_values:
                print(f'Матч {match_url} уже обрабатывался')
                continue

            start_time = time.time()

            match_id = match_url.split('#id:')[-1]
            # Словарь, в который будет записываться в финальный список, а этот список записываться в Excel
            matches_results_dict = {}
            try:
                  # Устанавливаем короткий таймаут
                driver.get(match_url)
                # driver.get('https://www.sofascore.com/football/match/pari-nizhny-novgorod-cska-moscow/AWseIFb#id:14038566')
                # print(f"Страница {match_url} загружена")
            except TimeoutException:
                # print(f"Принудительно останавливаем загрузку {match_url}...")
                driver.execute_script("window.stop();")  # Останавливаем загрузку
            
            # Дополнительная проверка - ждем появления основного контента
            try:
                wait = WebDriverWait(driver, 5)
                # Ждем появления любого значимого элемента страницы
                wait.until(EC.presence_of_element_located((By.TAG_NAME, "body")))
            except:
                print("Основной контент не загрузился")


            max_retries = 3
            retry_count = 0
            success = False

            while retry_count < max_retries and not success:
                try:
                    div_with_favourite = driver.find_element(By.XPATH, '//div[not(@class) and contains(., "Favourite")]')
                    data_list = div_with_favourite.find_elements(By.TAG_NAME, 'span')
                    success = True  # успешно
                except:
                    retry_count += 1
                    print(f"Попытка {retry_count}/{max_retries} не удалась, перезагружаем страницу {match_url}")
                    driver.refresh()
                    time.sleep(4)

            if not success:
                print(f"Не удалось за {max_retries} попыток, пропускаем {match_url}")
                continue  # пропускаем эту итерацию внешнего цикла


            date_pattern = r'\b\d{2}/\d{2}/\d{4}\b'
            time_pattern = r'\b\d{2}:\d{2}\b'

            found_date = None
            found_time = None
            day = 0
            month = 0
            year = 0


            # Итерируем по списку
            for i, item in enumerate(data_list):
                text = item.text.strip()
                
                # Ищем дату
                if not found_date and re.match(date_pattern, text):
                    found_date = text
                    day, month, year = found_date.split('/')
                    
                # Ищем время
                if not found_time and re.match(time_pattern, text):
                    found_time = text
                
                # Если нашли оба значения, можно прервать цикл
                if found_date and found_time:
                    break



            # Теперь ищем и кликаем H2H
            try:
                wait = WebDriverWait(driver, 5)
                h2h_button = wait.until(EC.element_to_be_clickable((By.CSS_SELECTOR, 'button[data-testid="tab-matches"]')))
                
                driver.execute_script("arguments[0].click();", h2h_button)
                # print("Кликнут элемент H2H")
                time.sleep(1)
            except Exception as e:
                print(f"Элемент H2H не найден: {e}")


            # Теперь ищем кнопку "Показать больше"
            try:
                target_div = driver.find_element(By.CSS_SELECTOR, 'div[class="w_100%"]')
                
                # Пробуем найти кнопку на русском
                try:
                    show_more_button = target_div.find_element(By.XPATH, './/button[text()="Показать больше"]')
                except:
                    # Если не нашли на русском, ищем на английском
                    show_more_button = target_div.find_element(By.XPATH, './/button[text()="Show more"]')
                
                # Кликаем через JavaScript
                driver.execute_script("arguments[0].click();", show_more_button)
                # print("Кликнули на кнопку 'Показать больше'")
                time.sleep(1)
            except Exception as e:
                pass
                # print(f"Кнопка 'Показать больше' не найдена: {e}")


            # Собираем все ссылки H2H
            try:
                # Находим div с классом w_100%
                main_div = driver.find_element(By.CSS_SELECTOR, 'div[class="w_100%"]')
                
                # Находим все ссылки с data-id внутри этого div
                match_links = main_div.find_elements(By.CSS_SELECTOR, 'a[data-id]')
                
                # Собираем данные о матчах
                matches_data = []
                games_h2h = 0  # переменная хранит количество очных игр
                games_h2h_25 = 0  # переменная хранит количество очных игр (25 игр, если есть)
                games_h2h_3 = 0  # переменная хранит количество очных игр (3 игры, если есть)
                games_h2h_5 = 0  # переменная хранит количество очных игр (5 игры, если есть)
                games_home_h2h = 0
                h2h_matches_over_2_5 = 0
                h2h_matches_under_2_5 = 0

                balls_h2h = 0  # переменная хранит сумму мячей в очных играх
                balls_h2h_25 = 0  # переменная хранит сумму мячей в очных играх (25 игр, если есть)
                balls_h2h_3 = 0  # переменная хранит сумму мячей в очных играх (3 игры, если есть)
                balls_h2h_5 = 0  # переменная хранит сумму мячей в очных играх (5 игры, если есть)

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

                for link in match_links:
                    
                    # Названия команд
                    first_team_element = link.find_element(By.XPATH, './/div[1]/div[4]/div[1]/div[1]/div[1]//bdi')
                    team1 = first_team_element.text
                    second_team_element = link.find_element(By.XPATH, './/div[1]/div[4]/div[1]/div[1]/div[2]//bdi')
                    team2 = second_team_element.text
                    
                    # Счет
                    try:
                        # Находим все элементы с currentScore внутри тега a
                        score_elements = link.find_elements(By.CSS_SELECTOR, '[class*="currentScore"]')

                        commands = link.find_elements(By.TAG_NAME, 'bdi')
                        commands = [cmd for cmd in commands if cmd.text.strip() and "Canceled" not in cmd.text and "Postponed" not in cmd.text] # удаляем пустой текст

                        команда1 = commands[2].text
                        команда2 = commands[3].text

                        date2 = commands[0].text
                        date2_obj = datetime.strptime(date2, "%d/%m/%y") 
                        if date2_obj >= date1_obj:
                            continue

                        coms = []
                        for elem in commands:
                            coms.append(elem.text.strip())

                        command_name_1 = driver.find_elements(By.TAG_NAME, 'bdi')[2].text
                        command_name_2 = driver.find_elements(By.TAG_NAME, 'bdi')[3].text
                        
                        # Фильтруем элементы, исключая те, у которых текст FT или AT
                        filtered_scores = []
                        for elem in score_elements:
                            text = elem.text.strip()
                            if text not in ['FT', 'AP', 'ПВ', 'ПП']:  # можно добавить другие статусы
                                filtered_scores.append(elem)

                        # Берем тексты первых двух отфильтрованных элементов
                        if len(filtered_scores) >= 2:

                            home_score = filtered_scores[0].text
                            away_score = filtered_scores[1].text

                            if '(' in home_score: home_score = int(home_score.split('(')[0])
                            if '(' in away_score: away_score = int(away_score.split('(')[0])

                            if команда1 in command_name_1 or command_name_1 in команда1:  # это условие означает, что команда играет дома
                                
                                games_home_h2h += 1

                                # Считаем победы, ничьи и поражения в очных встречах
                                if home_score > away_score:
                                    # print(f"ДОМАШНИЙ МАТЧ: {команда1} {home_score} //// {away_score} {команда2}")
                                    victories_on_the_home_field += 1
                                if home_score < away_score:
                                    loos_on_the_home_field += 1
                                if away_score > home_score:
                                    victories_on_the_guest_field += 1
                                if home_score == away_score:
                                    draw_on_the_home_field += 1

                                if games_home_h2h <= 25:
                                    if команда1 in command_name_1 or command_name_1 in команда1: # это условие означает, что команда играет дома
                                        количество_очных_игр_на_поле_хозяев_25 += 1
                                        сумма_мячей_в_очных_играх_на_поле_хозяев_25 += (int(home_score) + int(away_score))

                                if games_home_h2h <= 5:
                                    if команда1 in command_name_1 or command_name_1 in команда1: # это условие означает, что команда играет дома
                                        количество_очных_игр_на_поле_хозяев_5 += 1
                                        сумма_мячей_в_очных_играх_на_поле_хозяев_5 += (int(home_score) + int(away_score))

                                if games_home_h2h <= 3:
                                    if команда1 in command_name_1 or command_name_1 in команда1: # это условие означает, что команда играет дома
                                        количество_очных_игр_на_поле_хозяев_3 += 1
                                        сумма_мячей_в_очных_играх_на_поле_хозяев_3 += (int(home_score) + int(away_score))



                            games_h2h += 1
                            if games_h2h <= 25:
                                
                                games_h2h_25 += 1 # считаем количество очных игр (25 игр)
                                balls_h2h_25 += (int(home_score) + int(away_score)) # считаем сумму мячей в 25 очных играх
                                # print(f"{games_h2h_25} из {25} матчей. {команда1} {int(home_score)} //// {int(away_score)} {команда2}")

                                if команда1 in command_name_1 or command_name_1 in команда1: # это условие означает, что команда играет дома
                                    if int(home_score) + int(away_score) > 2.5: h2h_matches_over_2_5 += 1
                                    if int(home_score) + int(away_score) < 2.5: h2h_matches_under_2_5 += 1

                            if games_h2h <= 5:
                                games_h2h_5 += 1 # считаем количество очных игр (5 игры)
                                balls_h2h_5 += (int(home_score) + int(away_score)) # считаем сумму мячей 5 очных играх


                            if games_h2h <= 3:
                                games_h2h_3 += 1 # считаем количество очных игр (3 игры)
                                balls_h2h_3 += (int(home_score) + int(away_score)) # считаем сумму мячей 3 очных играх


                            balls_h2h += (int(home_score) + int(away_score)) # считаем сумму мячей во всех очных играх
                            
                        elif len(filtered_scores) == 1:
                            continue
                            home_score = filtered_scores[0].text
                            away_score = ""
                            # print(f"Только один счет: {home_score}")
                        else:
                            continue
                            home_score = ""
                            away_score = ""
                            # print("Счетов не найдено")
                            
                    except Exception as e:
                        # print(f"Ошибка поиска счета: {e}")
                        home_score = ""
                        away_score = ""
                        
                    # Дата матча
                    match_date = link.find_element(By.XPATH, './/div[1]/div[2]//bdi')
                    match_date = match_date.text

                    # Дата
                    try:
                        date_elem = link.find_element(By.CSS_SELECTOR, 'bdi.textStyle_body.small.c_neutrals.nLv3')
                        date = date_elem.text
                    except:
                        date = ""
                    
                    # Собираем все данные в словарь
                    match_info = {
                        'data_id': "",
                        'href': href,
                        'match_date': match_date,
                        'team1': team1,
                        'team2': team2,
                        'score1': home_score,
                        'score2': away_score,
                        'date': date,
                        'result': f"{team1} {home_score}-{away_score} {team2}"
                    }
                    
                    matches_data.append(match_info)


                # Посчитаем сколько каждая команда сыграла домашних и гостевых матчей
                team_stats = {}
                for match in matches_data:
                    home_team = match['team1']
                    away_team = match['team2']
                    
                    # Инициализируем запись для команды, если ее еще нет
                    if home_team not in team_stats: team_stats[home_team] = {'home': 0, 'away': 0}
                    if away_team not in team_stats: team_stats[away_team] = {'home': 0, 'away': 0}
                    
                    # Увеличиваем счетчики
                    team_stats[home_team]['home'] += 1
                    team_stats[away_team]['away'] += 1

                
            except Exception as e:
                print(f"Ошибка: {e}")
                matches_data = []
                matches_results_list = []



            # ________________________________Найдем ссылки на профиль команд и будем их парсить___________________________________
            def get_all_games(driver, n, command_name, command_name_1, command_name_2): # функция собирает все матчи

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
                
                def zeroing():
                    global games, games_home, games_away, games_home_for_total, games_away_for_total
                    global scan_games, scan_games_home, scan_games_away
                    global win_home, loss_home, draw_home, win_away, loss_away, draw_away
                    global matches_over_2_5, matches_under_2_5, matches_over_2_5_guest, matches_under_2_5_guest
                    global matches_over_2_5_home, matches_under_2_5_home
                    global победы_команды_2_в_гостях, поражения_команды_2_в_гостях, игра_в_гостях_ничьи
                    global количество_игр_на_поле_хозяев, сумма_мячей_в_играх_на_поле_хозяев
                    global количество_игр_на_поле_гостей, сумма_мячей_в_играх_на_поле_гостей
                    global количество_попыток
                    
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

                количество_игр_на_поле_хозяев =0
                сумма_мячей_в_играх_на_поле_хозяев = 0

                количество_игр_на_поле_гостей = 0
                сумма_мячей_в_играх_на_поле_гостей = 0

                количество_попыток = 0
                
                while scan_games or scan_games_home or scan_games_away:
                    
                    # print(scan_games, scan_games_home, scan_games_away)

                    if "Error 503 backend read error" in driver.page_source:
                        print('На странице ошибка: "Error 503 backend read error"')
                        driver.refresh()
                        time.sleep(4)
                        print('Error 503 backend read error. Обновили страницу')
                        # time.sleep(444)
                        scan_games = True
                        scan_games_home = True
                        scan_games_away = True
                        zeroing()

                    cards = driver.find_elements(By.CSS_SELECTOR, 'a[data-id]') # Находим все матчи

                    
                    for index, card in enumerate(cards):
                        
                        try:
                            # Ищем ВСЕ элементы с currentScore для отладки
                            all_score_elements = card.find_element(By.TAG_NAME, 'div').find_elements(By.CSS_SELECTOR, '[class*="currentScore"]')

                            commands = card.find_elements(By.TAG_NAME, 'bdi') # команды, которые играют между собой
                            commands = [cmd for cmd in commands if cmd.text.strip() and "Canceled" not in cmd.text and "Postponed" not in cmd.text] # удаляем пустой текст

                            команда1 = commands[2].text
                            команда2 = commands[3].text

                            # Пропускаем матч, если он был до заданной даты
                            date2 = commands[0].text
                            date2_obj = datetime.strptime(date2, "%d/%m/%y") 
                            if date2_obj >= date1_obj:
                                continue

                            scores = []
                            for span in all_score_elements: # Тут мы собираем количество голов 
                                text = span.text.strip()
                                if '(' in text: text = text.split('(')[0]
                                if text and text.isdigit() and ":" not in text:
                                    scores.append(text)
                            счет_первой_команды = int(scores[0])
                            счет_второй_команды = int(scores[1])
                        except: 
                            # print(commands)
                            continue


                        # Сканируем любые матчи (дома / в гостях)
                        if games < n:
                            try:
                                if команда1 in command_name or command_name in команда1:  # это условие означает, что команда играет дома
                                    games += 1
                                    balls_home_com1_com2 += счет_первой_команды + счет_второй_команды
                                    # print(f"{games} из {n} матчей. ДОМАШНИЙ МАТЧ: {команда1} {счет_первой_команды} - {счет_второй_команды} {команда2}  /  {balls_home_com1_com2}")

                                if команда2 in command_name or command_name in команда2:  # это условие означает, что команда играет в гостях
                                    games += 1
                                    balls_home_com1_com2 += счет_первой_команды + счет_второй_команды
                                    # print(f"{games} из {n} матчей. ГОСТЕВОЙ МАТЧ: {команда1} {счет_первой_команды} - {счет_второй_команды} {команда2}  /  {balls_home_com1_com2}")

                            except Exception as e:
                                pass
                        else: scan_games = False


                        # Сканируем домашние матчи
                        if games_home < n:
                            try:
                                if команда1 in command_name or command_name in команда1:  # это условие означает, что команда играет дома
                                    games_home += 1
                                    # print(f"{games_home} из {n} матчей. ДОМАШНИЙ МАТЧ: {команда1} //// {команда2}")
                                    balls_home += счет_первой_команды
                                    
                                    if счет_первой_команды > счет_второй_команды: win_home += 1
                                    if счет_второй_команды > счет_первой_команды: loss_home += 1
                                    if счет_второй_команды == счет_первой_команды: draw_home += 1
                                    if счет_первой_команды + счет_второй_команды > 2.5: matches_over_2_5 += 1
                                    if счет_первой_команды + счет_второй_команды < 2.5: matches_under_2_5 += 1
                                else:
                                    if счет_первой_команды + счет_второй_команды > 2.5: test1 += 1
                                    if счет_первой_команды + счет_второй_команды < 2.5: test2 += 1

                                if команда1 in command_name_1 or command_name_1 in команда1: # это условие означает, что команда играет дома
                                    количество_игр_на_поле_хозяев += 1
                                    сумма_мячей_в_играх_на_поле_хозяев += счет_первой_команды + счет_второй_команды
                                        # print(f"{games_home} из {n} матчей. ДОМАШНИЙ МАТЧ: {команда1} {счет_первой_команды} //// {счет_второй_команды} {команда2} сумма сячей {сумма_мячей_в_играх_на_поле_хозяев}")

                            except Exception as e:
                                pass

                        else: scan_games_home = False



                        # Сканируем гостевые матчи
                        if games_away < n:
                            try:
                                if команда2 in command_name or command_name in команда2:  # это условие означает, что команда играет в гостях
                                    games_away += 1
                                    # print(f"{games_away} из {n} матчей. ГОСТЕВОЙ МАТЧ: {команда1} {scores[0]} //// {scores[1]} {команда2}")
                                    balls_away += счет_второй_команды
                                    
                                    if scores[0] < scores[1]: win_away += 1
                                    if scores[0] > scores[1]: loss_away += 1
                                    if scores[0] == scores[1]: draw_away += 1


                                if команда2 in command_name_2 or command_name_2 in команда2: # это условие означает, что команда играет в гостях
                                    
                                    if scores[0] < scores[1]: победы_команды_2_в_гостях += 1
                                    if scores[0] > scores[1]: поражения_команды_2_в_гостях += 1
                                    if scores[0] == scores[1]: игра_в_гостях_ничьи += 1

                                # ____________________________________________________________
                                if games_away <= n:
                                    if команда2 in command_name_2 or command_name_2 in команда2: # это условие означает, что команда играет в гостях
                                        количество_игр_на_поле_гостей += 1
                                        сумма_мячей_в_играх_на_поле_гостей += счет_первой_команды + счет_второй_команды
                                        # print(f"{games_away} из {n} матчей. ГОСТЕВОЙ МАТЧ: {команда1} {счет_первой_команды} //// {счет_второй_команды} {команда2} сумма сячей {сумма_мячей_в_играх_на_поле_гостей}")


                            except Exception as e:
                                pass

                        else: scan_games_away = False

                    
                        # Сканируем домашние матчи для подсчета тотала
                        if games_home_for_total < n:
                            try:
                                if команда1 in command_name or command_name in команда1:  # это условие означает, что команда играет дома
                                    games_home_for_total += 1
                                    if счет_первой_команды + счет_второй_команды > 2.5:
                                        # print(f"ДОМАШНИЙ МАТЧ: {команда1} {счет_первой_команды} - {счет_второй_команды} {команда2}")
                                        matches_over_2_5_home += 1
                                    if счет_первой_команды + счет_второй_команды < 2.5: matches_under_2_5_home += 1


                            except Exception as e:
                                pass

                        else: scan_games_home = False


                        # Сканируем гостевые матчи для подсчета тотала

                        if games_away_for_total < n:
                            try:
                                if команда2 in command_name or command_name in команда2:  # это условие означает, что команда играет в гостях
                                    games_away_for_total += 1
                                    # print(f"{games_away_for_total} из {n} матчей. ГОСТЕВОЙ МАТЧ: {команда1} {счет_первой_команды} //// {счет_второй_команды} {команда2} сумма мячей {сумма_мячей_в_играх_на_поле_гостей}")
                                    if счет_первой_команды + счет_второй_команды > 2.5: matches_over_2_5_guest += 1
                                    if счет_первой_команды + счет_второй_команды < 2.5: matches_under_2_5_guest += 1
                            except Exception as e:
                                pass

                        else: scan_games_away = False



                    # Ищем стрелочку <--- чтобы собрать все матчи
                    if scan_games or scan_games_home or scan_games_away:

                        try:
                            # Находим div с текстом "Never miss a play"
                            div_element = driver.find_element(By.XPATH, '//div[.//span[contains(text(), "Never miss a play")]]')
                            
                            # Находим кнопку с текстом "Favourite" внутри этого div
                            favourite_button = div_element.find_element(By.XPATH, './/button[.//span[contains(text(), "Favourite")]]')
                            
                            favourite_button.click()
                            time.sleep(1)
                            print("Кнопка 'Favourite' найдена и нажата")

                        except:
                            pass


                        # print(f"количество попыток: {количество_попыток}")
                        if количество_попыток > 10:
                            scan_games = False
                            scan_games_home = False
                            scan_games_away = False
                        try:
                            buttons = driver.find_element(By.XPATH, '//span[text()="Matches"]/ancestor::div[@class="d_none md:d_block"]/following-sibling::div//button[1]')
                            if buttons:
                                buttons.click()
                                количество_попыток += 1
                                time.sleep(1)
                            else:
                                print(f"Стрелочка <--- не найдена")

                        except:
                            # check_and_refresh_full_text(driver)
                            print(f"Ошибка сервера, перезагружаем")
                            driver.refresh()
                            time.sleep(4)
                           

                            # driver.save_screenshot(r"C:\Users\user\Desktop\sofascore\screenshot.png")



                                


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



            try:
                url_team1 = driver.find_element(By.XPATH, '/html/body/div[1]/main/div[2]/div/div/div[1]/div[3]/div/div[2]/div/div[1]/div[1]/div/a').get_attribute("href")
                url_team2 = driver.find_element(By.XPATH, '/html/body/div[1]/main/div[2]/div/div/div[1]/div[3]/div/div[2]/div/div[1]/div[3]/div/a').get_attribute("href")
            
            except:
                print(f'команды не найдены {match_url}')
                continue
            
            command_name_1 = driver.find_elements(By.TAG_NAME, 'bdi')[2].text
            command_name_2 = driver.find_elements(By.TAG_NAME, 'bdi')[3].text

            # Список команд для обработки
            teams = [
                {'url': url_team1, 'name': command_name_1},
                {'url': url_team2, 'name': command_name_2}
            ]

            # Список значений n для перебора
            n_values = [25, 5, 3]
            results = {}

            # Обрабатываем обе команды
            for team in teams:
                
                team_results = {}
                
                for n in n_values:
                    try:
                        driver.set_page_load_timeout(3)
                        driver.get(f"{team['url']}")
                    except TimeoutException:
                        driver.execute_script("window.stop();")  # На всякий случай

                    team_results[n] = get_all_games(driver, n, team['name'], command_name_1, command_name_2)
                    
                    # Обновляем страницу только если это не последняя итерация
                    if n != n_values[-1]:
                        try:
                            driver.refresh()

                        except TimeoutException:
                            driver.execute_script("window.stop();")
                            
                
                results[team['name']] = team_results


            
            matches_results_dict['data_id'] = match_id
            matches_results_dict['href'] = match_url
            matches_results_dict['число'] = day
            matches_results_dict['месяц'] = month
            matches_results_dict['год'] = year
            matches_results_dict['время'] = found_time
            matches_results_dict['команда_1'] = team1
            matches_results_dict['команда_2'] = team2

            # matches_results_dict['Сумма'] = games_h2h
            # matches_results_dict['1'] = victories_on_the_home_field
            # matches_results_dict['x'] = draw
            # matches_results_dict['2'] = victories_on_the_guest_field
            matches_results_dict['Сумма очных игр на поле команда 1'] = games_home_h2h
            matches_results_dict['Побед на своем поле в очных играх на поле команда 1'] = victories_on_the_home_field
            matches_results_dict['Ничьи в очных играх на поле команда 1'] = draw_on_the_home_field
            matches_results_dict['Поражений на своем поле в очных играх на поле команда 1'] = loos_on_the_home_field

            # matches_results_dict['сумма'] = results[command_name_1][25]['Домашние игры']
            # matches_results_dict['1'] = results[command_name_1][25]['Победы дома']
            # matches_results_dict['x'] = results[command_name_1][25]['Ничьи']
            # matches_results_dict['2'] = results[command_name_1][25]['Поражения дома']
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

            matches_results_dict['Все встречи. Гостевое поле. Тотал 2.5 Б/М'] = results[command_name_1][25]['Все встречи. Гостевое поле. Тотал 2.5 Б/М']
            matches_results_dict['Все встречи. Гостевое поле. Тотал 2.5 Б'] = results[command_name_1][25]['Все встречи. Гостевое поле. Тотал 2.5 Б']
            matches_results_dict['Все встречи. Гостевое поле. Тотал 2.5 М'] = results[command_name_1][25]['Все встречи. Гостевое поле. Тотал 2.5 М']


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

            # cwd = os.path.dirname(__file__)
            # file_name = f'data.xlsx'
            # full_path = os.path.join(cwd, file_name)            
            full_path = get_file_path()
            # Создаем или загружаем существующий Excel файл
            try:
                df = pd.read_excel(full_path)
            except FileNotFoundError:
                df = pd.DataFrame()

            # Создаем DataFrame из текущего словаря
            new_row = pd.DataFrame([matches_results_dict])

            # Добавляем новую строку
            df = pd.concat([df, new_row], ignore_index=True)

            # Сохраняем в Excel
            df.to_excel(full_path, index=False)

            

            end_time = time.time()
            execution_time = end_time - start_time
            print(f'Сделано {count} из {len(list_matches)}      Время выполнения: {execution_time:.4f} секунд')

else:
    print("Основной div с классом 'mdDown:pt_sm' не найден")

# Дополнительные действия...
time.sleep(1)

# Закрываем драйвер
driver.quit()

# parser_page(sofascore_page)


