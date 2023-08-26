package ui

import (
	"github.com/mk-5/fjira/internal/app"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func Test_goIntoValidScreen(t *testing.T) {
	RegisterGoTo()
	app.InitTestApp(nil)

	type args struct {
		gotoMethod    func()
		viewPredicate func() bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"should switch view into text writer view", args{
			gotoMethod: func() { app.GoTo("text-writer", &TextWriterArgs{}) },
			viewPredicate: func() bool {
				_, ok := app.GetApp().CurrentView().(*TextWriterView)
				return ok
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			tt.args.gotoMethod()

			// then
			ok := tt.args.viewPredicate()
			assert2.New(t).True(ok, "Current view is invalid.")
		})
	}
}
