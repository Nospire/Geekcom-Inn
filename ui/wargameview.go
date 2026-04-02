package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// WargameHeader renders the progress header for a wargame room.
// Sits above the chat in wargame rooms.
func WargameHeader(wargame string, currentLevel, maxLevel, points, width int) string {
	accent := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true)
	highlight := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true)
	dim := lipgloss.NewStyle().Foreground(ColorDim)
	dimmer := lipgloss.NewStyle().Foreground(ColorDimmer)
	green := lipgloss.NewStyle().Foreground(ColorGreen)
	amber := lipgloss.NewStyle().Foreground(ColorAmber)

	var b strings.Builder

	// Title line
	name := strings.ToUpper(wargame)
	b.WriteString(accent.Render("  "+name) + "  ")

	// Level + points
	b.WriteString(highlight.Render(fmt.Sprintf("Lv.%d", currentLevel)))
	if maxLevel > 0 {
		b.WriteString(dimmer.Render(fmt.Sprintf("/%d", maxLevel)))
	}
	b.WriteString(dim.Render("  "))
	b.WriteString(amber.Render(fmt.Sprintf("%d pts", points)))

	// Right-align /submit hint
	hint := dim.Render("/submit")
	leftW := lipgloss.Width(name) + 2 + // "  NAME  "
		lipgloss.Width(fmt.Sprintf("Lv.%d", currentLevel)) +
		lipgloss.Width(fmt.Sprintf("/%d", maxLevel)) +
		2 + // gap
		lipgloss.Width(fmt.Sprintf("%d pts", points))
	hintW := lipgloss.Width(hint)
	gap := width - leftW - hintW - 4
	if gap < 1 {
		gap = 1
	}
	b.WriteString(strings.Repeat(" ", gap))
	b.WriteString(hint)
	b.WriteString("\n")

	// Progress bar
	barW := width - 6
	if barW < 10 {
		barW = 10
	}
	if barW > 50 {
		barW = 50
	}
	filled := 0
	if maxLevel > 0 {
		filled = currentLevel * barW / maxLevel
	}
	if filled > barW {
		filled = barW
	}

	bar := green.Render(strings.Repeat("█", filled)) +
		dimmer.Render(strings.Repeat("░", barW-filled))
	b.WriteString("  " + bar)
	b.WriteString("\n")

	// Separator
	b.WriteString("  " + dimmer.Render(strings.Repeat("─", width-6)))
	b.WriteString("\n")

	return b.String()
}
