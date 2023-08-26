package issues

import (
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestGoIntoIssues(t *testing.T) {
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
		{"should switch view into jql view", args{
			gotoMethod: func() { app.GoTo("jql", jira.NewJiraApiMock(nil)) },
			viewPredicate: func() bool {
				return app.CurrentScreenName() == "jql"
			},
		}},
		{"should switch view into search issues view", args{
			gotoMethod: func() { app.GoTo("issues-search", "ABC", nil, jira.NewJiraApiMock(nil)) },
			viewPredicate: func() bool {
				return app.CurrentScreenName() == "issues-search"
			},
		}},
		{"should switch view into issue view", args{
			gotoMethod: func() { app.GoTo("issue", "ABC-123", nil, jira.NewJiraApiMock(nil)) },
			viewPredicate: func() bool {
				return app.CurrentScreenName() == "issue"
			},
		}},
		{"should switch view into issues view with jql", args{
			gotoMethod: func() { app.GoTo("issues-search-jql", "test jql", jira.NewJiraApiMock(nil)) },
			viewPredicate: func() bool {
				return app.CurrentScreenName() == "issues-search-jql"
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
