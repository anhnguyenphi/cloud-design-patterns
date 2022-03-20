package ratelimit

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestLeakyRateLimit_Take(t *testing.T) {
	rl := NewLeakyRateLimit(100, 100)
	start := time.Now()
	prev := start
	wg := sync.WaitGroup{}

	for i := 0; i < 120; i++ {
		wg.Add(1)
		go func(y int) {
			err := rl.Receive(Task{ID: y})
			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}(i)
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(y int) {
			task, now := rl.Take()
			fmt.Println(task.ID, now.Sub(prev))
			prev = now
			wg.Done()
		}(i)
	}
	wg.Wait()

	fmt.Println("Total: ", time.Now().Sub(start))
}
