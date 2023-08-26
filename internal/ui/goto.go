package ui

import (
	"github.com/mk-5/fjira/internal/app"
)

func RegisterGoTo() {
	app.RegisterGoto("text-writer", func(args ...interface{}) {
		a := args[0].(*TextWriterArgs)
		view := NewTextWriterView(a)
		app.GetApp().SetView(view)
	})
}
