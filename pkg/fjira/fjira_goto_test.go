package fjira

import (
	"github.com/mk5/fjira/internal/app"
	"github.com/mk5/fjira/internal/jira"
	assert2 "github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func init() {
	CreateNewFjira(jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.String(), "issue") {
			w.WriteHeader(200)
			w.Write([]byte("{}")) //nolint:errcheck
			return
		}
		if strings.Contains(r.URL.String(), "project") {
			w.WriteHeader(200)
			w.Write([]byte("[]")) //nolint:errcheck
		}
	}))
}

func Test_goIntoValidScreen(t *testing.T) {
	type args struct {
		gotoMethod    func()
		viewPredicate func() bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"should switch view into assignment change view", args{
			gotoMethod: func() { goIntoChangeAssignment(&jira.JiraIssue{}) },
			viewPredicate: func() bool {
				_, ok := app.GetApp().CurrentView().(*fjiraAssignChangeView)
				return ok
			},
		}},
		{"should switch view into search projects view", args{
			gotoMethod: func() { goIntoProjectsSearch() },
			viewPredicate: func() bool {
				_, ok := app.GetApp().CurrentView().(*fjiraSearchProjectsView)
				return ok
			},
		}},
		{"should switch view into search issues view", args{
			gotoMethod: func() { goIntoIssuesSearch(&jira.JiraProject{}) },
			viewPredicate: func() bool {
				_, ok := app.GetApp().CurrentView().(*fjiraSearchIssuesView)
				return ok
			},
		}},
		{"should switch view into change status view", args{
			gotoMethod: func() { goIntoChangeStatus(&jira.JiraIssue{}) },
			viewPredicate: func() bool {
				_, ok := app.GetApp().CurrentView().(*fjiraStatusChangeView)
				return ok
			},
		}},
		{"should switch view into issue view", args{
			gotoMethod: func() { goIntoIssueView("ABC-123") },
			viewPredicate: func() bool {
				_, ok := app.GetApp().CurrentView().(*fjiraIssueView)
				return ok
			},
		}},
		{"should switch view into comment view", args{
			gotoMethod: func() { goIntoCommentView(&jira.JiraIssue{}) },
			viewPredicate: func() bool {
				_, ok := app.GetApp().CurrentView().(*fjiraCommentView)
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
