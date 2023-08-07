package app

import (
	"bytes"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestFlush(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()
	app := CreateNewAppWithScreen(screen)

	type args struct {
		message string
	}

	tests := []struct {
		name string
		args args
	}{
		{"should render flush messages", args{message: "Flush message!"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			Success(tt.args.message)
			app.flash[0].Draw(screen)

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

			// and when
			Error(tt.args.message)
			app.flash[0].Draw(screen)

			buffer.Reset()
			contents, x, y = screen.GetContents()
			screen.Show()
			for i := 0; i < x*y; i++ {
				if string(contents[i].Bytes) != "" {
					buffer.Write(contents[i].Bytes)
				}
			}
			result = strings.TrimSpace(buffer.String())

			// then
			assert.Contains(t, result, tt.args.message)
		})
	}
}
