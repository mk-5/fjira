package app

import (
	"github.com/gdamore/tcell/v2"
)

type Confirmation struct {
	Complete          chan bool
	message           string
	screenX           int
	screenY           int
	style             tcell.Style
	questionMarkStyle tcell.Style
}

const (
	Yes          = 'y'
	No           = 'n'
	QuestionMark = "? "
)

func Confirm(app *App, message string) bool {
	confirmation := newConfirmation(message)
	app.AddDrawable(confirmation)
	app.AddSystem(confirmation)
	if yesNo := <-confirmation.Complete; true {
		return yesNo
	}
	return false
}

func newConfirmation(message string) *Confirmation {
	return &Confirmation{
		Complete:          make(chan bool),
		message:           message,
		style:             DefaultStyle(),
		questionMarkStyle: DefaultStyle().Bold(true).Foreground(Color("finder.title")),
	}
}

func (c *Confirmation) Draw(screen tcell.Screen) {
	DrawText(screen, 0, c.screenY-2, c.questionMarkStyle, QuestionMark)
	DrawText(screen, 2, c.screenY-2, c.style, c.message)
}

func (c *Confirmation) Resize(screenX, screenY int) {
	c.screenX = screenX
	c.screenY = screenY
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
