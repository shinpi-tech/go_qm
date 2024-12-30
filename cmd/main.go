package main

import (
	"fmt"

	qm "github.com/shinpi-tech/go_qm"
)

func main() {
	query := map[string]string{
		"limit":             "30",
		"sort":              "price,-name",
		"page":              "2",
		"available":         "false",
		"category":          "tyres",
		"price":             "!5000~10000",
		"date":              "27-12-2024~29-12-2024",
		"params.load_index": "~!96",
		"params.diameter":   "14,15",
		"params.width":      "175,195,205",
		"brand":             "6302ac8a85bafafe377bd7de,6302ac8a85bafafe377bd7dd",
	}

	res, err := qm.Search(query)

	if err != nil {
		panic(err)
	}

	fmt.Println("Match", res.Match)
	fmt.Println("Sort", res.Sort)
	fmt.Println("Page / Limit", res.Page, "/", res.Limit)
}
