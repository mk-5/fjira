package app

import "github.com/gdamore/tcell/v2"

func InitTestApp(s tcell.SimulationScreen) *App {
	if s == nil {
		s = tcell.NewSimulationScreen("utf-8")
	}
	app := CreateNewAppWithScreen(s)
	return app
}
