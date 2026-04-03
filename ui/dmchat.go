package ui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"tavrn.sh/internal/dm"
)

type DMSendMsg struct {
	ToFP   string
	ToNick string
	Text   string
}

type DMBackToInboxMsg struct{}

// DMChat renders a conversation with one person.
type DMChat struct {
	peerFP        string
	peerNick      string
	ownFP         string
	ownNick       string
	ownColorIndex int
	viewport      viewport.Model
	input         textinput.Model
	messages      []dm.DirectMessage
	width, height int
}

func NewDMChat(peerFP, peerNick, ownFP, ownNick string, ownColorIndex int) DMChat {
	vp := viewport.New(viewport.WithWidth(80), viewport.WithHeight(10))
	ti := textinput.New()
	ti.Focus()
	ti.Placeholder = fmt.Sprintf("Message %s...", peerNick)
	ti.CharLimit = 500
	return DMChat{
		peerFP:        peerFP,
		peerNick:      peerNick,
		ownFP:         ownFP,
		ownNick:       ownNick,
		ownColorIndex: ownColorIndex,
		viewport:      vp,
		input:         ti,
	}
}

func (d *DMChat) SetSize(w, h int) {
	d.width = w
	d.height = h
	inputH := 1
	headerH := 3
	vpH := h - inputH - headerH
	if vpH < 2 {
		vpH = 2
	}
	d.viewport.SetWidth(w)
	d.viewport.SetHeight(vpH)
	d.input.SetWidth(w - 4)
	d.renderMessages()
}

func (d *DMChat) SetMessages(msgs []dm.DirectMessage) {
	d.messages = msgs
	d.renderMessages()
}

func (d *DMChat) AddMessage(msg dm.DirectMessage) {
	d.messages = append(d.messages, msg)
	d.renderMessages()
}

func (d *DMChat) renderMessages() {
	dim := lipgloss.NewStyle().Foreground(ColorDim)
	sand := lipgloss.NewStyle().Foreground(ColorSand)

	var lines []string
	for _, m := range d.messages {
		ts := m.CreatedAt.Format("15:04")
		timeStr := dim.Render(ts)

		var nickStyle lipgloss.Style
		if m.FromFP == d.ownFP {
			nickStyle = NickStyle(d.ownColorIndex)
		} else {
			nickStyle = lipgloss.NewStyle().Foreground(ColorAccent).Bold(true)
		}

		nick := nickStyle.Render(m.FromNick)
		text := sand.Render(m.Text)
		lines = append(lines, fmt.Sprintf(" %s %s %s", timeStr, nick, text))
	}

	content := strings.Join(lines, "\n")
	d.viewport.SetContent(content)
	d.viewport.GotoBottom()
}

func (d DMChat) Update(msg tea.Msg) (DMChat, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		switch keyMsg.String() {
		case "esc":
			return d, func() tea.Msg { return DMBackToInboxMsg{} }
		case "enter":
			text := strings.TrimSpace(d.input.Value())
			if text != "" {
				d.input.Reset()
				return d, func() tea.Msg {
					return DMSendMsg{
						ToFP:   d.peerFP,
						ToNick: d.peerNick,
						Text:   text,
					}
				}
			}
			return d, nil
		case "shift+up", "pgup":
			d.viewport.ScrollUp(3)
			return d, nil
		case "shift+down", "pgdown":
			d.viewport.ScrollDown(3)
			return d, nil
		}
	}
	var cmd tea.Cmd
	d.input, cmd = d.input.Update(msg)
	return d, cmd
}

func (d DMChat) View() string {
	accent := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true)
	dimmer := lipgloss.NewStyle().Foreground(ColorDimmer)
	dim := lipgloss.NewStyle().Foreground(ColorDim)

	contentW := d.width - 4
	if contentW < 20 {
		contentW = 20
	}

	// Header
	var header strings.Builder
	header.WriteString("  " + accent.Render("DM: "+d.peerNick))
	header.WriteString("\n")
	header.WriteString("  " + dimmer.Render(strings.Repeat("─", contentW)))
	header.WriteString("\n")

	// Footer with input
	prompt := dim.Render(" > ")
	inputLine := prompt + d.input.View()

	return header.String() + d.viewport.View() + "\n" + inputLine
}
