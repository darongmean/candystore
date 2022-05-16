# Usage

Install with the following command.

> go install

Then run one of the following command to get started.

> $GOPAHT/bin/candystore 

> $GOPAHT/bin/candystore help

> $GOPAHT/bin/candystore fav --data test/data.csv

# Objective

Given the following data

```csv
Name	Candy	Eaten
Annika	Geisha	100
Jonas	Geisha	200
Jonas	Kexchoklad	100
Aadya	Nötchoklad	2
Jonas	Nötchoklad	3
Jane	Nötchoklad	17
Annika	Geisha	100
Jonas	Geisha	700
Jane	Nötchoklad	4
Aadya	Center	7
Jonas	Geisha	900
Jane	Nötchoklad	1
Jonas	Kexchoklad	12
Jonas	Plopp	40
Jonas	Center	27
Aadya	Center	2
Annika	Center	8
```

Output result like

```json
[
    {
        "name": "Jonas",
        "favouriteSnack": "Geisha",
        "totalSnacks": 1982
    },
    {
        "name": "Annika",
        "favouriteSnack": "Geisha",
        "totalSnacks": 208
    },
    {
        "name": "Jane",
        "favouriteSnack": "Nötchoklad",
        "totalSnacks": 22
    },
    {
        "name": "Aadya",
        "favouriteSnack": "Center",
        "totalSnacks": 11
    }
]
```

# How

Output is a json of a list of objects. The object has the following properties
- `name` is the distinct value of `name` column.
- `totalSnacks` can be computed from grouping input data by `name` column and sum the `eaten` column.
- `favouriteSnack` is the distinct value of `candy` column. 
  It can be computed from grouping input data by `name` and `candy` column, sum the `eaten` column, and find the max value.

Thus to compute the output, we need to build two tables:

```csv
Name    TotalEaten
Jane    22
```

and

```csv
Name    Candy   TotalEaten
Jane    Geisha  15
```

Then we can join the two table by `name` and filter by max value per candy.
