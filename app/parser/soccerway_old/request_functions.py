import requests
from time import sleep
import json

from app.parser.soccerway.constants import *


# Получить запрос по ссылке "link"
def get_response(link, decode_json=False):
    """ get_response(str_link): return 'response' """

    for i in range(MAX_COUNT_RESPONSE):

        try:
            # Edited code start
            response = requests.get(link, headers={'User-Agent': USER_AGENT})
            # Edited code end
            if DEBUG_MODE:
                print("∟ Link:", link)

            if response.status_code == 200:
                break
            else:
                sleep(DELAY_RESPONSE)

        except:
            if DEBUG_MODE:
                print("[ - Error - ] Request:", link)

    try:
        if decode_json == True:
            response = json.loads(response.content)
    except:
        if DEBUG_MODE:
            print("[ - Error -] Make Json:", response)
    finally:
        return response
