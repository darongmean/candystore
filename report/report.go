package report

import (
	"sort"
	"strconv"
)

type TopCustomer struct {
	Name           string `json:"name"`
	FavouriteSnack string `json:"favouriteSnack"`
	TotalSnacks    uint64 `json:"totalSnacks"`
}

func ListTopCustomers(records [][]string) ([]TopCustomer, error) {
	eatenByNameAndCandy, eatenByName, maxEatenByName, err := aggregateEaten(dropHeaderRow(records))
	if err != nil {
		return nil, err
	}

	customers := findFavouriteSnack(eatenByNameAndCandy, eatenByName, maxEatenByName)

	sort.Slice(customers, func(i, j int) bool {
		return customers[i].TotalSnacks > customers[j].TotalSnacks
	})

	return customers, nil
}

func dropHeaderRow(records [][]string) [][]string {
	return records[1:]
}

func aggregateEaten(records [][]string) (eatenByNameAndCandy map[string]map[string]uint64, eatenByName map[string]uint64, maxEatenByName map[string]uint64, err error) {
	eatenByNameAndCandy = make(map[string]map[string]uint64)
	eatenByName = make(map[string]uint64)
	maxEatenByName = make(map[string]uint64)

	for _, record := range records {
		name := record[0]
		candy := record[1]
		eaten, err := strconv.ParseUint(record[2], 10, 32)
		if err != nil {
			return nil, nil, nil, err
		}

		eatenByName[name] += eaten

		if eatenByNameAndCandy[name] == nil {
			eatenByNameAndCandy[name] = make(map[string]uint64)
		}
		eatenByNameAndCandy[name][candy] += eaten

		if maxEatenByName[name] < eatenByNameAndCandy[name][candy] {
			maxEatenByName[name] = eatenByNameAndCandy[name][candy]
		}
	}

	return eatenByNameAndCandy, eatenByName, maxEatenByName, nil
}

func findFavouriteSnack(eatenByNameAndCandy map[string]map[string]uint64, eatenByName map[string]uint64, maxEatenByName map[string]uint64) []TopCustomer {
	customers := make([]TopCustomer, 0)

	for name, candyEaten := range eatenByNameAndCandy {
		for candy, eaten := range candyEaten {
			if maxEatenByName[name] == eaten {
				customers = append(customers, TopCustomer{
					Name:           name,
					FavouriteSnack: candy,
					TotalSnacks:    eatenByName[name]})
			}
		}
	}

	return customers
}
