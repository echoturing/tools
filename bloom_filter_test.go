package tools

import (
	"testing"
)

func TestStringHash_Hash(t *testing.T) {
	addr := "localhost:6379"
	password := ""
	db := 0
	keySuffix := "fade"
	rbc := NewRedisBloomFilter(addr, password, db, keySuffix)
	s1 := StringHash("StringHash1")
	var err error
	{
		err = rbc.AddItem(s1)
		if err != nil {
			t.Error("need nil")
		}
	}

}
