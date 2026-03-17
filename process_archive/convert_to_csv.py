"""
xlsx_to_csv.py

Конвертирует все .xlsx файлы из archive/excel в archive/csv
Структура выходных папок по году из имени файла:

    soccer_2022.xlsx  (10 листов)  →  archive/csv/2022/Лига1.csv
                                       archive/csv/2022/Лига2.csv
                                       ...
    data_2023.xlsx                 →  archive/csv/2023/...

Год извлекается из имени файла (первые 4 цифры подряд).
Если год не найден — файл кладётся в archive/csv/unknown/

Использование:
    python xlsx_to_csv.py
    python xlsx_to_csv.py --input archive/excel --output archive/csv
    python xlsx_to_csv.py --sep ";"
"""

import argparse
import os
import re
import glob
import pandas as pd
from pathlib import Path


# ─── Настройки по умолчанию ───────────────────────────────────────────────────

DEFAULT_INPUT  = "archive/excel"
DEFAULT_OUTPUT = "archive/csv"
DEFAULT_SEP    = ";"
DEFAULT_ENC    = "utf-8-sig"   # utf-8-sig — Excel открывает без проблем с кириллицей


# ─── Извлечение года из имени файла ──────────────────────────────────────────

def extract_year(filename: str) -> str:
    """
    Ищет 4-значный год в имени файла.
    soccer_2022.xlsx  → '2022'
    2019_data.xlsx    → '2019'
    archive.xlsx      → 'unknown'
    """
    match = re.search(r"(19|20)\d{2}", filename)
    return match.group(0) if match else "unknown"


# ─── Конвертация одного файла ─────────────────────────────────────────────────

def convert_file(xlsx_path: str, base_output_dir: str, sep: str, encoding: str) -> list[str]:
    """
    Конвертирует один xlsx файл.
    Создаёт подпапку по году, каждый лист — отдельный CSV.
    Возвращает список путей созданных CSV файлов.
    """
    stem = Path(xlsx_path).stem
    year = extract_year(stem)

    # Папка назначения: archive/csv/2022/
    year_dir = os.path.join(base_output_dir, year)
    os.makedirs(year_dir, exist_ok=True)

    created = []

    try:
        sheets: dict = pd.read_excel(
            xlsx_path,
            sheet_name=None,   # все листы
            header=None,       # сохраняем как есть
            dtype=str,         # без автоматических преобразований
        )
    except Exception as e:
        print(f"  [!] Ошибка чтения {xlsx_path}: {e}")
        return []

    for sheet_name, df in sheets.items():
        safe_sheet = (
            sheet_name
            .strip()
            .replace("/", "-")
            .replace("\\", "-")
            .replace(" ", "_")
            .replace(":", "-")
        )

        # Имя файла = название листа (имя xlsx не дублируем — оно уже в папке года)
        # Но если лист один — используем имя xlsx файла
        if len(sheets) == 1:
            csv_name = f"{stem}.csv"
        else:
            csv_name = f"{safe_sheet}.csv"

        csv_path = os.path.join(year_dir, csv_name)

        df = df.fillna("")
        df.to_csv(
            csv_path,
            sep=sep,
            index=False,
            header=False,
            encoding=encoding,
        )

        rows = len(df)
        cols = len(df.columns)
        print(f"    → {year}/{csv_name}  ({rows} строк, {cols} колонок)")
        created.append(csv_path)

    return created


# ─── Основная функция ─────────────────────────────────────────────────────────

def main(input_dir: str, output_dir: str, sep: str, encoding: str):
    if not os.path.isdir(input_dir):
        print(f"[!] Папка не найдена: {input_dir}")
        return

    os.makedirs(output_dir, exist_ok=True)

    pattern = os.path.join(input_dir, "**", "*.xlsx")
    files = sorted(glob.glob(pattern, recursive=True))
    files += sorted(glob.glob(os.path.join(input_dir, "*.xlsx")))
    files = sorted(set(files))

    if not files:
        print(f"[!] xlsx файлы не найдены в: {input_dir}")
        return

    print(f"[+] Найдено xlsx файлов: {len(files)}")
    print(f"[+] Выходная папка: {output_dir}\n")

    total_csv   = 0
    total_errors = 0
    years_seen: dict[str, int] = {}

    for xlsx_path in files:
        fname = os.path.relpath(xlsx_path, input_dir)
        year  = extract_year(Path(xlsx_path).stem)
        print(f"[>] {fname}  (год: {year})")

        created = convert_file(xlsx_path, output_dir, sep, encoding)

        if created:
            total_csv += len(created)
            years_seen[year] = years_seen.get(year, 0) + len(created)
        else:
            total_errors += 1

    # Итог
    print(f"\n[✓] Готово: создано {total_csv} CSV файлов")
    print(f"    Структура папок:")
    for year in sorted(years_seen):
        print(f"      {output_dir}/{year}/  — {years_seen[year]} файлов")
    if total_errors:
        print(f"    Ошибок: {total_errors}")


# ─── CLI ──────────────────────────────────────────────────────────────────────

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Конвертация xlsx → csv")
    parser.add_argument(
        "--input", default=DEFAULT_INPUT,
        help=f"Папка с xlsx файлами (по умолчанию: {DEFAULT_INPUT})"
    )
    parser.add_argument(
        "--output", default=DEFAULT_OUTPUT,
        help=f"Папка для CSV файлов (по умолчанию: {DEFAULT_OUTPUT})"
    )
    parser.add_argument(
        "--sep", default=DEFAULT_SEP,
        help=f"Разделитель CSV (по умолчанию: '{DEFAULT_SEP}')"
    )
    parser.add_argument(
        "--encoding", default=DEFAULT_ENC,
        help=f"Кодировка CSV (по умолчанию: {DEFAULT_ENC})"
    )
    args = parser.parse_args()

    main(args.input, args.output, args.sep, args.encoding)