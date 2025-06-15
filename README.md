# Harvest helper

These scripts help me compile data into a harvest spreadsheet.

Here are the steps:
1. Original record is a one-page per week list of crops and how much was harvested.
2. Copy the first column of the harvest spreadsheet and paste it into `crops.txt`. The expected input has two header rows, then one row per crop. There should be no escaping or quoting here.
3. Run `script/input --out weekNN.txt`. Go down the harvest sheet in order and type crop name (with handy autocomplete), hit enter, then type the quantity.
4. Run `script/process weekNN.txt | pbcopy`. 
5. Paste into the right week's column (below the header).
