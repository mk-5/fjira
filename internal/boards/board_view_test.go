package boards

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	"github.com/stretchr/testify/assert"
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
		{"should create a new board view", args{boardConfiguration: &jira.BoardConfiguration{Id: 1}, project: &jira.Project{}}},
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
			b := NewBoardView(&jira.Project{}, &jira.BoardConfiguration{Id: 1}, "", nil)
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
			view := NewBoardView(&jira.Project{Id: "1"}, &jira.BoardConfiguration{Id: 1}, "", api).(*boardView)
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
			view := NewBoardView(&jira.Project{Id: "1"}, &jira.BoardConfiguration{Id: 1}, "", api).(*boardView)

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

func Test_boardView_SetSprints_SelectsActiveSprint(t *testing.T) {
	// given
	view := NewBoardView(&jira.Project{}, &jira.BoardConfiguration{}, "", nil).(*boardView)
	sprints := []jira.SprintItem{
		{Id: 1, Name: "Sprint Future", State: "future"},
		{Id: 2, Name: "Sprint Active", State: "active"},
	}

	// when
	view.SetSprints(sprints)

	// then
	assert.NotNil(t, view.activeSprint)
	assert.Equal(t, 2, view.activeSprint.Id)
	assert.Equal(t, "active", view.activeSprint.State)
	assert.Equal(t, "Sprint Active", view.activeSprint.Name)
	assert.Equal(t, 2, len(view.sprints))
}

func Test_boardView_fetchIssues_WithoutSprint_UsesSearchJql(t *testing.T) {
	app.InitTestApp(nil)
	api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{
  "startAt": 0,
  "maxResults": 100,
  "total": 2,
  "issues": [
    { "id": "10001", "key": "K1" },
    { "id": "10002", "key": "K2" }
  ]
}`))
	})

	view := NewBoardView(&jira.Project{}, &jira.BoardConfiguration{}, "project = GEN", api).(*boardView)

	// when
	issues, err := view.fetchIssues()

	// then
	assert.NoError(t, err)
	assert.Len(t, issues, 2)
	assert.Equal(t, "K1", issues[0].Key)
	assert.Equal(t, "K2", issues[1].Key)
}

func Test_boardView_fetchIssues_WithActiveSprint_UsesSprintIssues(t *testing.T) {
	app.InitTestApp(nil)
	api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{
  "total": 2,
  "maxResults": 100,
  "issues": [
    { "id": "20001", "key": "S1" },
    { "id": "20002", "key": "S2" }
  ]
}`))
	})

	view := NewBoardView(&jira.Project{}, &jira.BoardConfiguration{Id: 1}, "project = GEN", api).(*boardView)
	view.activeSprint = &jira.SprintItem{Id: 7, Name: "Sprint 7", State: "active"}

	// when
	issues, err := view.fetchIssues()

	// then
	assert.NoError(t, err)
	assert.Len(t, issues, 2)
	assert.Equal(t, "S1", issues[0].Key)
	assert.Equal(t, "S2", issues[1].Key)
}

