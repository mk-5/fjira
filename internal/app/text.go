package app

import "github.com/gdamore/tcell"

type Text struct {
	x     int
	y     int
	style tcell.Style
	text  string
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
