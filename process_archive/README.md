uv add pandas

# базовый запуск — ищет CSV в текущей папке

uv run analyze_columns.py

# с указанием папки

uv run analyze_columns.py --folder ./data --output report.csv

# если шапка у тебя другой высоты

uv run analyze_columns.py --folder ./data --skiprows 3

```

**Что выдаёт на экран:**
```

[+] Найдено файлов: 8
file1.csv: 47 колонок
file2.csv: 52 колонок
...
[+] Всего уникальных колонок: 61
Есть во ВСЕХ файлах: 34
Есть в части файлов: 19
Только в одном файле: 8

--- Колонки общие для всех файлов ---
число
месяц
...
