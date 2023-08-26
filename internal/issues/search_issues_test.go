package issues

import (
	"bytes"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	"github.com/mk-5/fjira/internal/projects"
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
			assert.NotNil(t, NewIssuesSearchView(tt.args.project, nil, jira.NewJiraApiMock(nil)), "NewIssuesSearchView(%v)", tt.args.project)
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
			view := NewIssuesSearchView(&jira.Project{}, nil, jira.NewJiraApiMock(nil))
			view.Destroy()
		})
	}
}

func Test_fjiraSearchIssuesView_Init(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
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
			app.InitTestApp(screen)
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
			view := NewIssuesSearchView(&jira.Project{Key: "TEST", Name: "TEST"}, nil, api).(*searchIssuesView)

			// when
			view.Init()
			for view.fuzzyFind == nil {
				<-time.After(10 * time.Millisecond)
			}
			view.fuzzyFind.SetDebounceDisabled(true)
			query := "summary"
			for _, key := range query {
				view.HandleKeyEvent(tcell.NewEventKey(-1, key, tcell.ModNone))
			}
			view.Resize(screen.Size())
			view.Update()
			view.Draw(tt.args.screen)

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
			view := NewIssuesSearchView(&jira.Project{Key: "TEST", Name: "TEST"}, nil, jira.NewJiraApiMock(nil)).(*searchIssuesView)

			// when
			view.currentQuery = tt.query

			// then
			assert.Equalf(t, tt.want, view.queryHasIssueFormat(), "queryHasIssueFormat()")
		})
	}
}

func Test_fjiraSearchIssuesView_runSelectStatus(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()

	tests := []struct {
		name string
	}{
		{"should run select status view"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.InitTestApp(screen)
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(`[{"statuses" : [{"id": "1", "name": "Status1", "description": ""}, {"id": "2", "name": "xxx", "description": ""}]}]`))
			})
			view := NewIssuesSearchView(&jira.Project{Id: "TEST", Key: "TEST", Name: "TEST"}, nil, api).(*searchIssuesView)

			// when
			done := make(chan struct{})
			go func() {
				view.runSelectStatus()
				done <- struct{}{}
			}()
			for {
				if view.fuzzyFind != nil {
					break
				}
				<-time.NewTimer(10 * time.Millisecond).C
			}
			query := "xxx"
			for _, key := range query {
				view.fuzzyFind.HandleKeyEvent(tcell.NewEventKey(-1, key, tcell.ModNone))
			}
			view.fuzzyFind.Update()
			view.fuzzyFind.HandleKeyEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
			<-done

			// then
			assert.NotNil(t, searchForStatus)
			assert.Equal(t, "xxx", searchForStatus.Name)
		})
	}
}

func Test_fjiraSearchIssuesView_runSelectUser(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()

	tests := []struct {
		name string
	}{
		{"should run select user view"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.InitTestApp(screen)
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(`[{"id": "U1", "displayName": "Bob"}, {"id": "U2", "displayName": "John"}]`))
			})
			view := NewIssuesSearchView(&jira.Project{Id: "TEST", Key: "TEST", Name: "TEST"}, nil, api).(*searchIssuesView)

			// when
			done := make(chan struct{})
			go func() {
				view.runSelectUser()
				done <- struct{}{}
			}()
			for {
				if view.fuzzyFind != nil {
					break
				}
				<-time.NewTimer(10 * time.Millisecond).C
			}
			query := "John"
			for _, key := range query {
				view.fuzzyFind.HandleKeyEvent(tcell.NewEventKey(-1, key, tcell.ModNone))
			}
			view.fuzzyFind.Update()
			view.fuzzyFind.HandleKeyEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
			<-done

			// then
			assert.NotNil(t, searchForUser)
			assert.Equal(t, "John", searchForUser.DisplayName)
		})
	}
}

func Test_fjiraSearchIssuesView_runSelectLabel(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()

	tests := []struct {
		name string
	}{
		{"should run select label view"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.InitTestApp(screen)
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(`{"token":"","suggestions":[{"label":"SomethingElse","html":"<b></b>SomethingElse"},{"label":"TestLabel","html":"<b></b>TestLabel"},{"label":"Design","html":"<b></b>Design"},{"label":"Windows","html":"<b></b>Windows"}]}`))
			})
			view := NewIssuesSearchView(&jira.Project{Id: "TEST", Key: "TEST", Name: "TEST"}, nil, api).(*searchIssuesView)

			// when
			done := make(chan struct{})
			go func() {
				view.runSelectLabel()
				done <- struct{}{}
			}()
			for view.fuzzyFind == nil {
				<-time.After(10 * time.Millisecond)
			}
			view.fuzzyFind.SetDebounceDisabled(true)
			query := "de"
			for _, key := range query {
				view.fuzzyFind.HandleKeyEvent(tcell.NewEventKey(-1, key, tcell.ModNone))
				view.Update()
			}
			view.fuzzyFind.HandleKeyEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
			view.Update()
			<-done

			// then
			assert.NotNil(t, searchForLabel)
			assert.Equal(t, "Design", searchForLabel)
		})
	}
}

