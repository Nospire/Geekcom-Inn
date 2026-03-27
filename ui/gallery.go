package ui

import (
	"math/rand"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"tavrn/internal/store"
)

type GalleryNote struct {
	ID          int
	X, Y        int
	Text        string
	Nickname    string
	Fingerprint string
	ColorIndex  int
}

type GalleryView struct {
	notes       []GalleryNote
	selected    int
	dragging    bool
	dragOffsetX int
	dragOffsetY int
	width       int
	height      int
	fingerprint string
	rng         *rand.Rand
}

func NewGalleryView(fingerprint string) GalleryView {
	return GalleryView{
		notes:       make([]GalleryNote, 0),
		selected:    -1,
		fingerprint: fingerprint,
		rng:         rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (g *GalleryView) SetSize(w, h int)              { g.width = w; g.height = h }
func (g *GalleryView) AddNote(n GalleryNote)          { g.notes = append(g.notes, n) }
func (g *GalleryView) ClearAll()                      { g.notes = nil; g.selected = -1 }

func (g *GalleryView) LoadNotes(rows []store.NoteRow) {
	g.notes = make([]GalleryNote, len(rows))
	for i, r := range rows {
		g.notes[i] = GalleryNote{
			ID: r.ID, X: r.X, Y: r.Y, Text: r.Text,
			Nickname: r.Nickname, Fingerprint: r.Fingerprint, ColorIndex: r.ColorIndex,
		}
	}
}

func (g *GalleryView) MoveNote(id, x, y int) {
	for i := range g.notes {
		if g.notes[i].ID == id {
			g.notes[i].X = x
			g.notes[i].Y = y
			return
		}
	}
}

func (g *GalleryView) RemoveNote(id int) {
	for i := range g.notes {
		if g.notes[i].ID == id {
			g.notes = append(g.notes[:i], g.notes[i+1:]...)
			if g.selected >= len(g.notes) {
				g.selected = -1
			}
			return
		}
	}
}

func (g *GalleryView) RandomPosition() (int, int) {
	maxX := g.width - 30
	maxY := g.height - 8
	if maxX < 4 {
		maxX = 4
	}
	if maxY < 2 {
		maxY = 2
	}
	return 2 + g.rng.Intn(maxX), 1 + g.rng.Intn(maxY)
}

func (g GalleryView) Update(msg tea.Msg) (GalleryView, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "d", "delete", "backspace":
			if g.selected >= 0 && g.selected < len(g.notes) {
				note := g.notes[g.selected]
				if note.Fingerprint == g.fingerprint {
					return g, func() tea.Msg {
						return GalleryDeleteMsg{NoteID: note.ID}
					}
				}
			}
		case "tab":
			if len(g.notes) > 0 {
				g.selected = (g.selected + 1) % len(g.notes)
			}
		}

	case tea.MouseClickMsg:
		mx, my := msg.Mouse().X, msg.Mouse().Y
		hit := g.hitTest(mx, my)
		if hit >= 0 {
			g.selected = hit
			g.dragging = true
			g.dragOffsetX = mx - g.notes[hit].X
			g.dragOffsetY = my - g.notes[hit].Y
		} else {
			g.selected = -1
		}

	case tea.MouseReleaseMsg:
		if g.dragging && g.selected >= 0 && g.selected < len(g.notes) {
			note := g.notes[g.selected]
			if note.Fingerprint == g.fingerprint {
				g.dragging = false
				return g, func() tea.Msg {
					return GalleryMoveMsg{NoteID: note.ID, X: note.X, Y: note.Y}
				}
			}
		}
		g.dragging = false

	case tea.MouseMotionMsg:
		if g.dragging && g.selected >= 0 && g.selected < len(g.notes) {
			mx, my := msg.Mouse().X, msg.Mouse().Y
			newX := mx - g.dragOffsetX
			newY := my - g.dragOffsetY
			if newX < 0 {
				newX = 0
			}
			if newY < 0 {
				newY = 0
			}
			if newX > g.width-5 {
				newX = g.width - 5
			}
			if newY > g.height-3 {
				newY = g.height - 3
			}
			g.notes[g.selected].X = newX
			g.notes[g.selected].Y = newY
		}
	}

	return g, nil
}

func (g GalleryView) hitTest(mx, my int) int {
	// Reverse order: topmost (last rendered) first
	for i := len(g.notes) - 1; i >= 0; i-- {
		n := g.notes[i]
		w, h := noteSize(n.Text)
		if mx >= n.X && mx < n.X+w && my >= n.Y && my < n.Y+h {
			return i
		}
	}
	return -1
}

