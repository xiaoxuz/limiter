package limiter

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestLB(t *testing.T) {
	lb := NewLeakyBucket(&LbConfig{
		Rate:     1,
		MaxSlack: 0,
	})

	var handler = func() {
		var pass, unpass int64
		pass = 0
		for i := 0; i < 12; i++ {
			if err := lb.Take(); err != nil {
				atomic.AddInt64(&unpass, 1)
			} else {
				atomic.AddInt64(&pass, 1)
			}
			time.Sleep(100 * time.Millisecond)
		}
		t.Logf("pass:[%d] unpass:[%d]", pass, unpass)
	}

	handler()
	t.Logf("available cnt:%d", lb.Cnt())
	time.Sleep(5 * time.Second)
	t.Logf("available cnt:%d", lb.Cnt())

	handler()
}
