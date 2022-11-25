package app

import (
	"bytes"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestConfirm(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()

	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
	}{
		{"should render confirmation message", args{message: "Do you want to confirm xxx?"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app := CreateNewAppWithScreen(screen)
			go Confirm(app, tt.args.message)
			<-time.NewTimer(100 * time.Millisecond).C

			// when
			app.Render()
			var buffer bytes.Buffer
			contents, x, y := screen.GetContents()
			screen.Show()
			for i := 0; i < x*y; i++ {
				if string(contents[i].Bytes) != "" {
					buffer.Write(contents[i].Bytes)
				}
			}
			result := strings.TrimSpace(buffer.String())

			// then
			assert.Contains(t, result, tt.args.message)
		})
	}
}

func TestConfirmation_HandleKeyEvent(t *testing.T) {
	type args struct {
		ev *tcell.EventKey
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"should process ESC key", args{ev: tcell.NewEventKey(tcell.KeyEsc, 0, 0)}, false},
		{"should process NO key", args{ev: tcell.NewEventKey(0, No, 0)}, false},
		{"should process YES key", args{ev: tcell.NewEventKey(0, Yes, 0)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			completeChan := make(chan bool)
			c := &Confirmation{
				Complete: completeChan,
				message:  "abc",
				screenX:  2,
				screenY:  2,
			}
			go c.HandleKeyEvent(tt.args.ev)
			result := <-completeChan

			// then
			assert.Equal(t, tt.want, result)
		})
	}
}
