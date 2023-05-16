package fjira

import (
	"bytes"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
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
			app.CreateNewAppWithScreen(screen)
			CreateNewFjira(&fjiraSettings{})
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(`[{"statuses" : [{"id": "1", "name": "Status1", "description": ""}, {"id": "2", "name": "xxx", "description": ""}]}]`))
			})
			_ = SetApi(api)
			view := NewIssuesSearchView(&jira.Project{Id: "TEST", Key: "TEST", Name: "TEST"})

			// when
			go view.runSelectStatus()
			<-time.NewTimer(700 * time.Millisecond).C
			query := "xxx"
			for _, key := range query {
				view.HandleKeyEvent(tcell.NewEventKey(-1, key, tcell.ModNone))
			}
			view.fuzzyFind.Update()
			<-time.NewTimer(700 * time.Millisecond).C // fuzzy debounce
			view.fuzzyFind.Update()
			<-time.NewTimer(700 * time.Millisecond).C
			view.HandleKeyEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
			view.Update()
			<-time.NewTimer(700 * time.Millisecond).C

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
			app.CreateNewAppWithScreen(screen)
			CreateNewFjira(&fjiraSettings{})
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(`[{"id": "U1", "displayName": "Bob"}, {"id": "U2", "displayName": "John"}]`))
			})
			_ = SetApi(api)
			view := NewIssuesSearchView(&jira.Project{Id: "TEST", Key: "TEST", Name: "TEST"})

			// when
			go view.runSelectUser()
			<-time.NewTimer(700 * time.Millisecond).C
			query := "John"
			for _, key := range query {
				view.fuzzyFind.HandleKeyEvent(tcell.NewEventKey(-1, key, tcell.ModNone))
			}
			view.fuzzyFind.Update()
			<-time.NewTimer(700 * time.Millisecond).C // fuzzy debounce
			view.fuzzyFind.Update()
			<-time.NewTimer(700 * time.Millisecond).C
			view.HandleKeyEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
			view.Update()
			<-time.NewTimer(700 * time.Millisecond).C

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
			app.CreateNewAppWithScreen(screen)
			CreateNewFjira(&fjiraSettings{})
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(`{"token":"","suggestions":[{"label":"SomethingElse","html":"<b></b>SomethingElse"},{"label":"TestLabel","html":"<b></b>TestLabel"},{"label":"Design","html":"<b></b>Design"},{"label":"Windows","html":"<b></b>Windows"}]}`))
			})
			_ = SetApi(api)
			view := NewIssuesSearchView(&jira.Project{Id: "TEST", Key: "TEST", Name: "TEST"})

			// when
			go view.runSelectLabel()
			<-time.NewTimer(700 * time.Millisecond).C
			query := "de"
			for _, key := range query {
				view.fuzzyFind.HandleKeyEvent(tcell.NewEventKey(-1, key, tcell.ModNone))
			}
			view.fuzzyFind.Update()
			<-time.NewTimer(700 * time.Millisecond).C // fuzzy debounce
			view.fuzzyFind.Update()
			<-time.NewTimer(700 * time.Millisecond).C
			view.HandleKeyEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
			view.Update()
			<-time.NewTimer(700 * time.Millisecond).C

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

	tests := []struct {
		name string
	}{
		{"should run&select board view"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.CreateNewAppWithScreen(screen)
			CreateNewFjira(&fjiraSettings{})
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
			_ = SetApi(api)
			view := NewIssuesSearchView(&jira.Project{Id: "TEST", Key: "TEST", Name: "TEST"})

			// when
			go view.runSelectBoard()
			<-time.NewTimer(300 * time.Millisecond).C
			query := "Gen"
			for _, key := range query {
				view.HandleKeyEvent(tcell.NewEventKey(-1, key, tcell.ModNone))
			}
			view.Update()
			view.Update()
			view.HandleKeyEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
			<-time.NewTimer(300 * time.Millisecond).C
			_, switchedToBoardsView := app.GetApp().CurrentView().(*boardView)

			// then
			assert.True(t, switchedToBoardsView)
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
			app.CreateNewAppWithScreen(screen)
			CreateNewFjira(&fjiraSettings{})
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, err := w.Write([]byte(`{"token":"","suggestions":[{"label":"SomethingElse","html":"<b></b>SomethingElse"},{"label":"TestLabel","html":"<b></b>TestLabel"},{"label":"Design","html":"<b></b>Design"},{"label":"Windows","html":"<b></b>Windows"}]}`))
				println(err)
			})
			_ = SetApi(api)
			view := NewIssuesSearchView(&jira.Project{Id: "TEST", Key: "TEST", Name: "TEST"})

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
			CreateNewFjira(&fjiraSettings{})
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
			_ = SetApi(api)
			view := NewIssuesSearchView(&jira.Project{Id: "TEST", Key: "TEST", Name: "TEST"})

			// when
			boards := view.findBoards()

			// then
			assert.Equal(t, 1, len(boards))
			assert.Equal(t, "GEN board", boards[0].Name)
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
			app.CreateNewAppWithScreen(screen)
			CreateNewFjira(&fjiraSettings{})
			view := NewIssuesSearchView(&jira.Project{Id: "TEST", Key: "TEST", Name: "TEST"})

			// when
			view.goBack()
			<-time.After(200 * time.Millisecond)

			// then
			_, ok := app.GetApp().CurrentView().(*fjiraSearchProjectsView)
			assert.True(t, ok)
		})
	}
}

func Test_fjiraSearchIssuesView_goBackWithCustomJql(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()

	tests := []struct {
		name string
	}{
		{"should go back with custom jql"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			app.CreateNewAppWithScreen(screen)
			CreateNewFjira(&fjiraSettings{})
			view := NewIssuesSearchViewWithCustomJql("test jql")

			// when
			view.goBack()
			<-time.After(200 * time.Millisecond)

			// then
			_, ok := app.GetApp().CurrentView().(*fjiraJqlSearchView)
			assert.True(t, ok)
		})
	}
}
