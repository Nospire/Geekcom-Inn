package ratelimit

import (
	"testing"
	"time"
)

func TestChatLimiterAllowsBurst(t *testing.T) {
	l := NewChat()
	for i := 0; i < 5; i++ {
		if !l.Allow() {
			t.Errorf("burst message %d rejected", i)
		}
	}
}

func TestChatLimiterRejectsAfterBurst(t *testing.T) {
	l := NewChat()
	for i := 0; i < 5; i++ {
		l.Allow()
	}
	if l.Allow() {
		t.Error("6th message in burst should be rejected")
	}
}

func TestChatLimiterRefills(t *testing.T) {
	l := NewChat()
	for i := 0; i < 5; i++ {
		l.Allow()
	}
	l.mu.Lock()
	l.lastRefill = l.lastRefill.Add(-1 * time.Second)
	l.mu.Unlock()
	if !l.Allow() {
		t.Error("should allow after refill")
	}
}

func TestCanvasLimiterHigherRate(t *testing.T) {
	l := NewCanvas()
	allowed := 0
	for i := 0; i < 35; i++ {
		if l.Allow() {
			allowed++
		}
	}
	if allowed < 30 {
		t.Errorf("canvas allowed %d, want at least 30", allowed)
	}
}

func TestEscalation(t *testing.T) {
	l := NewChat()
	for i := 0; i < 10; i++ {
		l.Allow()
	}

	l.RecordViolation()
	cd := l.CooldownRemaining()
	if cd <= 0 || cd > 6*time.Second {
		t.Errorf("1st violation cooldown = %v, want ~5s", cd)
	}

	l.RecordViolation()
	cd = l.CooldownRemaining()
	if cd <= 5*time.Second || cd > 31*time.Second {
		t.Errorf("2nd violation cooldown = %v, want ~30s", cd)
	}
}

func TestEscalationResets(t *testing.T) {
	l := NewChat()
	l.RecordViolation()
	l.mu.Lock()
	l.lastViolation = l.lastViolation.Add(-6 * time.Minute)
	l.mu.Unlock()
	l.mu.Lock()
	l.maybeResetEscalation()
	v := l.violations
	l.mu.Unlock()
	if v != 0 {
		t.Errorf("violations = %d, want 0 after reset", v)
	}
}
