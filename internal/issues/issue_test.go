package issues

import (
	"bytes"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	assert2 "github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
	"time"
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
        "labels": ["TestLabel"],
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
        },
 		"comment": {
            "comments": [
                {
                    "author": {
                        "displayName": "Mateusz Kulawik"
                    },
                    "body": "Comment 123-ABC",
                    "created": "2022-06-09T22:53:42.057+0200",
                    "updated": "2022-06-09T22:53:42.057+0200"
                }
            ],
            "maxResults": 1,
            "total": 1,
            "startAt": 0
        }
    }
}`

func Test_shouldDisplayIssueView(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()
	app.InitTestApp(screen)
	RegisterGoTo()
	api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.String(), "issue") {
			w.WriteHeader(200)
			_, _ = w.Write([]byte(jiraIssueJson)) //nolint:errcheck
			return
		}
		if strings.Contains(r.URL.String(), "project") {
			w.WriteHeader(200)
			_, _ = w.Write([]byte("[]")) //nolint:errcheck
		}
	})

	assert := assert2.New(t)
	tests := []struct {
		name     string
		screen   tcell.Screen
		testFunc func()
	}{
		{"should crate valid issue view", screen, func() {
			// when
			app.GoTo("issue", "ABC-123", nil, api)
			view, ok := app.GetApp().CurrentView().(*issueView)

			// then
			assert2.True(t, ok)

			// when
			view.Draw(screen)
			var buffer bytes.Buffer
			contents, x, y := screen.GetContents()
			screen.Show()
			for i := 0; i < x*y; i++ {
				if len(contents[i].Bytes) != 0 {
					buffer.Write(contents[i].Bytes)
				}
			}
			result := buffer.String()

			// then
			assert.True(ok, "Invalid view has been set")
			assert.Equal("JWC-3", view.issue.Key, "Invalid issue key")
			assert.Equal("Lorem ipsum", view.issue.Fields.Description, "Invalid issue description")
			assert.Contains(result, "JWC-3", "should contain issue number")
			assert.Contains(result, "Lorem ipsum", "should contain ticket description")
			assert.Contains(result, "Mateusz Kulawik", "should contain comment")
			assert.Contains(result, "Comment 123-ABC", "should contain comment author")
			assert.Contains(result, "TestLabel", "should contain labels")

			// and then
			view.Destroy()
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc()
		})
	}
}

func Test_issueView_ActionBar(t *testing.T) {
	app.InitTestApp(nil)
	app.RegisterGoto("issues-search", func(args ...interface{}) {
	})
	app.RegisterGoto("status-change", func(args ...interface{}) {
	})
	app.RegisterGoto("users-assign", func(args ...interface{}) {
	})
	app.RegisterGoto("labels-add", func(args ...interface{}) {
	})
	app.RegisterGoto("text-writer", func(args ...interface{}) {
	})

	type args struct {
		key           tcell.Key
		char          rune
		viewPredicate func() bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"should handle exit action", args{key: tcell.KeyEscape, viewPredicate: func() bool {
			return app.CurrentScreenName() == "issues-search"
		}}},
		{"should handle status change action", args{char: 's', viewPredicate: func() bool {
			return app.CurrentScreenName() == "status-change"
		}}},
		{"should handle assign user action", args{char: 'a', viewPredicate: func() bool {
			return app.CurrentScreenName() == "users-assign"
		}}},
		{"should handle comment action", args{char: 'c', viewPredicate: func() bool {
			return app.CurrentScreenName() == "text-writer"
		}}},
		{"should handle label action", args{char: 'l', viewPredicate: func() bool {
			return app.CurrentScreenName() == "labels-add"
		}}},
		{"should handle open action", args{char: 'o', viewPredicate: func() bool {
			return true
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			})
			issue := &jira.Issue{Id: "1", Key: "ABC-1"}
			view := NewIssueView(issue, func() {
				app.GoTo("issues-search", "ABC", nil, jira.NewJiraApiMock(nil))
			}, api).(*issueView)
			done := make(chan struct{})
			started := make(chan struct{})
			go func() {
				started <- struct{}{}
				view.handleIssueAction()
				done <- struct{}{}
			}()
			<-started
			<-time.NewTimer(100 * time.Millisecond).C

			// when
			view.HandleKeyEvent(tcell.NewEventKey(tt.args.key, tt.args.char, tcell.ModNone))
			<-done
			result := tt.args.viewPredicate()

			// then
			assert2.True(t, result)
		})
	}
}

func Test_issueView_doComment(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()
	app.InitTestApp(screen)

	tests := []struct {
		name string
	}{
		{"should run doComment api"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			done := make(chan bool)
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(``)) //nolint:errcheck
				done <- true
			})
			view := NewIssueView(&jira.Issue{Key: "test"}, nil, api).(*issueView)

			// when
			view.Init()
			go view.doComment(view.issue, "abcde")

			// then
			select {
			case <-done:
			case <-time.After(3 * time.Second):
				t.Fail()
			}
		})
	}
}

func Test_fjiraIssueView_HandleKeyEvent(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()
	app.InitTestApp(screen)

	tests := []struct {
		name string
	}{
		{"should process scrollUp&scrollDown"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			view := NewIssueView(&jira.Issue{Key: "test"}, nil, jira.NewJiraApiMock(nil)).(*issueView)
			view.fuzzyFind = app.NewFuzzyFind("test", []string{})
			view.scrollY = 0
			view.maxScrollY = 100

			// when
			view.HandleKeyEvent(tcell.NewEventKey(tcell.KeyDown, 'k', tcell.ModNone))

			// then
			assert2.Equal(t, 1, view.scrollY)

			// and when
			view.HandleKeyEvent(tcell.NewEventKey(tcell.KeyUp, 'k', tcell.ModNone))

			// then
			assert2.Equal(t, 0, view.scrollY)
		})
	}
}
