package ratelimit

import (
	"sync"
	"time"
)

// Convert burst traffic to the fixed QPS
// Pros: smooths out bursts of requests and processes them at an approximately average rate.
// Cons:  a burst of traffic can fill up the queue with old requests and starve more recent requests from being processed. It also provides no guarantee that requests get processed in a fixed amount of time
type leakyRateLimit struct {
	tasks         chan Task
	queueSize     int
	requestPerSec int
	interval      time.Duration
	last          time.Time
	mutex         sync.RWMutex
}

func NewLeakyRateLimit(requestPerSec, queueSize int) RateLimit {
	return &leakyRateLimit{
		queueSize:     queueSize,
		tasks:         make(chan Task, queueSize),
		requestPerSec: requestPerSec,
		interval:      time.Second / time.Duration(requestPerSec),
	}
}

func (l *leakyRateLimit) Take() (Task, time.Time) {
	task := <-l.tasks
	l.mutex.Lock()
	defer l.mutex.Unlock()
	now := time.Now()
	if l.last.IsZero() {
		l.last = now
		return task, now
	}
	sleepFor := l.interval - now.Sub(l.last)
	if sleepFor > 0 {
		time.Sleep(sleepFor)
		l.last = now.Add(sleepFor)
	} else {
		l.last = now
	}
	return task, l.last
}

func (l *leakyRateLimit) Receive(task Task) error {
	if len(l.tasks) == l.queueSize {
		return DropErr
	}
	l.tasks <- task
	return nil
}
