package app

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell"
)

var (
	errorStyle   = DefaultStyle.Foreground(tcell.ColorRed).Background(tcell.ColorWhiteSmoke)
	successStyle = DefaultStyle.Foreground(tcell.ColorGreen).Background(tcell.ColorWhiteSmoke)
)

func Error(message string) {
	app := GetApp()
	errorMessage := fmt.Sprintf("Error! -%s", message)
	errorBox := NewTextBox(app.ScreenX/2-len(errorMessage)/2, app.ScreenY-1, errorStyle, errorStyle, errorMessage)
	GetApp().AddFlash(errorBox, 3*time.Second)
}

func Success(message string) {
	app := GetApp()
	successMessage := fmt.Sprintf("Success! %s", message)
	successBox := NewTextBox(app.ScreenX/2-len(successMessage)/2, app.ScreenY-1, successStyle, successStyle, successMessage)
	GetApp().AddFlash(successBox, 3*time.Second)
}
