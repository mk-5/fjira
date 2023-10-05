package boards

import (
	"bytes"
	"encoding/json"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	"github.com/stretchr/testify/assert"
	"net/http"
	"regexp"
	"testing"
)

func TestNewBoardView(t *testing.T) {
	type args struct {
		project            *jira.Project
		boardConfiguration *jira.BoardConfiguration
	}
	tests := []struct {
		name string
		args args
	}{
		{"should create a new board view", args{boardConfiguration: &jira.BoardConfiguration{}, project: &jira.Project{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, NewBoardView(tt.args.project, tt.args.boardConfiguration, "", nil), "NewBoardView(%v, %v)", tt.args.project, tt.args.boardConfiguration)
		})
	}
}

func Test_boardView_Destroy(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should run destroy without error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBoardView(&jira.Project{}, &jira.BoardConfiguration{}, "", nil)
			b.Destroy()
		})
	}
}

func Test_boardView_Draw(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()
	type args struct {
		screen tcell.Screen
	}
	app.InitTestApp(screen)
	tests := []struct {
		name string
		args args
	}{
		{"should draw board view", args{screen: screen}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			boardJson := `
{
    "id": 1,
    "name": "GEN board",
    "type": "kanban",
    "self": "https://test.net/rest/agile/1.0/board/1/configuration",
    "location": {
        "type": "project",
        "key": "GEN",
        "id": "10000",
        "self": "https://test.net/rest/api/2/project/10000",
        "name": "General"
    },
    "filter": {
        "id": "10000",
        "self": "https://test.net/rest/api/2/filter/10000"
    },
    "subQuery": {
        "query": "fixVersion in unreleasedVersions() OR fixVersion is EMPTY"
    },
    "columnConfig": {
        "columns": [
            {
                "name": "COL1",
                "statuses": []
            },
            {
                "name": "COL2",
                "statuses": [
                    {
                        "id": "10000",
                        "self": "https://test.net/rest/api/2/status/10000"
                    }
                ]
            },
            {
                "name": "COL3",
                "statuses": [
                    {
                        "id": "10001",
                        "self": "https://test.net/rest/api/2/status/10001"
                    }
                ]
            },
            {
                "name": "COL4",
                "statuses": [
                    {
                        "id": "10002",
                        "self": "https://test.net/rest/api/2/status/3"
                    }
                ]
            },
            {
                "name": "COL5",
                "statuses": [
                    {
                        "id": "10003",
                        "self": "https://test.net/rest/api/2/status/10002"
                    }
                ]
            }
        ],
        "constraintType": "issueCount"
    },
    "ranking": {
        "rankCustomFieldId": 10019
    }
}
`
			var board jira.BoardConfiguration
			_ = json.Unmarshal([]byte(boardJson), &board)
			view := NewBoardView(&jira.Project{}, &board, "", nil).(*boardView)
			view.issues = []jira.Issue{
				{Id: "1", Key: "GEN-1", Fields: jira.IssueFields{Status: jira.Status{Id: "10000"}}},
				{Id: "2", Key: "GEN-2", Fields: jira.IssueFields{Status: jira.Status{Id: "10001"}}},
				{Id: "3", Key: "GEN-3", Fields: jira.IssueFields{Status: jira.Status{Id: "10002"}}},
				{Id: "4", Key: "GEN-4", Fields: jira.IssueFields{Status: jira.Status{Id: "10003"}}},
			}

			// when
			view.SetColumnSize(10)
			view.Update()
			screenX, screenY := tt.args.screen.Size()
			view.Resize(screenX, screenY)
			view.Draw(tt.args.screen)
			var buffer bytes.Buffer
			contents, x, y := tt.args.screen.(tcell.SimulationScreen).GetContents()
			tt.args.screen.Show()
			for i := 0; i < x*y; i++ {
				if len(contents[i].Bytes) != 0 {
					buffer.Write(contents[i].Bytes)
				}
			}
			result := regexp.MustCompile(`\s+`).ReplaceAllString(buffer.String(), " ")

			// then
			assert.NotContains(t, result, "COL1")
			assert.Contains(t, result, "COL2")
			assert.Contains(t, result, "COL3")
			assert.Contains(t, result, "COL4")
			assert.Contains(t, result, "COL5")
			assert.Contains(t, result, "GEN-1")
			assert.Contains(t, result, "GEN-2")
			assert.Contains(t, result, "GEN-3")
			assert.Contains(t, result, "GEN-4")
		})
	}
}

