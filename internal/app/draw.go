package app

import (
	"github.com/gdamore/tcell/v2"
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
		if s != nil {
			s.SetContent(col, row, r, nil, style)
		}
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

func DrawBox(screen tcell.Screen, x1, y1, x2, y2 int, style tcell.Style) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	// Draw borders
	for col := x1; col <= x2; col++ {
		screen.SetContent(col, y1, tcell.RuneHLine, nil, style)
		screen.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		screen.SetContent(x1, row, tcell.RuneVLine, nil, style)
		screen.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		screen.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		screen.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		screen.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		screen.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
}
