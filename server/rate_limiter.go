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

	player_notif := strings.HasPrefix(cmd, pr.CmdNotifyPlayerPosition)
	player_position := strings.HasPrefix(cmd, pr.CmdGetPlayerPositions)
	item_notif := strings.HasPrefix(cmd, pr.CmdNotifyItemPosition)
	item_position := strings.HasPrefix(cmd, pr.CmdGetItemPositions)

	if player_notif || player_position || item_notif || item_position {
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
