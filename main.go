package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/darongmean/candystore/report"

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

func execFavouriteCommand(fileName string) {
	f, err := os.Open(fileName)
	check(err)
	defer f.Close()

	r := newCSVReader(f)
	records, err := r.ReadAll()
	check(err)

	favourites, err := report.ListTopCustomers(records)
	check(err)

	str, err := json.Marshal(favourites)
	check(err)

	fmt.Print(string(str))
}
