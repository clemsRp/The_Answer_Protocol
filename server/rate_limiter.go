package server

import (
	"strings"
	pr "tap/protocol"
	"time"
)

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

func (rl *RateLimiter) Allow(cmd string) bool {
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	tokensToAdd := int(elapsed / rl.refillRate)
	if strings.HasPrefix(cmd, pr.CmdNotifyPosition) {
		return true
	}
	if strings.HasPrefix(cmd, pr.CmdGetPositions) {
		return true
	}
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
