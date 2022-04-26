package app

import (
	"fmt"
	"github.com/gdamore/tcell"
	"time"
)

var (
	errorStyle = DefaultStyle.Foreground(tcell.ColorRed)
)

func Error(message string) {
	app := GetApp()
	errorMessage := fmt.Sprintf("Error! -%s", message)
	errorBox := NewTextBox(app.ScreenX/2-len(errorMessage)/2, app.ScreenY-1, errorStyle, errorStyle, fmt.Sprintf("Error! -%s", message))
	GetApp().AddDrawable(errorBox)
	GetApp().KeepAlive(errorBox)
	ticker := time.NewTimer(3 * time.Second)
	go func() {
		<-ticker.C
		GetApp().UnKeepAlive(errorBox)
		GetApp().RemoveDrawable(errorBox)
		GetApp()
	}()
}
