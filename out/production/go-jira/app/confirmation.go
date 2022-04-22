package app

import (
	"github.com/gdamore/tcell"
)

type Confirmation struct {
	Complete chan bool
	message  string
}

const (
	Yes          = 'y'
	No           = 'n'
	QuestionMark = "? "
)

var (
	QuestionMarkStyle = tcell.StyleDefault.Bold(true).Foreground(tcell.ColorYellowGreen)
)

func Confirm(app *App, message string) bool {
	confirmation := newConfirmation(message)
	app.AddDrawable(confirmation)
	app.AddSystem(confirmation)
	select {
	case yesNo := <-confirmation.Complete:
		return yesNo
	}
}

func newConfirmation(message string) *Confirmation {
	return &Confirmation{
		Complete: make(chan bool),
		message:  message,
	}
}

func (c *Confirmation) Draw(screen tcell.Screen) {
	DrawText(screen, 0, 0, QuestionMarkStyle, QuestionMark)
	DrawText(screen, 2, 0, tcell.StyleDefault, c.message)
}

func (c *Confirmation) Update() {
	// do nothing
}

func (c *Confirmation) HandleKeyEvent(ev *tcell.EventKey) {
	if ev.Key() == tcell.KeyEscape {
		c.Complete <- false
		return
	}
	switch ev.Rune() {
	case Yes:
		c.Complete <- true
	case No:
		c.Complete <- false
	}
}
