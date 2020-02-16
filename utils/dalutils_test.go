package utils

import "testing"

func TestWrapWith(t *testing.T) {
	source := []string{
		"a",
	}
	dst := WrapWith(source, "`")
	if dst[0] != "`a`" {
		t.Error("wrap with error")
	}
}

func TestStringJoin(t *testing.T) {
	source := []string{
		"a", "b", "c",
	}
	res := StringJoin(source, ",")
	if res != "a,b,c" {
		t.Error("string join error")
	}
}
