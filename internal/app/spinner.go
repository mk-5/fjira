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
			DefaultStyle(), DefaultStyle().Foreground(Color("spinner.accent")).Bold(true), DefaultStyle(),
		},
		textStyle:    DefaultStyle().Italic(true).Blink(true),
		spinnerIndex: new(int),
	}
}

func (t *SpinnerTCell) Draw(screen tcell.Screen) {
	screenX, screenY := screen.Size()
	index := (*t.spinnerIndex + 1) % len(t.spinner)
	*t.spinnerIndex = index
	row := screenY - 1
	col := screenX - len(t.spinner[index]) - 1
	if t.text != "" {
		col -= len(t.text) + 1
		DrawText(screen, screenX-1-len(t.text), screenY-1, t.textStyle, t.text)
	}
	for _, r := range t.spinner[index] {
		screen.SetContent(col, row, r, nil, t.styles[index])
		col++
	}
}
