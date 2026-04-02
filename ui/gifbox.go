package ui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"tavrn.sh/internal/chat"
)

const maxAnimatingGifs = 3

// RenderGifBox renders a GIF message as a bordered box with the current frame.
// The frame advances based on timing in TickGifAnimations.
func RenderGifBox(msg *chat.Message) string {
	if len(msg.GifFrames) == 0 {
		return ""
	}

	frame := msg.GifFrames[msg.GifFrame%len(msg.GifFrames)]

	border := lipgloss.NewStyle().Foreground(ColorBorder)
	dim := lipgloss.NewStyle().Foreground(ColorDim)
	accent := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true)

	// Find the widest line in the frame for border sizing
	frameLines := strings.Split(frame, "\n")
	maxWidth := 0
	for _, line := range frameLines {
		w := lipgloss.Width(line)
		if w > maxWidth {
			maxWidth = w
		}
	}

	var b strings.Builder

	// Top border
	title := " GIF "
	headerFill := maxWidth - lipgloss.Width(title) + 2
	if headerFill < 0 {
		headerFill = 0
	}
	b.WriteString(border.Render("╭─") + accent.Render(title) + border.Render(strings.Repeat("─", headerFill)+"╮"))
	b.WriteString("\n")

	// Frame content
	for _, line := range frameLines {
		padding := maxWidth - lipgloss.Width(line)
		if padding < 0 {
			padding = 0
		}
		b.WriteString(border.Render("│ ") + line + strings.Repeat(" ", padding) + border.Render(" │"))
		b.WriteString("\n")
	}

	// Bottom border with title
	label := fmt.Sprintf(" %s ", msg.GifTitle)
	if len(label) > maxWidth {
		label = label[:maxWidth]
	}
	footerFill := maxWidth - lipgloss.Width(label) + 2
	if footerFill < 0 {
		footerFill = 0
	}
	b.WriteString(border.Render("╰"+strings.Repeat("─", footerFill)) + dim.Render(label) + border.Render("╯"))

	return b.String()
}

// TickGifAnimations advances frames for the most recent N animated GIFs.
// Returns true if any frame changed (needs redraw).
func TickGifAnimations(messages []chat.Message) bool {
	now := time.Now()
	changed := false

	// Find the most recent animated GIFs (up to maxAnimatingGifs)
	animating := 0
	for i := len(messages) - 1; i >= 0 && animating < maxAnimatingGifs; i-- {
		msg := &messages[i]
		if !msg.IsGif || len(msg.GifFrames) <= 1 {
			continue
		}
		animating++

		delay := 100 // default
		if msg.GifFrame < len(msg.GifDelays) {
			delay = msg.GifDelays[msg.GifFrame]
		}

		if now.Sub(msg.GifLastTick) >= time.Duration(delay)*time.Millisecond {
			msg.GifFrame = (msg.GifFrame + 1) % len(msg.GifFrames)
			msg.GifLastTick = now
			changed = true
		}
	}

	return changed
}
