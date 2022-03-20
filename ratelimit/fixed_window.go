package ratelimit

import (
	"sync"
	"time"
)

// The system uses a window size of n seconds (typically using human-friendly values, such as 60 or 3600 seconds)
// to track the fixed window algorithm rate. Each incoming request increments the counter for the window. It discards
// the request if the counter exceeds a threshold.

// Pros: ensures more recent requests get processed without being starved by old requests.

// Cons: a single burst of traffic that occurs near the boundary of a window can result in the
// processing of twice the rate of requests because it will allow requests for both the current
// and next windows within a short time. Additionally, if many consumers wait for a reset window,
// they may stampede your API at the same time at the top of the hour.
type fixedWindowRateLimit struct {
	tasks     chan Task
	last      time.Time
	window    time.Duration
	allowance int
	capacity  int
	mutex     sync.Mutex
}

func NewFixedWindowRateLimit(window time.Duration, capacity int) RateLimit {
	return &fixedWindowRateLimit{
		tasks:    make(chan Task, capacity),
		capacity: capacity,
		last:     time.Now(),
		window:   window,
	}
}

func (f *fixedWindowRateLimit) Take() (Task, time.Time) {
	return <-f.tasks, time.Now()
}

func (f *fixedWindowRateLimit) Receive(task Task) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	now := time.Now()
	if now.After(f.last) {
		f.last = now.Add(f.window)
		f.allowance = 0
	}
	f.allowance += 1
	if f.allowance > f.capacity {
		return DropErr
	}
	f.tasks <- task
	return nil
}
