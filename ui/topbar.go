package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// Scrolling diagonal frames for the top bar decoration
var topBarDiagFrames = []string{
	"╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱",
	"╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱",
	"╲╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╲",
	"╱╲╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╲╱",
}

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

	// Line 1: Animated diagonal fill bar (full width)
	diagPattern := topBarDiagFrames[t.Frame%len(topBarDiagFrames)]
	// Repeat to fill width
	fillRunes := []rune(diagPattern)
	var fillBuf strings.Builder
	for i := 0; i < t.Width; i++ {
		fillBuf.WriteRune(fillRunes[i%len(fillRunes)])
	}
	pair := artGradientPairs[t.Frame%len(artGradientPairs)]
	diagLine := GradientText(fillBuf.String(), pair[0], pair[1], false)

	// Line 2: Stats — left: online/weekly, center: title, right: room
	onlineDot := lipgloss.NewStyle().Foreground(ColorGreen).Render("●")
	onlineNum := lipgloss.NewStyle().Foreground(ColorSand).Bold(true).Render(
		fmt.Sprintf("%d online", t.OnlineCount))
	weekly := lipgloss.NewStyle().Foreground(ColorDim).Render(
		fmt.Sprintf("%d this week", t.WeeklyCount))
	dot := lipgloss.NewStyle().Foreground(ColorDimmer).Render(" · ")

	statsLeft := fmt.Sprintf("  %s %s%s%s", onlineDot, onlineNum, dot, weekly)

	// Center: big gradient title
	titleText := "TAVRN.SH"
	title := GradientText(titleText, pair[1], pair[0], true)

	// Right: room name
	room := lipgloss.NewStyle().Foreground(ColorAmber).Bold(true).Render(
		fmt.Sprintf("#%s  ", t.Room))

	// Position title in center
	leftW := lipgloss.Width(statsLeft)
	rightW := lipgloss.Width(room)
	titleW := len(titleText)
	centerPos := (t.Width - titleW) / 2
	gapLeft := centerPos - leftW
	gapRight := t.Width - centerPos - titleW - rightW
	if gapLeft < 1 {
		gapLeft = 1
	}
	if gapRight < 1 {
		gapRight = 1
	}

	statsLine := statsLeft + strings.Repeat(" ", gapLeft) + title + strings.Repeat(" ", gapRight) + room

	// Line 3: bottom border
	border := lipgloss.NewStyle().Foreground(ColorBorder).Render(
		"  " + strings.Repeat("─", t.Width-4) + "  ")

	return lipgloss.JoinVertical(lipgloss.Left, diagLine, statsLine, border)
}
