package server

import "time"

type RateLimiter struct {
	tokens     int
	maxTokens  int
	refillRate time.Duration
	lastRefill time.Time
}

func NewRateLimiter(max int, rate time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:     max,
		maxTokens:  max,
		refillRate: rate,
		lastRefill: time.Now(),
	}
}

func (rl *RateLimiter) Allow() bool {
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	tokensToAdd := int(elapsed / rl.refillRate)

	if tokensToAdd > 0 {
		rl.tokens += tokensToAdd
		if rl.tokens > rl.maxTokens {
			rl.tokens = rl.maxTokens
		}
		rl.lastRefill = rl.lastRefill.Add(time.Duration(tokensToAdd) * rl.refillRate)
	}

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}
