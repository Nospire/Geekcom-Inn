package ui

type BottomBar struct {
	Width int
}

func NewBottomBar() BottomBar {
	return BottomBar{}
}

func (b BottomBar) View() string {
	return BottomBarStyle.Width(b.Width).MaxWidth(b.Width).Render(" /help for commands")
}
