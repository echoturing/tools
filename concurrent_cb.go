package tools

import "sync"

//ConcurrentCb call cb from start to start+concurrentCount
func ConcurrentCb(start, concurrentCount int, cb func(i int)) {
	for index := 0; index < concurrentCount; index++ {
		current := index + start
		go func() {
			cb(current)
		}()
	}
}

func concurrentCbWithWg(start, concurrentCount int, cb func(i int), wg *sync.WaitGroup) {
	for index := 0; index < concurrentCount; index++ {
		current := index + start
		go func() {
			defer wg.Done()
			cb(current)
		}()
	}
}

// do not need self control wait group
func SyncConcurrentCb(start, concurrentCount int, cb func(i int)) {
	wg := &sync.WaitGroup{}
	wg.Add(concurrentCount)
	concurrentCbWithWg(start, concurrentCount, cb, wg)
	wg.Wait()
}
