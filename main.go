package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

type Favourite struct {
	Name           string `json:"name"`
	FavouriteSnack string `json:"favouriteSnack"`
	TotalSnacks    uint64 `json:"totalSnacks"`
}

func main() {
	app := &cli.App{
		Name:  "candystore",
		Usage: "Candy store. - Example: candystore fav --data test/data.csv",
		Commands: []*cli.Command{
			{
				Name:  "fav",
				Usage: "List top customers and their favourties",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "data",
						Aliases:  []string{"d"},
						Usage:    "Path to a customer csv, tab delimited, file. The csv file should contain header: Name, Candy, and Eaten.",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					fileName, err := parseFavouriteCommand(c)
					if err != nil {
						return err
					}
					execFavouriteCommand(fileName)
					// TODO: return error
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	check(err)
}

func parseFavouriteCommand(c *cli.Context) (fileName string, err error) {
	fileName = strings.TrimSpace(c.String("data"))
	fmt.Printf("candystore fav --data %s\n", fileName)

	if len(fileName) == 0 {
		return "", errors.New("--data value should contain at least 1 non-whitespace character")
	} else {
		return fileName, nil
	}
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

func execFavouriteCommand(fileName string) {
	f, err := os.Open(fileName)
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