func (g GalleryView) View() string {
	if g.width == 0 || g.height == 0 {
		return ""
	}

	// Build background grid with dashed pattern
	grid := make([][]rune, g.height)
	gridColor := make([][]int, g.height) // -1=bg, -2=bgPattern, 0-11=note color
	for y := range grid {
		grid[y] = make([]rune, g.width)
		gridColor[y] = make([]int, g.width)
		for x := range grid[y] {
			// Diagonal hatch pattern like the bubbletea example
			if (x+y)%4 == 0 {
				grid[y][x] = '╲'
				gridColor[y][x] = -2 // bg pattern
			} else if (x+y)%4 == 2 {
				grid[y][x] = '╱'
				gridColor[y][x] = -2
			} else {
				grid[y][x] = ' '
				gridColor[y][x] = -1
			}
		}
	}

	// Empty state hint
	if len(g.notes) == 0 {
		hint := "CTRL+P to post a note"
		cx := (g.width - len(hint)) / 2
		cy := g.height / 2
		if cx < 0 {
			cx = 0
		}
		for i, ch := range hint {
			rx := cx + i
			if rx >= 0 && rx < g.width && cy >= 0 && cy < g.height {
				grid[cy][rx] = ch
				gridColor[cy][rx] = -3 // hint color
			}
		}
	}

	// Stamp notes onto grid
	for idx, note := range g.notes {
		isSelected := idx == g.selected
		isOwn := note.Fingerprint == g.fingerprint
		card := renderNoteCard(note, isSelected, isOwn)
		for dy, line := range card {
			ry := note.Y + dy
			if ry < 0 || ry >= g.height {
				continue
			}
			rx := note.X
			for _, ch := range line {
				if rx >= 0 && rx < g.width {
					grid[ry][rx] = ch
					gridColor[ry][rx] = note.ColorIndex
				}
				rx++
			}
		}
	}

	// Render to styled string
	bgPatternStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("237"))
	hintStyle := lipgloss.NewStyle().Foreground(ColorDim).Italic(true)

	var result []string
	for y := 0; y < g.height; y++ {
		var b strings.Builder
		prevKind := -99
		run := ""

		flushRun := func() {
			if run == "" {
				return
			}
			switch prevKind {
			case -1:
				b.WriteString(run) // plain space
			case -2:
				b.WriteString(bgPatternStyle.Render(run))
			case -3:
				b.WriteString(hintStyle.Render(run))
			default:
				b.WriteString(lipgloss.NewStyle().Foreground(NickColors[prevKind%len(NickColors)]).Render(run))
			}
			run = ""
		}

		for x := 0; x < g.width; x++ {
			kind := gridColor[y][x]
			ch := string(grid[y][x])
			if kind != prevKind {
				flushRun()
				prevKind = kind
			}
			run += ch
		}
		flushRun()
		result = append(result, b.String())
	}
	return strings.Join(result, "\n")
}

func renderNoteCard(note GalleryNote, selected, isOwn bool) []string {
	// Bigger cards with padding
	maxTextW := 24
	text := note.Text
	if len(text) > maxTextW {
		text = text[:maxTextW]
	}

	innerW := maxTextW + 2 // padding inside borders
	if innerW < 14 {
		innerW = 14
	}

	// Border chars
	tl, tr, bl, br := "╭", "╮", "╰", "╯"
	h, v := "─", "│"
	if selected {
		tl, tr, bl, br = "┏", "┓", "┗", "┛"
		h, v = "━", "┃"
	}

	var lines []string

	// Top border
	lines = append(lines, tl+strings.Repeat(h, innerW)+tr)

	// Empty padding line
	lines = append(lines, v+strings.Repeat(" ", innerW)+v)

	// Text lines (word wrapped)
	wrapped := wrapText(text, innerW-2)
	for _, wl := range wrapped {
		pad := innerW - 2 - len(wl)
		if pad < 0 {
			pad = 0
		}
		lines = append(lines, v+" "+wl+strings.Repeat(" ", pad)+" "+v)
	}

	// Empty padding line
	lines = append(lines, v+strings.Repeat(" ", innerW)+v)

	// Separator
	lines = append(lines, v+strings.Repeat("·", innerW)+v)

	// Nick line
	nick := note.Nickname
	if isOwn {
		nick = "~" + nick
	}
	nickPad := innerW - 2 - len(nick)
	if nickPad < 0 {
		nickPad = 0
	}
	lines = append(lines, v+" "+nick+strings.Repeat(" ", nickPad)+" "+v)

	// Bottom border
	lines = append(lines, bl+strings.Repeat(h, innerW)+br)

	return lines
}

func wrapText(text string, width int) []string {
	if width <= 0 || len(text) <= width {
		return []string{text}
	}
	var lines []string
	words := strings.Fields(text)
	cur := ""
	for _, w := range words {
		if cur == "" {
			cur = w
		} else if len(cur)+1+len(w) <= width {
			cur += " " + w
		} else {
			lines = append(lines, cur)
			cur = w
		}
	}
	if cur != "" {
		lines = append(lines, cur)
	}
	return lines
}

func noteSize(text string) (int, int) {
	maxTextW := 24
	innerW := maxTextW + 2
	if innerW < 14 {
		innerW = 14
	}
	w := innerW + 2 // + border chars
	wrapped := wrapText(text, innerW-2)
	h := len(wrapped) + 5 // top + pad + text + pad + separator + nick + bottom
	return w, h
}

type GalleryDeleteMsg struct{ NoteID int }
type GalleryMoveMsg struct{ NoteID, X, Y int }
