import argparse
import csv
import sys
from pathlib import Path


def read_table(path: str) -> tuple[list[str], list[dict]]:
    p = Path(path)
    if not p.exists():
        print(f"File not found: {path}", file=sys.stderr)
        sys.exit(1)

    suffix = p.suffix.lower()

    if suffix == ".csv":
        with open(p, newline="", encoding="utf-8") as f:
            reader = csv.DictReader(f)
            rows = list(reader)
            cols = list(reader.fieldnames or [])
        return cols, rows

    if suffix in (".xls", ".xlsx"):
        try:
            import openpyxl
        except ImportError:
            print("openpyxl required for Excel files: pip install openpyxl", file=sys.stderr)
            sys.exit(1)
        wb = openpyxl.load_workbook(p, read_only=True, data_only=True)
        ws = wb.active
        rows_raw = list(ws.iter_rows(values_only=True))
        if not rows_raw:
            return [], []
        cols = [str(c) if c is not None else "" for c in rows_raw[0]]
        rows = [dict(zip(cols, [str(v) if v is not None else "" for v in row])) for row in rows_raw[1:]]
        return cols, rows

    print(f"Unsupported file format: {suffix} (supported: .csv, .xlsx, .xls)", file=sys.stderr)
    sys.exit(1)


def compare_columns(cols_a: list[str], cols_b: list[str]) -> dict:
    set_a = set(cols_a)
    set_b = set(cols_b)
    return {
        "only_in_a": sorted(set_a - set_b),
        "only_in_b": sorted(set_b - set_a),
        "common": [c for c in cols_a if c in set_b],
    }


def compare_rows(rows_a: list[dict], rows_b: list[dict], common_cols: list[str], key_col: str | None) -> dict:
    if key_col:
        index_a = {row[key_col]: row for row in rows_a if key_col in row}
        index_b = {row[key_col]: row for row in rows_b if key_col in row}
        keys_a = set(index_a)
        keys_b = set(index_b)

        only_in_a = sorted(keys_a - keys_b)
        only_in_b = sorted(keys_b - keys_a)
        common_keys = keys_a & keys_b

        diffs = []
        for key in sorted(common_keys):
            row_a = index_a[key]
            row_b = index_b[key]
            col_diffs = []
            for col in common_cols:
                val_a = row_a.get(col, "")
                val_b = row_b.get(col, "")
                if val_a != val_b:
                    col_diffs.append((col, val_a, val_b))
            if col_diffs:
                diffs.append((key, col_diffs))

        return {
            "mode": "keyed",
            "key_col": key_col,
            "only_in_a": only_in_a,
            "only_in_b": only_in_b,
            "diffs": diffs,
            "total_a": len(rows_a),
            "total_b": len(rows_b),
            "matched": len(common_keys),
        }
    else:
        min_len = min(len(rows_a), len(rows_b))
        diffs = []
        for i in range(min_len):
            col_diffs = []
            for col in common_cols:
                val_a = rows_a[i].get(col, "")
                val_b = rows_b[i].get(col, "")
                if val_a != val_b:
                    col_diffs.append((col, val_a, val_b))
            if col_diffs:
                diffs.append((i + 1, col_diffs))

        return {
            "mode": "positional",
            "diffs": diffs,
            "total_a": len(rows_a),
            "total_b": len(rows_b),
            "matched": min_len,
        }


def print_report(file_a: str, file_b: str, col_cmp: dict, row_cmp: dict, show_values: bool):
    sep = "-" * 60

    print(sep)
    print("TABLE COMPARISON REPORT")
    print(sep)
    print(f"Table A: {file_a}")
    print(f"Table B: {file_b}")
    print()

    print("COLUMNS")
    print(sep)
    print(f"  Common columns   : {len(col_cmp['common'])}")
    print(f"  Only in A        : {len(col_cmp['only_in_a'])}")
    print(f"  Only in B        : {len(col_cmp['only_in_b'])}")
    if col_cmp["only_in_a"]:
        print(f"    A extra: {', '.join(col_cmp['only_in_a'])}")
    if col_cmp["only_in_b"]:
        print(f"    B extra: {', '.join(col_cmp['only_in_b'])}")
    print()

    print("ROWS")
    print(sep)
    print(f"  Rows in A        : {row_cmp['total_a']}")
    print(f"  Rows in B        : {row_cmp['total_b']}")

    if row_cmp["mode"] == "keyed":
        print(f"  Key column       : {row_cmp['key_col']}")
        print(f"  Matched by key   : {row_cmp['matched']}")
        print(f"  Only in A        : {len(row_cmp['only_in_a'])}")
        print(f"  Only in B        : {len(row_cmp['only_in_b'])}")
        if row_cmp["only_in_a"]:
            print(f"    A only keys: {', '.join(str(k) for k in row_cmp['only_in_a'][:20])}")
            if len(row_cmp["only_in_a"]) > 20:
                print(f"    ... and {len(row_cmp['only_in_a']) - 20} more")
        if row_cmp["only_in_b"]:
            print(f"    B only keys: {', '.join(str(k) for k in row_cmp['only_in_b'][:20])}")
            if len(row_cmp["only_in_b"]) > 20:
                print(f"    ... and {len(row_cmp['only_in_b']) - 20} more")
    else:
        print(f"  Compared rows    : {row_cmp['matched']}")
        if row_cmp["total_a"] != row_cmp["total_b"]:
            print(f"  Row count mismatch: A has {row_cmp['total_a']}, B has {row_cmp['total_b']}")
    print()

    diffs = row_cmp["diffs"]
    print("DIFFERENCES")
    print(sep)
    print(f"  Rows with diffs  : {len(diffs)}")
    print()

    if diffs and show_values:
        for entry in diffs:
            if row_cmp["mode"] == "keyed":
                key, col_diffs = entry
                print(f"  Key={key}")
            else:
                row_num, col_diffs = entry
                print(f"  Row {row_num}")
            for col, val_a, val_b in col_diffs:
                print(f"    [{col}] A={repr(val_a)}  B={repr(val_b)}")
            print()

    if not diffs:
        print("  No differences found in compared rows and columns.")
    elif not show_values:
        print("  Run with --values to see per-cell differences.")
    print(sep)


def main():
    parser = argparse.ArgumentParser(
        description="Compare two tables (CSV or Excel) and print a diff report.",
    )
    parser.add_argument("table_a", help="Path to first table (CSV or XLSX)")
    parser.add_argument("table_b", help="Path to second table (CSV or XLSX)")
    parser.add_argument(
        "--key",
        metavar="COLUMN",
        help="Column name to use as a unique key for row matching",
    )
    parser.add_argument(
        "--values",
        action="store_true",
        help="Show per-cell value differences in the report",
    )
    args = parser.parse_args()

    cols_a, rows_a = read_table(args.table_a)
    cols_b, rows_b = read_table(args.table_b)

    if args.key and args.key not in cols_a:
        print(f"Key column '{args.key}' not found in {args.table_a}", file=sys.stderr)
        sys.exit(1)
    if args.key and args.key not in cols_b:
        print(f"Key column '{args.key}' not found in {args.table_b}", file=sys.stderr)
        sys.exit(1)

    col_cmp = compare_columns(cols_a, cols_b)
    row_cmp = compare_rows(rows_a, rows_b, col_cmp["common"], args.key)

    print_report(args.table_a, args.table_b, col_cmp, row_cmp, args.values)


if __name__ == "__main__":
    main()
