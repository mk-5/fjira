package app

import "github.com/gdamore/tcell"

type View interface {
	Init()
	Destroy()
}

type Drawable interface {
	Draw(screen tcell.Screen)
}

type Resizable interface {
	Resize(screenX, screenY int)
}

type System interface {
	Update()
}

type KeyListener interface {
	HandleKeyEvent(keyEvent *tcell.EventKey)
}

type Dirtyable struct {
	dirty bool
}

type Text struct {
	Dirtyable
	x     int
	y     int
	style tcell.Style
	text  string
}

type TextBox struct {
	x            int
	y            int
	x2           int
	y2           int
	text         string
	textStyle    tcell.Style
	bgStyle      tcell.Style
	borderStyle  tcell.Style
	borderTop    bool
	borderBottom bool
}

func NewText(x, y int, style tcell.Style, text string) *Text {
	return &Text{
		x: x, y: y, style: style, text: text,
	}
}

func (t *Text) Draw(screen tcell.Screen) {
	row := t.y
	col := t.x
	for _, r := range []rune(t.text) {
		if r == '\n' {
			row++
			col = t.x
			continue
		}
		screen.SetContent(col, row, r, nil, t.style)
		col++
	}
}

func (t *Text) ChangeText(newText string) {
	t.text = newText
}

func (t *TextBox) Draw(screen tcell.Screen) {
	if t.y2 < t.y {
		t.y, t.y2 = t.y2, t.y
	}
	if t.x2 < t.x {
		t.x, t.x2 = t.x2, t.x
	}

	// Fill background
	for row := t.y; row <= t.y2; row++ {
		for col := t.x; col <= t.x2; col++ {
			screen.SetContent(col, row, ' ', nil, t.bgStyle)
		}
	}

	// Draw borders
	for col := t.x; col <= t.x2; col++ {
		screen.SetContent(col, t.y, tcell.RuneHLine, nil, t.borderStyle)
		screen.SetContent(col, t.y2, tcell.RuneHLine, nil, t.borderStyle)
	}
	for row := t.y + 1; row < t.y2; row++ {
		screen.SetContent(t.x, row, tcell.RuneVLine, nil, t.borderStyle)
		screen.SetContent(t.x2, row, tcell.RuneVLine, nil, t.borderStyle)
	}

	// Only draw corners if necessary
	if t.y != t.y2 && t.x != t.x2 {
		screen.SetContent(t.x, t.y, tcell.RuneULCorner, nil, t.borderStyle)
		screen.SetContent(t.x2, t.y, tcell.RuneURCorner, nil, t.borderStyle)
		screen.SetContent(t.x, t.y2, tcell.RuneLLCorner, nil, t.borderStyle)
		screen.SetContent(t.x2, t.y2, tcell.RuneLRCorner, nil, t.borderStyle)
	}

	if t.text != "" {
		DrawText(screen, t.x+2, t.y+1, t.textStyle, t.text)
	}
}
