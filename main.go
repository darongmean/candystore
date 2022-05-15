package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type Favourite struct {
	Name           string `json:"name"`
	FavouriteSnack string `json:"favouriteSnack"`
	TotalSnacks    uint64 `json:"totalSnacks"`
}

func main() {
	f, err := os.Open("test/data.csv")
	check(err)
	defer f.Close()

	r := newCSVReader(f)
	records, err := r.ReadAll()
	check(err)

	records = dropHeaderRow(records)

	totalSnackByName := make(map[string]uint64)
	totalEaten := make(map[string]map[string]uint64)
	maxEatenByName := make(map[string]uint64)
	for _, record := range records {
		name := record[0]
		candy := record[1]
		eaten, err := strconv.ParseUint(record[2], 10, 32)
		check(err)

		totalSnackByName[name] += eaten

		if totalEaten[name] == nil {
			totalEaten[name] = make(map[string]uint64)
		}
		totalEaten[name][candy] += eaten

		if maxEatenByName[name] < totalEaten[name][candy] {
			maxEatenByName[name] = totalEaten[name][candy]
		}
	}

	favourites := make([]Favourite, 0)
	for name, candyEaten := range totalEaten {
		for candy, eaten := range candyEaten {
			if maxEatenByName[name] == eaten {
				favourites = append(favourites,
					Favourite{Name: name,
						FavouriteSnack: candy,
						TotalSnacks:    totalSnackByName[name]})
			}
		}
	}

	sort.Slice(favourites, func(i, j int) bool {
		return favourites[i].TotalSnacks > favourites[j].TotalSnacks
	})

	str, err := json.Marshal(favourites)
	check(err)

	fmt.Print(string(str))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func newCSVReader(f *os.File) *csv.Reader {
	r := csv.NewReader(f)

	r.Comma = '\t'
	r.Comment = '#'
	r.TrimLeadingSpace = true

	return r
}

func dropHeaderRow(records [][]string) [][]string {
	return records[1:]
}
