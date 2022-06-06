package fjira

import (
	"github.com/mk5/fjira/internal/app"
	"github.com/mk5/fjira/internal/jira"
	assert2 "github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

const jiraIssueJson = `
{
    "expand": "renderedFields,names,schema,operations,editmeta,changelog,versionedRepresentations,customfield_10010.requestTypePractice",
    "id": "10011",
    "self": "https://test/rest/api/2/issue/10011",
    "key": "JWC-3",
    "fields": {
        "issuetype": {
            "id": "10013",
            "description": "A small, distinct piece of work.",
            "name": "Task",
            "subtask": false
        },
        "timespent": 14400,
        "project": {
            "id": "10003",
            "key": "JWC",
            "name": "JIRA WORK CHART",
            "projectTypeKey": "software"
        },
        "fixVersions": [],
        "aggregatetimespent": 14400,
        "resolutiondate": "2022-02-22T00:27:11.861+0100",
        "workratio": -1,
        "issuerestriction": {
            "issuerestrictions": {},
            "shouldDisplay": true
        },
        "lastViewed": "2022-02-22T00:27:17.356+0100",
        "created": "2021-10-02T22:34:22.521+0200",
        "aggregatetimeoriginalestimate": null,
        "timeestimate": 0,
        "versions": [],
        "issuelinks": [],
        "assignee": null,
        "updated": "2022-02-22T00:27:19.792+0100",
        "status": {
            "description": "",
            "name": "Done",
            "id": "10013"
        },
        "description": "Lorem ipsum",
        "summary": "Tutorial - create tutorial",
        "creator": {
            "accountId": "607f55ba074a0b006a6cb482",
            "emailAddress": "test@test.dev",
            "displayName": "Mateusz Kulawik",
            "active": true,
            "timeZone": "Europe/Warsaw",
            "accountType": "atlassian"
        },
        "subtasks": [],
        "reporter": {
			"accountId": "607f55ba074a0b006a6cb482",
            "emailAddress": "test@test.dev",
            "displayName": "Mateusz Kulawik",
            "active": true,
            "timeZone": "Europe/Warsaw",
            "accountType": "atlassian"
        }
    }
}`

func Test_shouldDisplayIssueView(t *testing.T) {
	fjira := CreateNewFjira(&fjiraSettings{})
	fjira.SetApi(jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.String(), "issue") {
			w.WriteHeader(200)
			w.Write([]byte(jiraIssueJson)) //nolint:errcheck
			return
		}
		if strings.Contains(r.URL.String(), "project") {
			w.WriteHeader(200)
			w.Write([]byte("[]")) //nolint:errcheck
		}
	}))

	assert := assert2.New(t)
	tests := []struct {
		name     string
		testFunc func()
	}{
		{"should crate valid issue view", func() {
			// when
			goIntoIssueView("ABC-123")
			view, ok := app.GetApp().CurrentView().(*fjiraIssueView)

			// then
			assert.True(ok, "Invalid view has been set")
			assert.Equal("JWC-3", view.issue.Key, "Invalid issue key")
			assert.Equal("Lorem ipsum", view.issue.Fields.Description, "Invalid issue description")
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc()
		})
	}
}
