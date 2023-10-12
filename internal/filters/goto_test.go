package filters

import (
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestGoIntoFilters(t *testing.T) {
	app.InitTestApp(nil)
	RegisterGoTo()

	type args struct {
		gotoMethod    func()
		viewPredicate func() bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"should switch view into filters view", args{
			gotoMethod: func() { app.GoTo("filters", jira.NewJiraApiMock(nil)) },
			viewPredicate: func() bool {
				return app.CurrentScreenName() == "filters"
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
