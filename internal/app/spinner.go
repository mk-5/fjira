package app

import "github.com/gdamore/tcell/v2"

type SpinnerTCell struct {
	spinner      []string
	text         string
	textStyle    tcell.Style
	styles       []tcell.Style
	spinnerIndex *int
}

func NewSimpleSpinner() *SpinnerTCell {
	return &SpinnerTCell{
		spinner: []string{".....", "....", ".."},
		styles: []tcell.Style{
			tcell.StyleDefault, tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true), tcell.StyleDefault,
		},
		textStyle:    tcell.StyleDefault.Italic(true).Blink(true),
		spinnerIndex: new(int),
	}
}

func (t *SpinnerTCell) Draw(screen tcell.Screen) {
	screenX, screenY := screen.Size()
	index := (*t.spinnerIndex + 1) % len(t.spinner)
	*t.spinnerIndex = index
	row := screenY - 2
	col := screenX - len(t.spinner[index]) - 1
	for _, r := range []rune(t.spinner[index]) {
		screen.SetContent(col, row, r, nil, t.styles[index])
		col++
	}
	if t.text != "" {
		DrawText(screen, screenX-1-len(t.text), screenY-3, t.textStyle, t.text)
	}
}
