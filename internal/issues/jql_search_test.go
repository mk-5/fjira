package issues

import (
	"bytes"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	os2 "github.com/mk-5/fjira/internal/os"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
	"time"
)

func prepareTestScreen(t *testing.T) tcell.SimulationScreen {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	app.InitTestApp(screen)
	tempDir := t.TempDir()
	_ = os2.SetUserHomeDir(tempDir)
	return screen
}

func TestNewJqlSearchView(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should render issues returned by jql query"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			screen := prepareTestScreen(t)
			defer screen.Fini()
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
			view := NewJqlSearchView(api).(*jqlSearchView)

			// when
			view.Init()
			<-time.After(time.Millisecond * 200)
			view.HandleKeyEvent(tcell.NewEventKey(tcell.KeyACK, 't', tcell.ModNone))
			// keep going app for a while
			i := 0
			for {
				view.Update()
				view.Draw(screen)
				i++
				if i > 1000000 {
					break
				}
			}
			var buffer bytes.Buffer
			contents, x, y := screen.GetContents()
			screen.Show()
			for i := 0; i < x*y; i++ {
				if len(contents[i].Bytes) != 0 {
					buffer.Write(contents[i].Bytes)
				}
			}
			result := strings.TrimSpace(buffer.String())

			// then
			assert.Contains(t, result, "ISSUE-1")
			assert.Contains(t, result, "ISSUE-2")
			assert.Contains(t, result, "ISSUE-3")

			// and then
			view.Destroy()
		})
	}
}
