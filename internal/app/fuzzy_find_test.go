package app

import (
	"bytes"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestFuzzyFind_Draw(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	screen.Init() //nolint:errcheck
	defer screen.Fini()

	type args struct {
		records []string
		query   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"should show valid results", args{records: []string{"abc"}, query: "abc"}, "abc"},
		{"should show valid results", args{records: []string{"Brzęczyszczykiewicz"}, query: "c"}, "Brzęczyszczykiewicz"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			fuzzyFind := NewFuzzyFind("test", tt.args.records)

			// when
			for _, key := range tt.args.query {
				fuzzyFind.HandleKeyEvent(tcell.NewEventKey(-1, key, tcell.ModNone))
			}
			fuzzyFind.Update()
			fuzzyFind.Resize(screen.Size())
			fuzzyFind.Draw(screen)
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
			assert.Contains(t, result, tt.want)
		})
	}
}
