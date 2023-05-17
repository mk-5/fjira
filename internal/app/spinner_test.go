package app

import (
	"bytes"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSpinnerTCell_Draw(t1 *testing.T) {
	tests := []struct {
		name string
	}{
		{"should draw the spinner"},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			// given
			screen := tcell.NewSimulationScreen("utf-8")
			_ = screen.Init() //nolint:errcheck
			defer screen.Fini()
			spin := NewSimpleSpinner()
			spin.text = "LOADING"

			// when
			spin.Draw(screen)

			// then
			var buffer bytes.Buffer
			contents, x, y := screen.GetContents()
			screen.Show()
			for i := 0; i < x*y; i++ {
				buffer.Write(contents[i].Bytes)
			}
			result := buffer.String()
			assert.Contains(t1, result, "LOADING")
		})
	}
}
