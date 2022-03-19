package ratelimit

import (
	"sync"
	"time"
)

type leakyRateLimit struct {
	requestPerSec int64
	interval      time.Duration
	last          time.Time
	mutex         sync.RWMutex
}

func NewLeakyRateLimit(requestPerSec int64) RateLimit {
	return &leakyRateLimit{
		requestPerSec: requestPerSec,
		interval:      time.Second / time.Duration(requestPerSec),
	}
}

func (l *leakyRateLimit) Take() time.Time {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	now := time.Now()
	if l.last.IsZero() {
		l.last = now
		return now
	}
	sleepFor := l.interval - now.Sub(l.last)
	if sleepFor > 0 {
		time.Sleep(sleepFor)
		l.last = now.Add(sleepFor)
	} else {
		l.last = now
	}
	return l.last
}
