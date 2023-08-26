package labels

import (
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestGoIntoAddLabelView(t *testing.T) {
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
		{"should switch view into add label view", args{
			gotoMethod: func() {
				app.GoTo("labels-add", &jira.Issue{}, func() {}, jira.NewJiraApiMock(nil))
			},
			viewPredicate: func() bool {
				_, ok := app.GetApp().CurrentView().(*addLabelView)
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
