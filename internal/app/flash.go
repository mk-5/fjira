package app

import (
	"fmt"
	"time"
)

func Error(message string) {
	app := GetApp()
	errorMessage := fmt.Sprintf("Error! -%s", message)
	errorStyle := DefaultStyle().Foreground(Color("alerts.error.foreground")).Background(Color("alerts.error.background"))
	errorBox := NewTextBox(app.ScreenX/2-len(errorMessage)/2, app.ScreenY-1, errorStyle, errorStyle, errorMessage)
	GetApp().AddFlash(errorBox, 5*time.Second)
}

func Success(message string) {
	app := GetApp()
	successMessage := fmt.Sprintf("Success! %s", message)
	successStyle := DefaultStyle().Foreground(Color("alerts.success.foreground")).Background(Color("alerts.success.background"))
	successBox := NewTextBox(app.ScreenX/2-len(successMessage)/2, app.ScreenY-1, successStyle, successStyle, successMessage)
	GetApp().AddFlash(successBox, 3*time.Second)
}
