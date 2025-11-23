#!/usr/bin/env python3
"""
Extract the first table from a Numbers spreadsheet and convert it to CSV.
Performs sanity checks to ensure it's the correct table.
"""

import sys
import csv
import argparse
from pathlib import Path
from numbers_parser import Document


def check_header_rows(table, max_header_rows=2):
    """Check if the table has 1-2 header rows with expected content."""
    if len(table.rows()) < 1:
        return False, "Table has no rows"
    
    header_texts = []
    for i in range(min(max_header_rows, len(table.rows()))):
        row = table.rows()[i]
        row_text = " ".join(str(cell.value or "") for cell in row).lower()
        header_texts.append(row_text)
    
    # Check for expected header keywords
    has_week_number = any("week number" in text for text in header_texts)
    has_week_ending = any("week ending on" in text or "week ending" in text for text in header_texts)
    
    if not (has_week_number or has_week_ending):
        return False, f"Header rows don't contain 'week number' or 'week ending on'. Found: {header_texts}"
    
    return True, f"Header check passed: week_number={has_week_number}, week_ending={has_week_ending}"


def check_formulas_only_in_bottom_row(table):
    """Check that formulas only appear in the bottom row."""
    rows = list(table.rows())
    if len(rows) < 2:
        return True, "Table has fewer than 2 rows, skipping formula check"
    
    # Check all rows except the last one
    for row_idx, row in enumerate(rows[:-1]):
        for col_idx, cell in enumerate(row):
            if cell.formula is not None:
                return False, f"Formula found in row {row_idx + 1}, column {col_idx + 1} (expected only in bottom row)"
    
    # Check that the bottom row has at least one formula
    bottom_row = rows[-1]
    has_formula = any(cell.formula is not None for cell in bottom_row)
    
    if not has_formula:
        return True, "No formulas found in bottom row (this might be okay)"
    
    return True, "Formula check passed: formulas only in bottom row"


def extract_table_to_csv(numbers_file, output_csv):
    """Extract the first table from a Numbers file and save as CSV."""
    doc = Document(numbers_file)
    
    if len(doc.sheets) == 0:
        raise ValueError("Document has no sheets.")
    
    print(f"Found {len(doc.sheets)} sheets.")
    sheet = doc.sheets[0]
    
    print(f"Using sheet '{sheet.name}'.")
    if len(sheet.tables) == 0:
        raise ValueError("Sheet has no tables")
    
    print(f"Found {len(sheet.tables)} tables.")
    table = sheet.tables[0]
    
    print(f"Using table '{table.name}', with {len(table.rows())} rows.")
    
    # Perform sanity checks
    print("\nPerforming sanity checks...")
    
    # Check header rows
    header_ok, header_msg = check_header_rows(table)
    print(f"  Header check: {header_msg}")
    if not header_ok:
        raise ValueError(f"Header check failed: {header_msg}")
    
    # Check formulas
    formula_ok, formula_msg = check_formulas_only_in_bottom_row(table)
    print(f"  Formula check: {formula_msg}")
    if not formula_ok:
        raise ValueError(f"Formula check failed: {formula_msg}")
    
    print("\nAll sanity checks passed! Extracting data...")
    
    # Extract data and write to CSV
    with open(output_csv, 'w', newline='', encoding='utf-8') as csvfile:
        writer = csv.writer(csvfile)
        
        for row in table.rows():
            csv_row = []
            for cell in row:
                # Get the value, handling formulas by getting their computed value
                if cell.formula is not None:
                    # For formulas, use the value (computed result)
                    value = cell.value
                else:
                    value = cell.value
                
                # Convert to string, handling None
                csv_row.append(str(value) if value is not None else "")
            
            writer.writerow(csv_row)
    
    print(f"Successfully extracted {len(table.rows())} rows to {output_csv}")


def main():
    parser = argparse.ArgumentParser(
        description="Extract the first table from a Numbers spreadsheet to CSV"
    )
    parser.add_argument(
        "numbers_file",
        type=str,
        help="Path to the Numbers file (e.g., '2025 harvest.numbers')"
    )
    parser.add_argument(
        "-o", "--output",
        type=str,
        help="Output CSV file path (default: same as input with .csv extension)"
    )
    
    args = parser.parse_args()
    
    numbers_path = Path(args.numbers_file)
    if not numbers_path.exists():
        print(f"Error: File not found: {numbers_path}", file=sys.stderr)
        sys.exit(1)
    
    if args.output:
        output_path = Path(args.output)
    else:
        output_path = numbers_path.with_suffix('.csv')
    
    extract_table_to_csv(str(numbers_path), str(output_path))


if __name__ == "__main__":
    main()
