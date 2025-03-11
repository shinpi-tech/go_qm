package qm

import (
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSearch(t *testing.T) {
	oid, _ := primitive.ObjectIDFromHex("6302ac8a85bafafe377bd7dd")

	dateStrFrom := "12-12-2024"
	dateStrTo := "29-12-2024"
	tFrom, _ := time.Parse("02-01-2006", dateStrFrom)
	tFrom = tFrom.UTC()
	tTo, _ := time.Parse("02-01-2006", dateStrTo)
	tTo = tTo.UTC()

	tests := []struct {
		name     string
		query    map[string]string
		expected Query
	}{
		{
			name: "Сортировка по убыванию цены",
			query: map[string]string{
				"sort": "-price",
			},
			expected: Query{
				Match: bson.M{},
				Sort:  bson.D{{Key: "price", Value: -1}, {Key: "_id", Value: 1}},
				Page:  1,
				Limit: 30,
			},
		},
		{
			name: "Сортировка по увеличению цены",
			query: map[string]string{
				"sort": "price",
			},
			expected: Query{
				Match: bson.M{},
				Sort:  bson.D{{Key: "price", Value: 1}, {Key: "_id", Value: 1}},
				Page:  1,
				Limit: 30,
			},
		},
		{
			name: "Постраничная навигация",
			query: map[string]string{
				"limit": "50",
				"page":  "2",
			},
			expected: Query{
				Match: bson.M{},
				Sort:  bson.D{{Key: "_id", Value: 1}},
				Limit: 50,
				Page:  2,
			},
		},
		{
			name: "Фильтрация по цене",
			query: map[string]string{
				"price": "100~200",
			},
			expected: Query{
				Match: bson.M{"price": bson.M{"$gte": 100, "$lte": 200}},
				Sort:  bson.D{{Key: "_id", Value: 1}},
				Page:  1,
				Limit: 30,
			},
		},
		{
			name: "Фильтрация по цене от",
			query: map[string]string{
				"price": "100~",
			},
			expected: Query{
				Match: bson.M{"price": bson.M{"$gte": 100}},
				Sort:  bson.D{{Key: "_id", Value: 1}},
				Page:  1,
				Limit: 30,
			},
		},
		{
			name: "Фильтрация по ID",
			query: map[string]string{
				"id": "6302ac8a85bafafe377bd7dd",
			},
			expected: Query{
				Match: bson.M{"_id": oid},
				Sort:  bson.D{{Key: "_id", Value: 1}},
				Page:  1,
				Limit: 30,
			},
		},
		{
			name: "Многочисленная фильтрация по ID",
			query: map[string]string{
				"brand": "6302ac8a85bafafe377bd7dd,6302ac8a85bafafe377bd7dd",
			},
			expected: Query{
				Match: bson.M{"brand": bson.M{"$in": []any{oid, oid}}},
				Sort:  bson.D{{Key: "_id", Value: 1}},
				Page:  1,
				Limit: 30,
			},
		},
		{
			name: "Филтрация по интервалу дат",
			query: map[string]string{
				"params.load_index": "12-12-2024~!29-12-2024",
			},
			expected: Query{
				Match: bson.M{"params.load_index": bson.M{"$gte": tFrom, "$lt": tTo}},
				Sort:  bson.D{{Key: "_id", Value: 1}},
				Page:  1,
				Limit: 30,
			},
		},
		{
			name: "Фильтрация по нескольким параметрам",
			query: map[string]string{
				"params.width": "test1,test2,test3",
			},
			expected: Query{
				Match: bson.M{"params.width": bson.M{"$in": []string{"test1", "test2", "test3"}}},
				Sort:  bson.D{{Key: "_id", Value: 1}},
				Page:  1,
				Limit: 30,
			},
		},
		{
			name: "Фильтрация по булевому значению",
			query: map[string]string{
				"available": "true",
			},
			expected: Query{
				Match: bson.M{"available": true},
				Sort:  bson.D{{Key: "_id", Value: 1}},
				Page:  1,
				Limit: 30,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Search(tt.query)
			if err != nil {
				t.Fatalf("Search() вернул ошибку: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Ожидается %+v", tt.expected)
				t.Errorf("Результат %+v", result)
			}
		})
	}
}
