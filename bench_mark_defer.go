package tools

//go:noinline
func defers() (r int) {
	defer func() {
		r = 42
	}()
	return
}

func normal() (r int) {
	r = 42
	return
}
