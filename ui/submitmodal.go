package ui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// SubmitFlagMsg carries the submitted flag back to the app.
type SubmitFlagMsg struct {
	Flag string
}

type SubmitModal struct {
	input    textinput.Model
	wargame  string
	level    int // next level to clear
	maxLevel int
	err      string
}

func NewSubmitModal(wargame string, currentLevel, maxLevel int) SubmitModal {
	ti := textinput.New()
	ti.Placeholder = strSubmitPlaceholder
	ti.Focus()
	ti.CharLimit = 200
	ti.Prompt = "> "
	ti.EchoMode = textinput.EchoPassword

	return SubmitModal{
		input:    ti,
		wargame:  wargame,
		level:    currentLevel + 1,
		maxLevel: maxLevel,
	}
}

func (s SubmitModal) Update(msg tea.Msg) (SubmitModal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			return s, func() tea.Msg { return CloseModalMsg{} }
		case "enter":
			flag := strings.TrimSpace(s.input.Value())
			if flag == "" {
				s.err = strSubmitFlagEmpty
				return s, nil
			}
			return s, func() tea.Msg { return SubmitFlagMsg{Flag: flag} }
		}
	}
	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)
	return s, cmd
}

func (s SubmitModal) View(width, height int) string {
	accent := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true)
	highlight := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true)
	dim := lipgloss.NewStyle().Foreground(ColorDim)
	dimmer := lipgloss.NewStyle().Foreground(ColorDimmer)
	errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	var b strings.Builder

	// Header
	b.WriteString(highlight.Render(strSubmitTitle))
	b.WriteString("\n\n")

	// Wargame + level info
	b.WriteString(accent.Render(strings.ToUpper(s.wargame)))
	b.WriteString(dim.Render(fmt.Sprintf(strSubmitLevelFmt, s.level)))
	if s.maxLevel > 0 {
		b.WriteString(dimmer.Render(fmt.Sprintf("/%d", s.maxLevel)))
	}
	b.WriteString("\n")
	b.WriteString(dimmer.Render(strings.Repeat("─", 30)))
	b.WriteString("\n\n")

	// Input
	b.WriteString(dim.Render(strSubmitFlagLabel))
	b.WriteString("\n")
	b.WriteString(s.input.View())
	b.WriteString("\n\n")

	if s.err != "" {
		b.WriteString(errStyle.Render(s.err))
		b.WriteString("\n\n")
	}

	// Controls
	b.WriteString(dimmer.Render("ENTER") + dim.Render(" "+strSubmitBtn+"  ") +
		dimmer.Render("ESC") + dim.Render(" "+strClose))

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(1, 2).
		Render(b.String())
}
