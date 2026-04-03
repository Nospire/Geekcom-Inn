package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"tavrn.sh/internal/dm"
)

type DMOpenConvoMsg struct {
	PeerFP   string
	PeerNick string
}

// DMInbox shows the list of DM conversations.
type DMInbox struct {
	conversations []dm.Conversation
	cursor        int
	width, height int
}

func NewDMInbox(convos []dm.Conversation) DMInbox {
	return DMInbox{
		conversations: convos,
	}
}

func (d *DMInbox) SetSize(w, h int) {
	d.width = w
	d.height = h
}

func (d *DMInbox) SetConversations(convos []dm.Conversation) {
	d.conversations = convos
	if d.cursor >= len(convos) {
		d.cursor = len(convos) - 1
	}
	if d.cursor < 0 {
		d.cursor = 0
	}
}

func (d DMInbox) Update(msg tea.Msg) (DMInbox, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		switch keyMsg.String() {
		case "up", "k":
			if d.cursor > 0 {
				d.cursor--
			}
		case "down", "j":
			if d.cursor < len(d.conversations)-1 {
				d.cursor++
			}
		case "enter":
			if len(d.conversations) > 0 {
				c := d.conversations[d.cursor]
				return d, func() tea.Msg {
					return DMOpenConvoMsg{PeerFP: c.PeerFP, PeerNick: c.PeerNick}
				}
			}
		}
	}
	return d, nil
}

func (d DMInbox) View() string {
	accent := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true)
	dimmer := lipgloss.NewStyle().Foreground(ColorDimmer)
	dim := lipgloss.NewStyle().Foreground(ColorDim)
	sand := lipgloss.NewStyle().Foreground(ColorSand)
	highlight := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true)
	amber := lipgloss.NewStyle().Foreground(ColorAmber).Bold(true)

	var b strings.Builder

	contentW := d.width - 4
	if contentW < 20 {
		contentW = 20
	}

	b.WriteString(accent.Render("  DIRECT MESSAGES"))
	b.WriteString("\n")
	b.WriteString("  " + dimmer.Render(strings.Repeat("─", contentW)))
	b.WriteString("\n\n")

	if len(d.conversations) == 0 {
		b.WriteString("  " + dim.Render("No conversations yet."))
		b.WriteString("\n\n")
		b.WriteString("  " + dim.Render("Use /dm <name> in the tavern to start one."))
		b.WriteString("\n")
	}

	maxNameW := 16
	maxPreviewW := contentW - maxNameW - 10
	if maxPreviewW < 10 {
		maxPreviewW = 10
	}

	for i, c := range d.conversations {
		isCurrent := i == d.cursor

		name := c.PeerNick
		if len(name) > maxNameW {
			name = name[:maxNameW-1] + "."
		}

		preview := c.LastMessage
		if len(preview) > maxPreviewW {
			preview = preview[:maxPreviewW-1] + "."
		}
		// Strip newlines from preview
		preview = strings.ReplaceAll(preview, "\n", " ")

		var line string
		if isCurrent {
			indicator := highlight.Render("▸ ")
			nameStr := amber.Render(name)
			previewStr := sand.Render(preview)
			line = indicator + nameStr + "  " + previewStr
		} else {
			nameStr := sand.Render(name)
			previewStr := dim.Render(preview)
			line = "  " + nameStr + "  " + previewStr
		}

		if c.UnreadCount > 0 {
			badge := lipgloss.NewStyle().Foreground(ColorAmber).Bold(true).
				Render(fmt.Sprintf(" (%d)", c.UnreadCount))
			line += badge
		}

		b.WriteString(line)
		b.WriteString("\n")
	}

	// Footer hints
	b.WriteString("\n")
	b.WriteString("  " + dimmer.Render("↑↓ navigate  ENTER open  TAB back to tavern"))

	return b.String()
}
