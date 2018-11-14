package tools

import (
	"sync/atomic"
	"testing"
)

func TestWaitGroupWrapper_Wrap(t *testing.T) {
	wgw := WaitGroupWrapper{}
	var s int32 = 0
	for i := 0; i < 100; i++ {
		wgw.Wrap(func() {
			atomic.AddInt32(&s, 1)
		})
	}
	wgw.Wait()
	if s != 100 {
		t.Error("need 100")
	}

}
