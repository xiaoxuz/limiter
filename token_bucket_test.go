package limiter

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	tb := NewTokenBucket(&TbConfig{
		QPS:    10,
		MaxCap: 20,
	})
	var handler = func(cnt int) {
		var wg sync.WaitGroup
		var pass, unpass int64
		wg.Add(cnt)
		for i := 0; i < cnt; i++ {
			go func() {
				defer wg.Done()
				if err := tb.Take(1); err != nil {
					atomic.AddInt64(&unpass, 1)
				} else {
					atomic.AddInt64(&pass, 1)
				}
			}()

		}
		wg.Wait()
		t.Logf("cnt:[%d] qps:[%d] maxcap:[%d] pass:[%d] unpass:[%d]", cnt, 10, 20, pass, unpass)
	}

	handler(50)

	time.Sleep(2 * time.Second)

	handler(100)
}
