package tools

import "sync"

// MergeCh that help us merge many channes in one chan
func MergeChs(chs ...<-chan interface{}) <-chan interface{} {
	merged := make(chan interface{})
	var wg sync.WaitGroup
	wg.Add(len(chs))
	go func() {
		wg.Wait()
		close(merged)
	}()
	for _, ch := range chs {
		go func(c <-chan interface{}) {
			for v := range c {
				merged <- v
			}
			wg.Done()
		}(ch)
	}
	return merged
}