func Test_fjiraSearchIssuesView_runSelectBoard(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()
	app.RegisterGoto("boards", func(args ...interface{}) {
	})

	tests := []struct {
		name string
	}{
		{"should run&select board view"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.InitTestApp(screen)
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, err := w.Write([]byte(`{
    "maxResults": 50,
    "startAt": 0,
    "total": 1,
    "isLast": true,
    "values": [
        {
            "id": 1,
            "self": "https://test/rest/agile/1.0/board/1",
            "name": "GEN board",
            "type": "kanban",
            "location": {
                "projectId": 10000,
                "displayName": "General (GEN)",
                "projectName": "General",
                "projectKey": "GEN",
                "projectTypeKey": "software",
                "avatarURI": "https://test/rest/api/2/universal_avatar/view/type/project/avatar/10416?size=small",
                "name": "General (GEN)"
            }
        }
    ]
}
`))
				println(err)
			})
			view := NewIssuesSearchView(&jira.Project{Id: "TEST", Key: "TEST", Name: "TEST"}, nil, api).(*searchIssuesView)

			// when
			done := make(chan struct{})
			go func() {
				view.runSelectBoard()
				done <- struct{}{}
			}()
			for view.fuzzyFind == nil {
				<-time.After(10 * time.Millisecond)
			}
			query := "Gen"
			for _, key := range query {
				view.HandleKeyEvent(tcell.NewEventKey(-1, key, tcell.ModNone))
			}
			view.HandleKeyEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
			<-done

			// then
			assert.True(t, app.CurrentScreenName() == "boards")
		})
	}
}

func Test_fjiraSearchIssuesView_findLabels(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()

	tests := []struct {
		name string
	}{
		{"should find project labels"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.InitTestApp(screen)
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, err := w.Write([]byte(`{"token":"","suggestions":[{"label":"SomethingElse","html":"<b></b>SomethingElse"},{"label":"TestLabel","html":"<b></b>TestLabel"},{"label":"Design","html":"<b></b>Design"},{"label":"Windows","html":"<b></b>Windows"}]}`))
				println(err)
			})
			view := NewIssuesSearchView(&jira.Project{Id: "TEST", Key: "TEST", Name: "TEST"}, nil, api).(*searchIssuesView)

			// when
			view.findLabels("")

			// then: should contain 4 labels + "All" label
			assert.Equal(t, 5, len(view.labels))
			assert.Contains(t, view.labels, "SomethingElse")
			assert.Contains(t, view.labels, "TestLabel")
			assert.Contains(t, view.labels, "Windows")
			assert.Contains(t, view.labels, "Design")
		})
	}
}

func Test_fjiraSearchIssuesView_findBoards(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()

	tests := []struct {
		name string
	}{
		{"should find project labels"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.CreateNewAppWithScreen(screen)
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, err := w.Write([]byte(`{
    "maxResults": 50,
    "startAt": 0,
    "total": 1,
    "isLast": true,
    "values": [
        {
            "id": 1,
            "self": "https://test/rest/agile/1.0/board/1",
            "name": "GEN board",
            "type": "kanban",
            "location": {
                "projectId": 10000,
                "displayName": "General (GEN)",
                "projectName": "General",
                "projectKey": "GEN",
                "projectTypeKey": "software",
                "avatarURI": "https://test/rest/api/2/universal_avatar/view/type/project/avatar/10416?size=small",
                "name": "General (GEN)"
            }
        }
    ]
}
`))
				println(err)
			})
			view := NewIssuesSearchView(&jira.Project{Id: "TEST", Key: "TEST", Name: "TEST"}, nil, api).(*searchIssuesView)

			// when
			bs := view.findBoards()

			// then
			assert.Equal(t, 1, len(bs))
			assert.Equal(t, "GEN board", bs[0].Name)
		})
	}
}

func Test_fjiraSearchIssuesView_goBack(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()

	tests := []struct {
		name string
	}{
		{"should go back"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			app.InitTestApp(screen)
			done := make(chan struct{})
			view := NewIssuesSearchView(&jira.Project{Id: "TEST", Key: "TEST", Name: "TEST"}, func() {
				done <- struct{}{}
			}, jira.NewJiraApiMock(nil)).(*searchIssuesView)

			// when
			go view.goBack()

			// then
			<-done
			assert.True(t, true)
		})
	}
}

func TestNewIssuesSearchView_fjiraIssueView_goIntoIssueVIew(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()

	type args struct {
		project *jira.Project
	}
	tests := []struct {
		name string
		args args
	}{
		{"should go into issue view, and preserve custom JQL via custom goBackFn", args{project: &jira.Project{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				body := `
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
            "name": "TEST",
            "projectTypeKey": "software"
        },
        "lastViewed": "2022-02-22T00:27:17.356+0100",
        "created": "2021-10-02T22:34:22.521+0200",
        "issuelinks": [],
        "assignee": null,
        "updated": "2022-02-22T00:27:19.792+0100",
        "status": {
            "description": "",
            "name": "Done",
            "id": "10013"
        },
		"labels": ["TestLabel"],
        "description": "Lorem ipsum",
        "summary": "Tutorial - create tutorial",
        "creator": {
        },
        "subtasks": [],
        "reporter": {
        },
 		"comment": {}
    }
}
`
				w.Write([]byte(body)) //nolint:errcheck
			})
			a := app.InitTestApp(screen)
			view := NewIssuesSearchView(tt.args.project, nil, api).(*searchIssuesView)

			// when
			view.goToIssueView("ABC")
			a.SetView(projects.NewProjectsSearchView(api))
			view.customJql = "TEST"
			view.goToIssueView("ABC")

			// then
			if v, ok := app.GetApp().CurrentView().(*issueView); ok {
				assert.True(t, v.goBackFn != nil)
			}
		})
	}
}
