package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"tavrn/internal/chat"
)

type ChatView struct {
	viewport viewport.Model
	input    textinput.Model
	messages []chat.Message
	width    int
	height   int
}

func NewChatView() ChatView {
	ti := textinput.New()
	ti.Placeholder = "say something..."
	ti.Focus()
	ti.CharLimit = 500
	ti.TextStyle = InputStyle
	ti.PromptStyle = lipgloss.NewStyle().Foreground(ColorHighlight)
	ti.Prompt = "> "

	vp := viewport.New(80, 10)

	return ChatView{
		viewport: vp,
		input:    ti,
		messages: make([]chat.Message, 0),
	}
}

func (c *ChatView) SetSize(width, height int) {
	c.width = width
	c.height = height
	inputHeight := 1
	borderHeight := 2
	c.viewport.Width = width - borderHeight
	c.viewport.Height = height - inputHeight - borderHeight
	if c.viewport.Height < 1 {
		c.viewport.Height = 1
	}
	c.input.Width = width - 4
}

func (c *ChatView) AddMessage(msg chat.Message) {
	c.messages = append(c.messages, msg)
	c.renderMessages()
	c.viewport.GotoBottom()
}

func (c *ChatView) renderMessages() {
	var lines []string
	for _, msg := range c.messages {
		if msg.IsSystem {
			line := SystemMsgStyle.Render("* " + msg.Text)
			lines = append(lines, line)
		} else {
			nick := NickStyle(msg.ColorIndex).Render(msg.Nickname)
			line := fmt.Sprintf("%s: %s", nick, msg.Text)
			lines = append(lines, line)
		}
	}
	c.viewport.SetContent(strings.Join(lines, "\n"))
}

func (c ChatView) Update(msg tea.Msg) (ChatView, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	c.input, cmd = c.input.Update(msg)
	cmds = append(cmds, cmd)

	c.viewport, cmd = c.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return c, tea.Batch(cmds...)
}

func (c ChatView) View() string {
	chatContent := c.viewport.View()
	inputLine := c.input.View()

	inner := lipgloss.JoinVertical(lipgloss.Left, chatContent, inputLine)
	return ChatBorderStyle.Width(c.width).Height(c.height).Render(inner)
}

// InputValue returns current input text and clears it.
func (c *ChatView) InputValue() string {
	val := c.input.Value()
	c.input.Reset()
	return val
}
