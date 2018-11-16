package tools

import (
	"fmt"
	"sync"
	"testing"
)

func TestConcurrentCb(t *testing.T) {
	start := 1
	concurrentCount := 10
	wg := sync.WaitGroup{}
	wg.Add(concurrentCount)
	cb := func(i int) {
		fmt.Println(i)
		wg.Done()
	}
	ConcurrentCb(start, concurrentCount, cb)
	wg.Wait()
}

func TestSyncConcurrentCb(t *testing.T) {
	start := 1
	concurrentCount := 20
	cb := func(i int) {
		fmt.Println(i)
	}
	SyncConcurrentCb(start, concurrentCount, cb)

}
