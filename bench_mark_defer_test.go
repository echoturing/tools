package tools

import "testing"

func BenchmarkDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		defers()
	}
}

func BenchmarkNormal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		normal()
	}
}
