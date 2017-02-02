package bot

import (
	"sync"
	"time"

	rate "github.com/beefsack/go-rate"
)

type RateLimiter struct {
	count    int
	duration time.Duration

	limiters map[string]*rate.RateLimiter
	mu       sync.Mutex
}

func NewRateLimiter(count int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		count:    count,
		duration: duration,
		limiters: map[string]*rate.RateLimiter{},
	}
}

func (rl *RateLimiter) Limited(name string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, ok := rl.limiters[name]
	if !ok {
		limiter = rate.New(rl.count, rl.duration)
		rl.limiters[name] = limiter
	}

	ok, _ = limiter.Try()
	return !ok
}
