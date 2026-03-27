package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

type TopBar struct {
	Room        string
	OnlineCount int
	WeeklyCount int
	NowPlaying  string
	Width       int
	Frame       int
}

func NewTopBar() TopBar {
	return TopBar{Room: "lounge"}
}

func (t TopBar) View() string {
	if t.Width < 20 {
		return ""
	}

	// Line 1: ╱╱╱ decorative fill with title centered
	titleText := " TAVRN.SH "
	titleRendered := GradientText(titleText, ColorHighlight, ColorAmber, true)
	titleWidth := len(titleText)

	fillTotal := t.Width - titleWidth - 4
	if fillTotal < 4 {
		fillTotal = 4
	}
	leftFillN := fillTotal / 2
	rightFillN := fillTotal - leftFillN

	leftFill := lipgloss.NewStyle().Foreground(ColorBorder).Render(
		"  " + strings.Repeat("╱", leftFillN))
	rightFill := lipgloss.NewStyle().Foreground(ColorBorder).Render(
		strings.Repeat("╱", rightFillN) + "  ")

	brandLine := leftFill + titleRendered + rightFill

	// Line 2: Stats
	onlineDot := lipgloss.NewStyle().Foreground(ColorGreen).Render("●")
	onlineNum := lipgloss.NewStyle().Foreground(ColorSand).Bold(true).Render(
		fmt.Sprintf("%d online", t.OnlineCount))
	weekly := lipgloss.NewStyle().Foreground(ColorDim).Render(
		fmt.Sprintf("%d this week", t.WeeklyCount))
	room := lipgloss.NewStyle().Foreground(ColorAmber).Bold(true).Render(
		fmt.Sprintf("#%s", t.Room))
	dot := lipgloss.NewStyle().Foreground(ColorDimmer).Render(" · ")

	statsLeft := fmt.Sprintf("  %s %s%s%s", onlineDot, onlineNum, dot, weekly)

	// Right-align room name
	statsRight := room + "  "
	gap := t.Width - lipgloss.Width(statsLeft) - lipgloss.Width(statsRight)
	if gap < 0 {
		gap = 0
	}

	statsLine := statsLeft + strings.Repeat(" ", gap) + statsRight

	// Line 3: Bottom border
	border := lipgloss.NewStyle().Foreground(ColorBorder).Render(
		"  " + strings.Repeat("─", t.Width-4) + "  ")

	return lipgloss.JoinVertical(lipgloss.Left, brandLine, statsLine, border)
}
