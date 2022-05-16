package report

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListTopCustomers(t *testing.T) {
	for _, tc := range validRecordsTestCases {
		actual, err := ListTopCustomers(tc.records)

		assert.Nil(t, err)
		assert.ElementsMatch(t, tc.expected, actual)
	}
}

func TestListTopCustomersSortOutput(t *testing.T) {
	for _, tc := range validRecordsTestCases {
		actual, _ := ListTopCustomers(tc.records)

		// empty slice is sorted
		if len(actual) == 0 {
			continue
		}

		head := actual[0]
		for _, next := range actual[1:] {
			assert.True(t, head.TotalSnacks >= next.TotalSnacks, "Output is not sorted")
			head = next
		}
	}
}

func TestListTopCustomersAssumeEatenIsInteger(t *testing.T) {
	records := [][]string{
		{"Name", "Candy", "Eaten"},
		{"Jane", "Nötchoklad", "not-integer"},
	}

	_, err := ListTopCustomers(records)

	assert.NotNil(t, err)
}

var validRecordsTestCases = []struct {
	records  [][]string
	expected []TopCustomer
}{
	{nil, []TopCustomer{}},
	{[][]string{{"Name", "Candy", "Eaten"}}, []TopCustomer{}},
	{ // general output
		[][]string{
			{"Name", "Candy", "Eaten"},
			{"Jane", "Nötchoklad", "1"},
			{"Annika", "Center", "2"},
			{"Jane", "Kexchoklad", "3"}},
		[]TopCustomer{
			{Name: "Jane", FavouriteSnack: "Kexchoklad", TotalSnacks: 4},
			{Name: "Annika", FavouriteSnack: "Center", TotalSnacks: 2},
		},
	},
	{ // multiple favourites output
		[][]string{
			{"Name", "Candy", "Eaten"},
			{"Annika", "Center", "2"},
			{"Annika", "Kexchoklad", "2"},
		},
		[]TopCustomer{
			{Name: "Annika", FavouriteSnack: "Center", TotalSnacks: 4},
			{Name: "Annika", FavouriteSnack: "Kexchoklad", TotalSnacks: 4},
		},
	},
	{ // same total snacks output
		[][]string{
			{"Name", "Candy", "Eaten"},
			{"Jane", "Nötchoklad", "1"},
			{"Annika", "Center", "2"},
			{"Annika", "Kexchoklad", "2"},
			{"Jane", "Kexchoklad", "3"},
		},
		[]TopCustomer{
			{Name: "Jane", FavouriteSnack: "Kexchoklad", TotalSnacks: 4},
			{Name: "Annika", FavouriteSnack: "Center", TotalSnacks: 4},
			{Name: "Annika", FavouriteSnack: "Kexchoklad", TotalSnacks: 4},
		},
	},
}