func Test_boardView_HandleKeyEvent(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()
	type args struct {
		ev []*tcell.EventKey
	}
	tests := []struct {
		name        string
		args        args
		wantCursorX int
		wantCursorY int
	}{
		{"should handle key events and move cursor, and move issue", args{ev: []*tcell.EventKey{
			tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone),
			tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone),
			tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone),
			tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModNone),
			tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone),
			tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone),
			tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone),
		}}, 1, 1},
		{"should handle key events and move cursor, and move issue using VIM keys", args{ev: []*tcell.EventKey{
			tcell.NewEventKey(0, 'l', tcell.ModNone),
			tcell.NewEventKey(0, 'l', tcell.ModNone),
			tcell.NewEventKey(0, 'l', tcell.ModNone),
			tcell.NewEventKey(0, 'h', tcell.ModNone),
			tcell.NewEventKey(0, 'j', tcell.ModNone),
			tcell.NewEventKey(0, 'j', tcell.ModNone),
			tcell.NewEventKey(0, 'k', tcell.ModNone),
		}}, 1, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.InitTestApp(screen)
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				body := `
{
    "transitions": [
        {
            "id": "0",
            "name": "To Do"
        },
        {
            "id": "1",
            "name": "In Progress",
            "to": {
              "id": "2"
            }
        },
        {
            "id": "2",
            "name": "Done"
        }
    ]
}
`
				w.Write([]byte(body)) //nolint:errcheck
			})
			view := NewBoardView(&jira.Project{Id: "1"}, &jira.BoardConfiguration{}, "", api).(*boardView)
			view.columnStatusesMap[0] = []string{"0"}
			view.columnStatusesMap[1] = []string{"1"}
			view.columnStatusesMap[2] = []string{"2"}
			view.columnStatusesMap[3] = []string{"3"}
			view.statusesColumnsMap["0"] = 0
			view.statusesColumnsMap["1"] = 1
			view.statusesColumnsMap["2"] = 2
			view.statusesColumnsMap["3"] = 3
			view.issues = []jira.Issue{
				{Id: "1", Key: "I1", Fields: jira.IssueFields{Status: jira.Status{Id: "1"}}},
				{Id: "2", Key: "I2", Fields: jira.IssueFields{Status: jira.Status{Id: "0"}}},
				{Id: "3", Key: "I3", Fields: jira.IssueFields{Status: jira.Status{Id: "2"}}},
			}
			view.highlightedIssue = &view.issues[0]
			view.columns = []string{"a", "b", "c", "d"}
			view.Refresh()

			// when
			app.GetApp().Loading(false)
			for _, key := range tt.args.ev {
				view.HandleKeyEvent(key)
			}

			// then
			assert.Equal(t, tt.wantCursorX, view.cursorX)
			assert.Equal(t, tt.wantCursorY, view.cursorY)

			// when
			view.moveIssue(view.highlightedIssue, 1)

			// then
			assert.False(t, view.issueSelected)
		})
	}
}

func Test_boardView_Init(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should initialize view and set issues"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.InitTestApp(nil)
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				body := `
{
    "expand": "schema,names",
    "startAt": 0,
    "maxResults": 100,
    "total": 3,
    "issues": [
        {
            "key": "ISSUE-1",
            "fields": {
                "summary": "Issue summary 1",
				"description": "Desc1",
                "status": {
                    "name": "Status1"
                }
            }
        },
        {
            "key": "ISSUE-2",
            "fields": {
                "summary": "Issue summary 2",
				"description": "Desc2",
                "status": {
                    "name": "Status2"
                }
            }
        },
        {
            "key": "ISSUE-3",
            "fields": {
                "summary": "Issue summary 3",
				"description": "Desc3",
                "status": {
                    "name": "Status3"
                }
            }
        }
    ]
}
`
				_, _ = w.Write([]byte(body)) //nolint:errcheck
			})
			view := NewBoardView(&jira.Project{Id: "1"}, &jira.BoardConfiguration{}, "", api).(*boardView)

			// when
			view.Init()

			// then
			assert.Equal(t, 3, len(view.issues))
			assert.Equal(t, "ISSUE-1", view.issues[0].Key)
			assert.Equal(t, "ISSUE-2", view.issues[1].Key)
			assert.Equal(t, "ISSUE-3", view.issues[2].Key)
		})
	}
}
