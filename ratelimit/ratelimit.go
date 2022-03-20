package ratelimit

import (
	"errors"
	"time"
)

type (
	RateLimit interface {
		Take() (Task, time.Time)
		Receive(task Task) error
	}

	Task struct {
		ID int
	}
)

var (
	DropErr = errors.New("drop message")
)
