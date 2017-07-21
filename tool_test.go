package tools

import (
	"testing"
)

//summary
func TestIsNil(t *testing.T) {
	var i interface{}
	if IsNil(i) != true {
		t.Error("test failed")
	}
	i = nil
	if IsNil(i) != true {
		t.Error("test failed")
	}

	//struct test
	type Ts struct {
	}
	var instanceTs Ts
	if IsNil(instanceTs) == true {
		t.Error("test failed")
	}
	var pointerTs *Ts
	if IsNil(pointerTs) != true {
		t.Error("test failed")
	}
	var pointerTsWithNew *Ts = new(Ts)
	if IsNil(pointerTsWithNew) == true {
		t.Error("test failed")
	}
	//slice test
	var sList []string
	if IsNil(sList) != true {
		t.Error("test failed")
	}
	var sListWithMake []string = make([]string, 0)
	if IsNil(sListWithMake) == true {
		t.Error("test failed")
	}
	//map test
	var sMap map[string]string
	if IsNil(sMap) != true {
		t.Error("test failed")
	}
	var sMapWithMake map[string]string = make(map[string]string)
	if IsNil(sMapWithMake) == true {
		t.Error("test failed")
	}
	//chan test
	var sChan chan string
	if IsNil(sChan) != true {
		t.Error("test failed")
	}
	var sChanWithMake chan string = make(chan string)
	if IsNil(sChanWithMake) == true {
		t.Error("test failed")
	}
	// struct never be nil
	var s struct{}
	if IsNil(s) == true {
		t.Error("test failed")
	}
}
