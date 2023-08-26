package app

import (
	"bytes"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNewTextBox(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()

	tests := []struct {
		name string
	}{
		{"should render text box"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			box := NewTextBox(0, 0, tcell.StyleDefault, tcell.StyleDefault, "abc")
			box.SetX(0)
			box.SetY(0)

			// when
			box.Draw(screen)
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
			assert.Contains(t, result, "└─────┘")
		})
	}
}
