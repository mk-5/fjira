package projects

import (
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestGoIntoProjectsSearch(t *testing.T) {
	app.InitTestApp(nil)
	RegisterGoto()
	type args struct {
		gotoMethod    func()
		viewPredicate func() bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"should switch view into search projects view", args{
			gotoMethod: func() { app.GoTo("projects", jira.NewJiraApiMock(nil)) },
			viewPredicate: func() bool {
				return app.CurrentScreenName() == "projects"
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