func Test_boardView_assigneeFiltering(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()

	type args struct {
		assigneeFilter string
	}
	tests := []struct {
		name           string
		args           args
		expectedIssues int
	}{
		{"should filter by specific assignee", args{assigneeFilter: "John Doe"}, 2},
		{"should show all issues when 'All' is selected", args{assigneeFilter: "All"}, 4},
		{"should show unassigned issues", args{assigneeFilter: "Unassigned"}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.InitTestApp(screen)
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				body := `{
    "expand": "schema,names",
    "startAt": 0,
    "maxResults": 100,
    "total": 4,
    "issues": [
        {
            "id": "1",
            "key": "TEST-1",
            "fields": {
                "summary": "First issue",
                "assignee": {
                    "accountId": "john-doe",
                    "displayName": "John Doe"
                },
                "status": {
                    "id": "10000"
                }
            }
        },
        {
            "id": "2", 
            "key": "TEST-2",
            "fields": {
                "summary": "Second issue",
                "assignee": {
                    "accountId": "john-doe",
                    "displayName": "John Doe"
                },
                "status": {
                    "id": "10001"
                }
            }
        },
        {
            "id": "3",
            "key": "TEST-3", 
            "fields": {
                "summary": "Third issue",
                "assignee": {
                    "accountId": "jane-smith",
                    "displayName": "Jane Smith"
                },
                "status": {
                    "id": "10002"
                }
            }
        },
        {
            "id": "4",
            "key": "TEST-4",
            "fields": {
                "summary": "Fourth issue",
                "assignee": {
                    "accountId": "",
                    "displayName": ""
                },
                "status": {
                    "id": "10003"
                }
            }
        }
    ]
}`
				w.Write([]byte(body)) //nolint:errcheck
			})

			boardConfig := &jira.BoardConfiguration{
				Id: 1,
				ColumnConfig: struct {
					Columns []struct {
						Name     string `json:"name"`
						Statuses []struct {
							Id   string `json:"id"`
							Self string `json:"self"`
						} `json:"statuses"`
					} `json:"columns"`
					ConstraintType string `json:"constraintType"`
				}{
					Columns: []struct {
						Name     string `json:"name"`
						Statuses []struct {
							Id   string `json:"id"`
							Self string `json:"self"`
						} `json:"statuses"`
					}{
						{Name: "Col1", Statuses: []struct {
							Id   string `json:"id"`
							Self string `json:"self"`
						}{{Id: "10000"}}},
						{Name: "Col2", Statuses: []struct {
							Id   string `json:"id"`
							Self string `json:"self"`
						}{{Id: "10001"}}},
						{Name: "Col3", Statuses: []struct {
							Id   string `json:"id"`
							Self string `json:"self"`
						}{{Id: "10002"}}},
						{Name: "Col4", Statuses: []struct {
							Id   string `json:"id"`
							Self string `json:"self"`
						}{{Id: "10003"}}},
					},
				},
			}

			view := NewBoardView(&jira.Project{}, boardConfig, "", api).(*boardView)
			view.Init()

			// Simulate applying assignee filter
			if tt.args.assigneeFilter == "John Doe" {
				user := &jira.User{
					AccountId:   "john-doe",
					DisplayName: "John Doe",
				}
				view.applyAssigneeFilter(user)
			} else if tt.args.assigneeFilter == "All" {
				view.clearAssigneeFilter()
			} else if tt.args.assigneeFilter == "Unassigned" {
				user := &jira.User{
					AccountId:   "",
					DisplayName: "Unassigned",
				}
				view.applyAssigneeFilter(user)
			}

			// Verify the number of issues after filtering
			assert.Equal(t, tt.expectedIssues, len(view.issues))
		})
	}
}

func Test_boardView_getBoardAssignees(t *testing.T) {
	tests := []struct {
		name              string
		expectedAssignees int
		hasUnassigned     bool
	}{
		{"should extract unique assignees including unassigned", 2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := &boardView{
				allIssues: []jira.Issue{
					{
						Fields: jira.IssueFields{
							Assignee: struct {
								AccountId   string `json:"accountId"`
								DisplayName string `json:"displayName"`
							}{
								AccountId:   "john-doe",
								DisplayName: "John Doe",
							},
						},
					},
					{
						Fields: jira.IssueFields{
							Assignee: struct {
								AccountId   string `json:"accountId"`
								DisplayName string `json:"displayName"`
							}{
								AccountId:   "jane-smith",
								DisplayName: "Jane Smith",
							},
						},
					},
					{
						Fields: jira.IssueFields{
							Assignee: struct {
								AccountId   string `json:"accountId"`
								DisplayName string `json:"displayName"`
							}{
								AccountId:   "",
								DisplayName: "",
							},
						},
					},
				},
			}

			assignees := view.getBoardAssignees()

			// Should have unique assignees plus unassigned
			assert.Equal(t, tt.expectedAssignees+1, len(assignees)) // +1 for Unassigned

			// Check if unassigned option is included
			foundUnassigned := false
			for _, assignee := range assignees {
				if assignee.DisplayName == "Unassigned" {
					foundUnassigned = true
					break
				}
			}
			assert.Equal(t, tt.hasUnassigned, foundUnassigned)
		})
	}
}

func Test_boardView_actionBarAssigneeFilter(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()

	tests := []struct {
		name string
	}{
		{"should create board view with assignee filter functionality"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.InitTestApp(screen)

			boardConfig := &jira.BoardConfiguration{
				ColumnConfig: struct {
					Columns []struct {
						Name     string `json:"name"`
						Statuses []struct {
							Id   string `json:"id"`
							Self string `json:"self"`
						} `json:"statuses"`
					} `json:"columns"`
					ConstraintType string `json:"constraintType"`
				}{
					Columns: []struct {
						Name     string `json:"name"`
						Statuses []struct {
							Id   string `json:"id"`
							Self string `json:"self"`
						} `json:"statuses"`
					}{
						{Name: "Col1", Statuses: []struct {
							Id   string `json:"id"`
							Self string `json:"self"`
						}{{Id: "10000"}}},
					},
				},
			}

			view := NewBoardView(&jira.Project{}, boardConfig, "", nil).(*boardView)

			// The board view should be created successfully with assignee filter functionality
			assert.NotNil(t, view)
			assert.NotNil(t, view.bottomBar)
			assert.Nil(t, view.assigneeFilter) // should be nil initially
		})
	}
}
