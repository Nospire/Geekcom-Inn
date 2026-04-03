package ui

import (
	"fmt"

	"charm.land/lipgloss/v2"
)

type BottomBar struct {
	Width        int
	IsGallery    bool
	IsTankard    bool
	IsDMMode     bool
	MentionCount int
	DMUnread     int
}

func NewBottomBar() BottomBar {
	return BottomBar{}
}

func (b BottomBar) View() string {
	k := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true)
	d := lipgloss.NewStyle().Foreground(ColorDim)
	sep := lipgloss.NewStyle().Foreground(ColorDimmer).Render("  ·  ")

	// DM badge
	dmBadge := ""
	if b.DMUnread > 0 {
		dmBadge = fmt.Sprintf("(%d)", b.DMUnread)
	}

	var content string
	if b.IsDMMode {
		content = "  " +
			k.Render("TAB") + " " + d.Render(strBBarTavern) + sep +
			k.Render("ESC") + " " + d.Render(strBBarBack) + sep +
			k.Render("↑↓") + " " + d.Render(strBBarNavigate) + sep +
			k.Render("ENTER") + " " + d.Render(strBBarOpen)
	} else if b.IsTankard {
		content = "  " +
			k.Render("SPACE") + " " + d.Render(strBBarDrink) + sep +
			k.Render("ESC") + " " + d.Render(strBBarBack)
	} else if b.IsGallery {
		content = "  " +
			k.Render("P") + " " + d.Render(strBBarPost) + sep +
			k.Render("E") + " " + d.Render(strBBarExpand) + sep +
			k.Render("D") + " " + d.Render(strBBarDelete) + sep +
			k.Render("TAB") + " " + d.Render(strBBarSelect)
	} else {
		f4 := k.Render("F4") + " " + d.Render(strBBarMentions)
		if b.MentionCount > 0 {
			f4 = k.Render("F4") + " " + d.Render(fmt.Sprintf(strBBarMentionsFmt, b.MentionCount))
		}
		tabDM := k.Render("TAB") + " " + d.Render(strBBarDMs+dmBadge)
		content = "  " +
			k.Render("F1") + " " + d.Render(strBBarHelp) + sep +
			k.Render("F2") + " " + d.Render(strBBarNick) + sep +
			k.Render("F3") + " " + d.Render(strBBarRooms) + sep +
			f4 + sep +
			tabDM + sep +
			k.Render("F6") + " " + d.Render(strBBarTankard) + sep +
			k.Render("F7") + " " + d.Render(strBBarLeaderboard) + sep +
			k.Render("SHIFT+↑↓") + " " + d.Render(strBBarScroll)
	}

	return BottomBarStyle.Width(b.Width).MaxWidth(b.Width).Render(content)
}
