package ratelimit

import "time"

type (
	RateLimit interface {
		Take() time.Time
	}
)
