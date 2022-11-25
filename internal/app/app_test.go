package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestApp(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should create app without error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screen := tcell.NewSimulationScreen("utf-8")
			_ = screen.Init() //nolint:errcheck

			// when
			a := &App{
				screen:          screen,
				spinnerIndex:    0,
				keyEvent:        make(chan *tcell.EventKey),
				runOnAppRoutine: make([]func(), 0, 64),
				drawables:       make([]Drawable, 0, 256),
				systems:         make([]System, 0, 128),
				flash:           make([]Drawable, 0, 5),
				keepAlive:       make(map[interface{}]bool),
				dirty:           make(chan bool),
			}
			go a.Start()
			<-time.NewTimer(100 * time.Millisecond).C
			a.Quit()

			// then
			assert.True(t, a.quit)
		})
	}
}

func TestApp_KeepAlive(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()

	tests := []struct {
		name string
	}{
		{"should mark drawables as keep-alive"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screen := tcell.NewSimulationScreen("utf-8")
			_ = screen.Init() //nolint:errcheck
			defer screen.Fini()

			// given
			a := &App{
				screen:          screen,
				spinnerIndex:    0,
				keyEvent:        make(chan *tcell.EventKey),
				runOnAppRoutine: make([]func(), 0, 64),
				drawables:       make([]Drawable, 0, 256),
				systems:         make([]System, 0, 128),
				flash:           make([]Drawable, 0, 5),
				keepAlive:       make(map[interface{}]bool),
				dirty:           make(chan bool),
			}
			drawable := NewText(0, 0, tcell.StyleDefault, "test")
			a.AddDrawable(drawable)

			// when
			a.KeepAlive(drawable)

			// then
			assert.Equal(t, true, a.keepAlive[drawable])

			// and then
			a.UnKeepAlive(drawable)

			// then
			assert.Equal(t, false, a.keepAlive[drawable])
		})
	}
}
