#!/bin/bash

set -e
set -o nounset

mkdir -p harvest
for src in ~/Dropbox/self/Farm/20*harvest.numbers; do
  y="$(basename "$src" | awk '{print $1}')"
  dst="harvest/$y.csv"
  if ! (set -x; uv run cat-numbers -b -s 'Sheet 1' -t 'Table 1' "$src" > "$dst"); then
    rm -f "$dst"
  fi
done
