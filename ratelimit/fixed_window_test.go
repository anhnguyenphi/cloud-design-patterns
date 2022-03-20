package ratelimit

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestFixedWindowRateLimit_Receive(t *testing.T) {
	rl := NewFixedWindowRateLimit(time.Second*5, 5)
	start := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		count := 0
		for count < 10 {
			task, now := rl.Take()
			fmt.Println("Complete Task:", task.ID, now.Sub(start))
			count++
		}
		wg.Done()
	}()

	for i := 0; i < 20; i++ {
		err := rl.Receive(Task{ID: i})
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Millisecond * 500)
	}
	wg.Wait()
	fmt.Println("Total: ", time.Now().Sub(start))
}
