package qm

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Query struct {
	Match bson.M
	Sort  bson.M
	Page  int
	Limit int
}

func Search(query map[string]string) (Query, error) {
	res := Query{
		Match: bson.M{},
		Sort:  bson.M{},
	}

	for k, v := range query {
		switch k {
		case "sort":
			if v == "" {
				continue
			}
			processSort(v, &res)
		case "limit":
			processPagination(v, query["page"], &res)
		case "page":
			continue
		default:
			processMatch(k, v, &res)
		}
	}

	return res, nil
}
