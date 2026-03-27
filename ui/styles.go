package ui

import "github.com/charmbracelet/lipgloss"

// Cantina palette — ANSI 256 colors for max terminal compatibility.
var (
	ColorBackground = lipgloss.Color("235")
	ColorSand       = lipgloss.Color("180")
	ColorDim        = lipgloss.Color("243")
	ColorBorder     = lipgloss.Color("94")
	ColorHighlight  = lipgloss.Color("179")
	ColorAmber      = lipgloss.Color("172")

	NickColors = []lipgloss.Color{
		lipgloss.Color("174"), // dusty rose
		lipgloss.Color("109"), // faded teal
		lipgloss.Color("137"), // aged copper
		lipgloss.Color("138"), // soft clay
		lipgloss.Color("108"), // pale sage
		lipgloss.Color("179"), // weathered gold
		lipgloss.Color("140"), // dim lavender
		lipgloss.Color("67"),  // smoky blue
		lipgloss.Color("131"), // muted coral
		lipgloss.Color("144"), // warm stone
		lipgloss.Color("136"), // quiet amber
		lipgloss.Color("97"),  // dusk violet
	}
)

var (
	TopBarStyle = lipgloss.NewStyle().
			Foreground(ColorSand).
			Background(lipgloss.Color("236")).
			Bold(true).
			Padding(0, 1)

	BottomBarStyle = lipgloss.NewStyle().
			Foreground(ColorDim).
			Background(lipgloss.Color("236")).
			Padding(0, 1)

	ChatBorderStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorBorder).
			Foreground(ColorSand)

	SidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(ColorBorder).
			Foreground(ColorSand).
			Padding(0, 1)

	CanvasPlaceholderStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(ColorBorder).
				Foreground(ColorDim).
				Align(lipgloss.Center, lipgloss.Center)

	SystemMsgStyle = lipgloss.NewStyle().
			Foreground(ColorDim).
			Italic(true)

	InputStyle = lipgloss.NewStyle().
			Foreground(ColorSand)
)

func NickStyle(colorIndex int) lipgloss.Style {
	idx := colorIndex % len(NickColors)
	return lipgloss.NewStyle().Foreground(NickColors[idx]).Bold(true)
}
