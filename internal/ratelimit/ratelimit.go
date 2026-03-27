package ratelimit

import (
	"sync"
	"time"
)

type Limiter struct {
	mu            sync.Mutex
	tokens        float64
	maxTokens     float64
	refillRate    float64 // tokens per second
	lastRefill    time.Time
	violations    int
	cooldownUntil time.Time
	lastViolation time.Time
}

// NewChat creates a limiter: 2 msgs/sec sustained, burst of 5.
func NewChat() *Limiter {
	return &Limiter{
		tokens:     5,
		maxTokens:  5,
		refillRate: 2,
		lastRefill: time.Now(),
	}
}

// NewCanvas creates a limiter: 30 ops/sec, burst of 35.
func NewCanvas() *Limiter {
	return &Limiter{
		tokens:     35,
		maxTokens:  35,
		refillRate: 30,
		lastRefill: time.Now(),
	}
}

func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if now.Before(l.cooldownUntil) {
		return false
	}

	elapsed := now.Sub(l.lastRefill).Seconds()
	l.tokens += elapsed * l.refillRate
	if l.tokens > l.maxTokens {
		l.tokens = l.maxTokens
	}
	l.lastRefill = now

	if l.tokens < 1 {
		return false
	}
	l.tokens--
	return true
}

// RecordViolation escalates the cooldown tier.
// Tiers: 5s, 30s, 2min, 10min.
func (l *Limiter) RecordViolation() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.maybeResetEscalation()
	l.violations++
	l.lastViolation = time.Now()

	var cooldown time.Duration
	switch l.violations {
	case 1:
		cooldown = 5 * time.Second
	case 2:
		cooldown = 30 * time.Second
	case 3:
		cooldown = 2 * time.Minute
	default:
		cooldown = 10 * time.Minute
	}
	l.cooldownUntil = time.Now().Add(cooldown)
}

func (l *Limiter) CooldownRemaining() time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()
	remaining := time.Until(l.cooldownUntil)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// maybeResetEscalation resets violations if 5 minutes passed since last violation.
func (l *Limiter) maybeResetEscalation() {
	if !l.lastViolation.IsZero() && time.Since(l.lastViolation) > 5*time.Minute {
		l.violations = 0
	}
}
