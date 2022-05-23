package app

import (
	"github.com/gdamore/tcell"
)

const (
	EmptyLine = ""
)

func DrawText(screen tcell.Screen, x, y int, style tcell.Style, text string) {
	row := y
	col := x
	for _, r := range []rune(text) {
		if r == '\n' {
			row++
			col = x
			continue
		}
		screen.SetContent(col, row, r, nil, style)
		col++
	}
}

func DrawTextLimited(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) int {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		if r == '\n' {
			row++
			col = x1
			continue
		}
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
	return row
}
