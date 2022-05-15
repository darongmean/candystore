package report

import (
	"sort"
	"strconv"
)

type Favourite struct {
	Name           string `json:"name"`
	FavouriteSnack string `json:"favouriteSnack"`
	TotalSnacks    uint64 `json:"totalSnacks"`
}

func ListTopCustomers(records [][]string) ([]Favourite, error) {
	records = dropHeaderRow(records)

	totalSnackByName := make(map[string]uint64)
	totalEaten := make(map[string]map[string]uint64)
	maxEatenByName := make(map[string]uint64)
	for _, record := range records {
		name := record[0]
		candy := record[1]
		eaten, err := strconv.ParseUint(record[2], 10, 32)
		if err != nil {
			return nil, err
		}

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

	return favourites, nil
}

func dropHeaderRow(records [][]string) [][]string {
	return records[1:]
}
