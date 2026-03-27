package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type RoomInfo struct {
	Name  string
	Count int
}

type Sidebar struct {
	Rooms  []RoomInfo
	Width  int
	Height int
}

func NewSidebar() Sidebar {
	return Sidebar{
		Rooms: []RoomInfo{{Name: "lounge", Count: 0}},
	}
}

func (s Sidebar) View() string {
	header := lipgloss.NewStyle().Bold(true).Foreground(ColorSand)

	content := header.Render("Rooms") + "\n"
	for _, r := range s.Rooms {
		line := fmt.Sprintf(" #%-10s %d", r.Name, r.Count)
		content += line + "\n"
	}
	content += "\n"
	content += header.Render("Up Next") + "\n"
	content += lipgloss.NewStyle().Foreground(ColorDim).Render(" (coming soon)") + "\n"

	return SidebarStyle.
		Width(s.Width).
		Height(s.Height).
		MaxHeight(s.Height).
		Render(content)
}
