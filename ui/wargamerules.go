package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// WargameSignupMsg signals the user wants to join wargames.
type WargameSignupMsg struct{}

type WargameRulesModal struct {
	wargame       string
	isParticipant bool
}

func NewWargameRulesModal(wargame string, isParticipant bool) WargameRulesModal {
	return WargameRulesModal{wargame: wargame, isParticipant: isParticipant}
}

func (w WargameRulesModal) Update(msg tea.Msg) (WargameRulesModal, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		switch keyMsg.String() {
		case "esc", "q":
			return w, func() tea.Msg { return CloseModalMsg{} }
		case "enter":
			return w, func() tea.Msg { return CloseModalMsg{} }
		case "y":
			if !w.isParticipant {
				w.isParticipant = true
				return w, func() tea.Msg { return WargameSignupMsg{} }
			}
		}
	}
	return w, nil
}

func (w WargameRulesModal) View(width, height int) string {
	highlight := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true)
	accent := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true)
	dim := lipgloss.NewStyle().Foreground(ColorDim)
	dimmer := lipgloss.NewStyle().Foreground(ColorDimmer)
	green := lipgloss.NewStyle().Foreground(ColorGreen)
	amber := lipgloss.NewStyle().Foreground(ColorAmber)

	name := strings.ToUpper(w.wargame)

	var b strings.Builder

	b.WriteString(highlight.Render("ВАРГЕЙМ: "+name) + "\n")
	b.WriteString(dimmer.Render(strings.Repeat("─", 38)) + "\n\n")

	b.WriteString(accent.Render(strWargameCatWhat) + "\n")
	b.WriteString(dim.Render(strWargameWhat1) + "\n")
	b.WriteString(dim.Render(strWargameWhat2) + "\n")
	b.WriteString(dim.Render(strWargameWhat3) + "\n")
	b.WriteString(dim.Render(strWargameWhat4) + "\n\n")

	b.WriteString(accent.Render(strWargameCatHow) + "\n")
	b.WriteString(dim.Render(strWargameHow1Pre) + green.Render("Y") + dim.Render(strWargameHow1Suf) + "\n")
	b.WriteString(dim.Render(strWargameHow2Pre) + amber.Render("overthewire.org") + "\n")
	b.WriteString(dim.Render(strWargameHow3) + "\n")
	b.WriteString(dim.Render(strWargameHow4Pre) + green.Render("/submit") + dim.Render(strWargameHow4Suf) + "\n")
	b.WriteString(dim.Render(strWargameHow5) + "\n\n")

	b.WriteString(accent.Render(strWargameCatPoints) + "\n")
	b.WriteString(dim.Render(strWargamePoints1) + "\n")
	b.WriteString(dim.Render(strWargamePoints2) + "\n")
	b.WriteString(dim.Render(strWargamePoints3) + "\n\n")

	// Status + controls
	b.WriteString(dimmer.Render(strings.Repeat("─", 38)) + "\n")
	if w.isParticipant {
		b.WriteString(green.Render(strWargameSignedUp) + dim.Render(strWargameSignedUpSuf) + "\n\n")
		b.WriteString(dimmer.Render("ENTER") + dim.Render(" "+strContinue+"  ") +
			dimmer.Render("ESC") + dim.Render(" "+strClose))
	} else {
		b.WriteString(amber.Render(strWargameNotSigned) + "\n\n")
		b.WriteString(green.Bold(true).Render("Y") + dim.Render(" "+strSignUp+"  ") +
			dimmer.Render("ESC") + dim.Render(" "+strClose))
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(1, 2).
		Render(b.String())
}
