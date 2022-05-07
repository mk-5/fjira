package fjira

import (
	"github.com/mk5/fjira/internal/app"
	"github.com/mk5/fjira/internal/jira"
	assert2 "github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

var tests = []struct {
	gotoMethod    func()
	viewPredicate func() bool
}{
	{
		gotoMethod: func() {
			goIntoProjectsSearch()
		},
		viewPredicate: func() bool {
			_, ok := app.GetApp().CurrentView().(*fjiraSearchProjectsView)
			return ok
		},
	},
	{
		gotoMethod: func() {
			goIntoIssuesSearch(&jira.JiraProject{})
		},
		viewPredicate: func() bool {
			_, ok := app.GetApp().CurrentView().(*fjiraSearchIssuesView)
			return ok
		},
	},
	{
		gotoMethod: func() {
			goIntoChangeAssignment(&jira.JiraIssue{})
		},
		viewPredicate: func() bool {
			_, ok := app.GetApp().CurrentView().(*fjiraAssignChangeView)
			return ok
		},
	},
	{
		gotoMethod: func() {
			goIntoChangeStatus(&jira.JiraIssue{})
		},
		viewPredicate: func() bool {
			_, ok := app.GetApp().CurrentView().(*fjiraStatusChangeView)
			return ok
		},
	},
	{
		gotoMethod: func() {
			goIntoIssueViewFetchIssue("ABC-123")
		},
		viewPredicate: func() bool {
			_, ok := app.GetApp().CurrentView().(*fjiraIssueView)
			return ok
		},
	},
}

func Test_shouldGotoCorrectScreens(t *testing.T) {
	assert := assert2.New(t)

	// given
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

	for _, test := range tests {
		// when
		test.gotoMethod()

		// then
		ok := test.viewPredicate()
		assert.True(ok, "Current view is invalid.")
	}
}
