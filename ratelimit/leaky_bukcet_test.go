package ratelimit

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestLeakyRateLimit_Take(t *testing.T) {
	rl := NewLeakyRateLimit(100)
	start := time.Now()
	prev := start
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(y int) {
			now := rl.Take()
			fmt.Println(y, now.Sub(prev))
			prev = now
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("Total: ", time.Now().Sub(start))
}
