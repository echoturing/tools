package utils

import (
	"fmt"
	"testing"
)

func TestInsertOrUpdate(t *testing.T) {
	value, err := InsertOrUpdate(&InsertOrUpdateBody{
		Table:                  "item",
		Columns:                []string{"a", "b"},
		DuplicateUpdateColumns: []string{"a"},
	})
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(value)
}
