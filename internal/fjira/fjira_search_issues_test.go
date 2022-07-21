package fjira

import (
	"bytes"
	"github.com/gdamore/tcell/v2"
	"github.com/mk5/fjira/internal/app"
	"github.com/mk5/fjira/internal/jira"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestNewIssuesSearchView(t *testing.T) {
	type args struct {
		project *jira.Project
	}
	tests := []struct {
		name string
		args args
	}{
		{"should create new search issues view", args{project: &jira.Project{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, NewIssuesSearchView(tt.args.project), "NewIssuesSearchView(%v)", tt.args.project)
		})
	}
}

func Test_fjiraSearchIssuesView_Destroy(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should run destroy without error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := NewIssuesSearchView(&jira.Project{})
			view.Destroy()
		})
	}
}

func Test_fjiraSearchIssuesView_Init(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	screen.Init() //nolint:errcheck
	defer screen.Fini()

	type args struct {
		screen tcell.Screen
	}
	tests := []struct {
		name string
		args args
	}{
		{"should draw issues search view", args{screen: screen}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.CreateNewAppWithScreen(screen)
			CreateNewFjira(&fjiraSettings{})
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(`{
    "expand": "schema,names",
    "startAt": 0,
    "maxResults": 100,
    "total": 3,
    "issues": [
        {
            "key": "ISSUE-1",
            "fields": {
                "summary": "Issue summary 1",
                "issuetype": {
                    "name": "Task"
                },
                "project": {
                    "key": "ISSUE",
                    "name": "Issues"
                },
                "reporter": {
                    "emailAddress": "test@test.pl",
                    "displayName": "Test",
                    "active": true,
                    "timeZone": "Europe/Warsaw",
                    "accountType": "atlassian"
                },
                "assignee": {
                    "emailAddress": "test@test.pl",
                    "displayName": "Test",
                    "active": true,
                    "timeZone": "Europe/Warsaw",
                    "accountType": "atlassian"
                },
                "status": {
                    "name": "In Progress"
                }
            }
        },
        {
            "key": "ISSUE-2",
            "fields": {
                "summary": "Issue summary 2",
                "issuetype": {
                    "name": "Task"
                },
                "project": {
                    "key": "ISSUE",
                    "name": "Issues"
                },
                "reporter": {
                    "emailAddress": "test@test.pl",
                    "displayName": "Test",
                    "active": true,
                    "timeZone": "Europe/Warsaw",
                    "accountType": "atlassian"
                },
                "assignee": {
					"emailAddress": "test@test.pl",
                    "displayName": "Test",
                    "active": true,
                    "timeZone": "Europe/Warsaw",
                    "accountType": "atlassian"
                },
                "status": {
                    "name": "In Progress"
                }
            }
        },
        {
            "key": "ISSUE-3",
            "fields": {
                "summary": "Issue test 3",
                "issuetype": {
                    "name": "Task"
                },
                "project": {
                    "key": "ISSUE",
                    "name": "Issues"
                },
                "reporter": {
                    "emailAddress": "test@test.pl",
                    "displayName": "Test",
                    "active": true,
                    "timeZone": "Europe/Warsaw",
                    "accountType": "atlassian"
                },
                "assignee": {
					"emailAddress": "test@test.pl",
                    "displayName": "Test",
                    "active": true,
                    "timeZone": "Europe/Warsaw",
                    "accountType": "atlassian"
                },
                "status": {
                    "name": "In Progress"
                }
            }
        }
    ]
}`)) //nolint:errcheck
			})
			_ = SetApi(api)
			view := NewIssuesSearchView(&jira.Project{Key: "TEST", Name: "TEST"})

			// when
			view.Init()
			<-time.NewTimer(1 * time.Second).C
			query := "summary"
			for _, key := range query {
				view.HandleKeyEvent(tcell.NewEventKey(-1, key, tcell.ModNone))
			}
			view.Update()
			view.Resize(screen.Size())
			<-time.NewTimer(1 * time.Second).C
			view.Update()
			view.Draw(tt.args.screen)
			<-time.NewTimer(100 * time.Millisecond).C

			// and when
			var buffer bytes.Buffer
			contents, x, y := tt.args.screen.(tcell.SimulationScreen).GetContents()
			tt.args.screen.Show()
			for i := 0; i < x*y; i++ {
				if string(contents[i].Bytes) != " " {
					buffer.Write(contents[i].Bytes)
				}
			}
			result := buffer.String()

			// then
			assert.Contains(t, result, "ISSUE-1")
			assert.Contains(t, result, "ISSUE-2")
			assert.NotContains(t, result, "ISSUE-3")
		})
	}
}

func Test_fjiraSearchIssuesView_queryHasIssueFormat(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  bool
	}{
		{"should recognize jira issue key", "ISS-1", true},
		{"should recognize jira issue key", "ISS1", false},
		{"should recognize jira issue key", "", false},
		{"should recognize jira issue key", "TEST-312313", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			view := NewIssuesSearchView(&jira.Project{Key: "TEST", Name: "TEST"})

			// when
			view.currentQuery = tt.query

			// then
			assert.Equalf(t, tt.want, view.queryHasIssueFormat(), "queryHasIssueFormat()")
		})
	}
}

func Test_fjiraSearchIssuesView_runSelectStatus(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should run select status view"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.CreateNewAppWithScreen(tcell.NewSimulationScreen("utf-8"))
			CreateNewFjira(&fjiraSettings{})
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, err := w.Write([]byte(`[{"statuses" : [{"id": "1", "name": "Status1", "description": ""}, {"id": "2", "name": "xxx", "description": ""}]}]`))
				println(err)
			})
			_ = SetApi(api)
			view := NewIssuesSearchView(&jira.Project{Id: "TEST", Key: "TEST", Name: "TEST"})

			// when
			go view.runSelectStatus()
			<-time.NewTimer(100 * time.Millisecond).C
			query := "xxx"
			for _, key := range query {
				view.HandleKeyEvent(tcell.NewEventKey(-1, key, tcell.ModNone))
			}
			view.Update()
			view.Update()
			view.HandleKeyEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
			<-time.NewTimer(100 * time.Millisecond).C

			// then
			assert.NotNil(t, searchForStatus)
			assert.Equal(t, "xxx", searchForStatus.Name)
		})
	}
}

func Test_fjiraSearchIssuesView_runSelectUser(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should run select user view"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.CreateNewAppWithScreen(tcell.NewSimulationScreen("utf-8"))
			CreateNewFjira(&fjiraSettings{})
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, err := w.Write([]byte(`[{"id": "U1", "displayName": "Bob"}, {"id": "U2", "displayName": "John"}]`))
				println(err)
			})
			_ = SetApi(api)
			view := NewIssuesSearchView(&jira.Project{Id: "TEST", Key: "TEST", Name: "TEST"})

			// when
			go view.runSelectUser()
			<-time.NewTimer(100 * time.Millisecond).C
			query := "John"
			for _, key := range query {
				view.HandleKeyEvent(tcell.NewEventKey(-1, key, tcell.ModNone))
			}
			view.Update()
			view.Update()
			view.HandleKeyEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
			<-time.NewTimer(100 * time.Millisecond).C

			// then
			assert.NotNil(t, searchForUser)
			assert.Equal(t, "John", searchForUser.DisplayName)
		})
	}
}
