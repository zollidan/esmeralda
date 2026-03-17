import pandas as pd

df = pd.read_csv('../archive/csv/2026/календарь.csv', delimiter=";")

print(df.head(5))