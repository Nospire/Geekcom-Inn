package ui

import (
	"fmt"
	"testing"
	"time"

	"charm.land/lipgloss/v2"
)

func TestFormatTimestamp_JustNow(t *testing.T) {
	now := time.Now()
	got := formatTimestamp(now, now)
	if got != strTsJustNow {
		t.Errorf("got %q, want %q", got, strTsJustNow)
	}
}

func TestFormatTimestamp_FewSecondsAgo(t *testing.T) {
	now := time.Now()
	got := formatTimestamp(now.Add(-5*time.Second), now)
	if got != strTsJustNow {
		t.Errorf("got %q, want %q (under 10s)", got, strTsJustNow)
	}
}

func TestFormatTimestamp_SecondsAgo(t *testing.T) {
	now := time.Now()
	got := formatTimestamp(now.Add(-30*time.Second), now)
	want := fmt.Sprintf(strTsSecsAgo, 30)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFormatTimestamp_OneMinAgo(t *testing.T) {
	now := time.Now()
	got := formatTimestamp(now.Add(-90*time.Second), now)
	if got != strTs1MinAgo {
		t.Errorf("got %q, want %q", got, strTs1MinAgo)
	}
}

func TestFormatTimestamp_MinutesAgo(t *testing.T) {
	now := time.Now()
	got := formatTimestamp(now.Add(-10*time.Minute), now)
	want := fmt.Sprintf(strTsMinsAgo, 10)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFormatTimestamp_HoursAgo(t *testing.T) {
	now := time.Now()
	ts := now.Add(-3 * time.Hour)
	got := formatTimestamp(ts, now)
	expected := ts.Format("15:04")
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestFormatTimestamp_DaysAgo(t *testing.T) {
	now := time.Now()
	ts := now.Add(-48 * time.Hour)
	got := formatTimestamp(ts, now)
	expected := ts.Format("Jan 02 15:04")
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestFormatTimestamp_ExactBoundary10s(t *testing.T) {
	now := time.Now()
	got := formatTimestamp(now.Add(-10*time.Second), now)
	want := fmt.Sprintf(strTsSecsAgo, 10)
	if got != want {
		t.Errorf("at exactly 10s boundary got %q, want %q", got, want)
	}
}

func TestFormatTimestamp_ExactBoundary1Min(t *testing.T) {
	now := time.Now()
	got := formatTimestamp(now.Add(-60*time.Second), now)
	if got != strTs1MinAgo {
		t.Errorf("at exactly 60s boundary got %q, want %q", got, strTs1MinAgo)
	}
}

func TestWordWrap_Short(t *testing.T) {
	lines := wordWrap("hello", 80)
	if len(lines) != 1 || lines[0] != "hello" {
		t.Errorf("got %v", lines)
	}
}

func TestWordWrap_Long(t *testing.T) {
	text := "the quick brown fox jumps over the lazy dog"
	lines := wordWrap(text, 20)
	if len(lines) < 2 {
		t.Errorf("expected wrapping, got %d lines", len(lines))
	}
	for _, line := range lines {
		if lipgloss.Width(line) > 20 {
			t.Errorf("line %q display width %d exceeds 20", line, lipgloss.Width(line))
		}
	}
}

func TestWordWrap_Emoji(t *testing.T) {
	// Emojis are 2 cells wide — "🍺" takes 2 columns
	// With width 10: "hello 🍺" = 5+1+2 = 8 cols, "world" = 5 cols
	// Total "hello 🍺 world" = 14 cols, should wrap
	text := "hello 🍺 world"
	lines := wordWrap(text, 10)
	if len(lines) < 2 {
		t.Errorf("expected wrapping with emoji, got %d lines: %v", len(lines), lines)
	}
	for _, line := range lines {
		if lipgloss.Width(line) > 10 {
			t.Errorf("line %q display width %d exceeds 10", line, lipgloss.Width(line))
		}
	}
}

func TestWordWrap_ZeroWidth(t *testing.T) {
	lines := wordWrap("hello world", 0)
	if len(lines) != 1 {
		t.Errorf("zero width should return single line, got %d", len(lines))
	}
}

func TestWordWrap_ExactFit(t *testing.T) {
	lines := wordWrap("hello", 5)
	if len(lines) != 1 || lines[0] != "hello" {
		t.Errorf("got %v", lines)
	}
}

func TestWordWrap_SingleLongWord(t *testing.T) {
	// A single word longer than width can't be split
	lines := wordWrap("superlongword", 5)
	if len(lines) != 1 || lines[0] != "superlongword" {
		t.Errorf("got %v", lines)
	}
}
