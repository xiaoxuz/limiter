package limiter

import (
	"sync"
	"time"
)

type TbConfig struct {
	QPS    int64
	MaxCap int64
}

type TokenBucket struct {
	*TbConfig
	m          sync.Mutex
	available  int64
	lastTime time.Time
}

var _ Limiter = &TokenBucket{}

func NewTokenBucket(c *TbConfig) Limiter {
	return &TokenBucket{
		TbConfig:   c,
		m:          sync.Mutex{},
		available:  c.QPS,
		lastTime: time.Now(),
	}
}

func (tb *TokenBucket) Take() error {
	tb.m.Lock()
	defer tb.m.Unlock()

	tb.fill()

	if 1 <= tb.available {
		tb.available -= 1
		return nil
	}

	return ErrNoTEnoughToken
}

func (tb *TokenBucket) Put(cnt int64) error {
	tb.m.Lock()
	defer tb.m.Unlock()

	tb.fill()

	tb.available += cnt
	if tb.MaxCap > 0 && tb.available > tb.MaxCap {
		tb.available = tb.MaxCap
	}
	return nil
}

func (tb *TokenBucket) Cnt() int64 {
	tb.m.Lock()
	defer tb.m.Unlock()

	tb.fill()

	return tb.available
}

func (tb *TokenBucket) fill() error {
	n := time.Now()
	timeUnit := n.Sub(tb.lastTime).Seconds()

	fillCnt := int64(timeUnit) * tb.QPS
	if fillCnt <= 0 {
		return nil
	}

	tb.available += fillCnt

	if tb.MaxCap > 0 && tb.available > tb.MaxCap {
		tb.available = tb.MaxCap
	}

	tb.lastTime = n
	return nil
}
