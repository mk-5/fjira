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
	_ = screen.Init() //nolint:errcheck
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
			assert.NotEmpty(t, fuzzyFind.GetQuery())
			assert.Contains(t, result, tt.want)
		})
	}
}

func TestFuzzyFind_HandleKeyEvent(t *testing.T) {
	type args struct {
		ev []*tcell.EventKey
	}
	tests := []struct {
		name string
		args args
	}{
		{"should go up, and go down", args{ev: []*tcell.EventKey{
			tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone),
			tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone),
			tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone),
			tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone),
		}}},
		{"should go up, and go down using tab/tab-shift", args{ev: []*tcell.EventKey{
			tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone),
			tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone),
			tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone),
			tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone),
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			fuzzyFind := NewFuzzyFind("test", []string{"test1", "test2", "test3", "test4"})
			fuzzyFind.HandleKeyEvent(tcell.NewEventKey(0, 't', tcell.ModNone))
			fuzzyFind.Update()

			// when
			for _, key := range tt.args.ev {
				fuzzyFind.HandleKeyEvent(key)
			}

			// then
			assert.Equal(t, 2, fuzzyFind.selected)
		})
	}
}
