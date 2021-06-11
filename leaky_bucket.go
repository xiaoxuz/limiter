package limiter

import (
	"sync"
	"time"
)

type LbConfig struct {
	Rate     float64
	MaxSlack int64
}

type LeakyBucket struct {
	*LbConfig
	m          sync.Mutex
	perRequest time.Duration
	bufferTime time.Duration
	slackTime  time.Duration
	lastTime   time.Time
}

var _ Limiter = &LeakyBucket{}

func NewLeakyBucket(c *LbConfig) Limiter {
	lb := &LeakyBucket{
		LbConfig:   c,
		m:          sync.Mutex{},
		bufferTime: 0,
		lastTime:   time.Time{},
	}
	// 速率
	lb.perRequest = time.Second / time.Duration(c.Rate)
	// 松弛度
	lb.slackTime = (^time.Duration(c.MaxSlack) * time.Second / time.Duration(c.Rate)) + 1
	return lb
}

func (lb *LeakyBucket) Take() error {
	lb.m.Lock()
	defer lb.m.Unlock()

	n := time.Now()

	if lb.lastTime.Second() == 0 {
		lb.lastTime = n
		return nil
	}

	if lb.MaxSlack > 0 {
		return lb.withSlack()
	} else {
		return lb.withoutSlack()
	}
}

func (lb *LeakyBucket) withoutSlack() error {
	n := time.Now()
	lb.bufferTime = lb.perRequest - n.Sub(lb.lastTime)
	// 请求间隔时间+buffer时间 < 速率, 拒绝服务
	if lb.bufferTime > 0 {
		return ErrNoTEnoughToken
	} else {
		lb.lastTime = n
	}
	return nil
}

func (lb *LeakyBucket) withSlack() error{
	n := time.Now()
	lb.bufferTime += lb.perRequest - n.Sub(lb.lastTime)

	// 请求间隔时间+buffer时间 < 速率, 拒绝服务
	if lb.bufferTime > 0 {
		return ErrNoTEnoughToken
	} else {
		lb.lastTime = n
	}

	// 允许抵消的最长时间
	if lb.bufferTime < lb.slackTime {
		lb.bufferTime = lb.slackTime
	}
	return nil
}

func (lb *LeakyBucket) Cnt() int64 {
	n := time.Now()
	bufferTime := lb.bufferTime + lb.perRequest - n.Sub(lb.lastTime)

	return int64(bufferTime / lb.perRequest)
}
