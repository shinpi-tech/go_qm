package qm

import (
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func processSort(sortValue string, res *Query) {
	sortParams := strings.Split(sortValue, ",")
	for _, param := range sortParams {
		if strings.HasPrefix(param, "-") {
			field := strings.TrimPrefix(param, "-")
			res.Sort[field] = -1
		} else {
			res.Sort[param] = 1
		}
	}
}

func processPagination(limitStr, pageStr string, res *Query) {
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return
	}

	res.Page = page
	res.Limit = limit
}

func processMatch(key, value string, res *Query) {
	if value == "true" || value == "false" { // bool
		res.Match[key] = (value == "true")
	} else if strings.Contains(value, "~") { // range (int / date)
		processRange(key, value, res)
	} else { // string / oid
		params := strings.Split(value, ",")
		result := make([]interface{}, 0, len(params))

		for _, param := range params {
			if oid, err := primitive.ObjectIDFromHex(param); err == nil {
				result = append(result, oid)
			} else {
				result = append(result, param)
			}
		}

		res.Match[key] = bson.M{"$in": result}
	}
}

func processRange(key, value string, res *Query) {
	values := strings.Split(value, "~")
	if len(values) != 2 {
		return
	}

	rangeQuery := bson.M{}

	processParam := func(param string, operator string) bson.M {
		if param == "" {
			return nil
		}

		inclusive := !strings.HasPrefix(param, "!")
		if !inclusive {
			param = strings.TrimPrefix(param, "!")
		}

		createCondition := func(value interface{}) bson.M {
			if inclusive {
				switch operator {
				case "$gt":
					return bson.M{"$gte": value}
				case "$lt":
					return bson.M{"$lte": value}
				}
			}
			return bson.M{operator: value}
		}

		// DATE
		if paramDate, err := time.Parse("02-01-2006", param); err == nil {
			return createCondition(paramDate)
		}
		// INT
		if paramInt, err := strconv.Atoi(param); err == nil {
			return createCondition(paramInt)
		}

		return nil
	}

	if gtInt := processParam(values[0], "$gt"); gtInt != nil {
		for k, v := range gtInt {
			rangeQuery[k] = v
		}
	}
	if ltInt := processParam(values[1], "$lt"); ltInt != nil {
		for k, v := range ltInt {
			rangeQuery[k] = v
		}
	}

	// Если есть хотя бы одно условие, добавляем его в Match
	if len(rangeQuery) > 0 {
		res.Match[key] = rangeQuery
	}
}
