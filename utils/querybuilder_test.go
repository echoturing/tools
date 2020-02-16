package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestInsertOrUpdate(t *testing.T) {
	query, values, err := InsertOrUpdate(NewInsertOrUpdateBody(
		"item",
		[]string{"aa", "b"},
		[]string{"aa"},
		struct {
			A int       `db:"aa"`
			B time.Time `db:"b"`
		}{2, time.Now()},
	))
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(query)
	fmt.Println(values)
}
