import xlwt
from datetime import datetime

from constants import *

# Создание шапки в EXEL файле
def make_header(worksheet):
    
    x = y = 0
    for label in HEADER_EXEL:
        worksheet.write(x, y, label, STYLE_HEAD)
        y += 1