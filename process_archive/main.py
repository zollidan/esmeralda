"""
analyze_columns.py

Анализирует все CSV файлы в папке:
- Собирает полный набор колонок по всем файлам
- Показывает какие колонки есть в каждом файле
- Выводит отчёт в новый CSV файл

Использование:
    python analyze_columns.py --folder ./data --output report.csv
    python analyze_columns.py  # по умолчанию ищет CSV в текущей папке
"""

import argparse
import glob
import os
import csv
import pandas as pd
from collections import defaultdict


# ─── Настройки ────────────────────────────────────────────────────────────────

SEPARATOR   = ";"          # разделитель в исходных CSV
ENCODING    = "cp1251"     # кодировка (cp1251 для кириллицы из Excel)
SKIP_ROWS   = 8            # строк мусорной шапки сверху (подбери под свои файлы)
MAX_HEADER_ROWS = 3        # сколько строк объединённой шапки читать


# ─── Чтение шапки файла ───────────────────────────────────────────────────────

def read_columns(filepath: str) -> list[str]:
    """
    Читает колонки из CSV файла.
    Пробует несколько вариантов skiprows если стандартный не даёт результата.
    """
    for skip in [SKIP_ROWS, 0, 1, 2, 3]:
        try:
            df = pd.read_csv(
                filepath,
                sep=SEPARATOR,
                skiprows=skip,
                nrows=0,
                encoding=ENCODING,
                low_memory=False,
            )
            cols = [c for c in df.columns if not c.startswith("Unnamed")]
            if len(cols) > 3:
                return cols
        except Exception:
            continue
    return []


# ─── Нормализация имён колонок ────────────────────────────────────────────────

def normalize(name: str) -> str:
    """Убирает лишние пробелы, переводит в нижний регистр для сравнения."""
    return name.strip().lower()


# ─── Основной анализ ──────────────────────────────────────────────────────────

def analyze(folder: str, output: str):
    pattern = os.path.join(folder, "*.csv")
    files = sorted(glob.glob(pattern))

    if not files:
        print(f"[!] CSV файлы не найдены в папке: {folder}")
        return

    print(f"[+] Найдено файлов: {len(files)}")

    # Имя файла → список колонок
    file_columns: dict[str, list[str]] = {}

    for filepath in files:
        fname = os.path.basename(filepath)
        cols = read_columns(filepath)
        file_columns[fname] = cols
        print(f"    {fname}: {len(cols)} колонок")

    # Полный набор (union) — нормализованное имя → оригинальные варианты
    all_norm: dict[str, list[str]] = defaultdict(list)
    for fname, cols in file_columns.items():
        for col in cols:
            norm = normalize(col)
            if col not in all_norm[norm]:
                all_norm[norm].append(col)

    all_norm_keys = sorted(all_norm.keys())
    filenames = sorted(file_columns.keys())

    print(f"\n[+] Всего уникальных колонок: {len(all_norm_keys)}")

    # Считаем в скольких файлах встречается каждая колонка
    presence: dict[str, dict[str, bool]] = {}
    for norm_col in all_norm_keys:
        presence[norm_col] = {}
        for fname in filenames:
            file_norms = {normalize(c) for c in file_columns[fname]}
            presence[norm_col][fname] = norm_col in file_norms

    count_present = {
        norm_col: sum(presence[norm_col].values())
        for norm_col in all_norm_keys
    }

    common_cols = [c for c in all_norm_keys if count_present[c] == len(files)]
    partial_cols = [c for c in all_norm_keys if 0 < count_present[c] < len(files)]
    unique_cols = [c for c in all_norm_keys if count_present[c] == 1]

    print(f"    Есть во ВСЕХ файлах:     {len(common_cols)}")
    print(f"    Есть в части файлов:     {len(partial_cols)}")
    print(f"    Только в одном файле:    {len(unique_cols)}")

    # ─── Запись отчёта ────────────────────────────────────────────────────────

    with open(output, "w", newline="", encoding="utf-8-sig") as f:
        writer = csv.writer(f, delimiter=";")

        # Шапка
        header = [
            "колонка_норм",
            "оригинальные_варианты",
            "кол_во_файлов",
            "во_всех",
            "только_в_одном",
        ] + filenames
        writer.writerow(header)

        # Строки — сначала общие, потом частичные, потом уникальные
        def write_group(cols, label):
            if cols:
                writer.writerow([f"=== {label} ({len(cols)}) ==="])
            for norm_col in cols:
                variants = " | ".join(all_norm[norm_col])
                in_all = "да" if count_present[norm_col] == len(files) else ""
                only_one = "да" if count_present[norm_col] == 1 else ""
                row = [
                    norm_col,
                    variants,
                    count_present[norm_col],
                    in_all,
                    only_one,
                ]
                for fname in filenames:
                    row.append("✓" if presence[norm_col][fname] else "")
                writer.writerow(row)

        write_group(common_cols,  "ЕСТЬ ВО ВСЕХ ФАЙЛАХ")
        write_group(partial_cols, "ЕСТЬ В ЧАСТИ ФАЙЛОВ")
        write_group(unique_cols,  "ТОЛЬКО В ОДНОМ ФАЙЛЕ")

    print(f"\n[✓] Отчёт сохранён: {output}")
    print(f"\n--- Колонки общие для всех файлов ---")
    for c in common_cols:
        print(f"  {all_norm[c][0]}")


# ─── CLI ──────────────────────────────────────────────────────────────────────

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Анализ колонок CSV файлов")
    parser.add_argument(
        "--folder", default=".", help="Папка с CSV файлами (по умолчанию: .)"
    )
    parser.add_argument(
        "--output", default="columns_report.csv", help="Имя выходного файла"
    )
    parser.add_argument(
        "--sep", default=SEPARATOR, help=f"Разделитель (по умолчанию: '{SEPARATOR}')"
    )
    parser.add_argument(
        "--encoding", default=ENCODING, help=f"Кодировка (по умолчанию: {ENCODING})"
    )
    parser.add_argument(
        "--skiprows", type=int, default=SKIP_ROWS,
        help=f"Строк шапки пропустить (по умолчанию: {SKIP_ROWS})"
    )
    args = parser.parse_args()

    SEPARATOR = args.sep
    ENCODING  = args.encoding
    SKIP_ROWS = args.skiprows

    analyze(args.folder, args.output)