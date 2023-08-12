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

func TestNewStatusChangeView(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()

	type args struct {
		issue *jira.Issue
	}
	tests := []struct {
		name string
		args args
	}{
		{"should initialize & draw status change view", args{issue: &jira.Issue{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.CreateNewAppWithScreen(screen)
			CreateNewFjira(&fjiraWorkspaceSettings{})
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(`{
    "transitions": [
        {
            "id": "11",
            "name": "To Do"
        },
        {
            "id": "21",
            "name": "In Progress"
        }
    ]
}`)) //nolint:errcheck
			})
			_ = SetApi(api)
			view := NewStatusChangeView(tt.args.issue)

			// when
			view.Init()
			for view.fuzzyFind == nil {
				<-time.After(10 * time.Millisecond)
			}
			query := "in progress"
			for _, key := range query {
				view.HandleKeyEvent(tcell.NewEventKey(-1, key, tcell.ModNone))
			}
			i := 0 // keep app going for a while
			view.Resize(screen.Size())
			for {
				view.Update()
				view.Draw(screen)
				i++
				if i > 100000 {
					break
				}
			}

			// then
			var buffer bytes.Buffer
			contents, x, y := screen.GetContents()
			screen.Show()
			for i := 0; i < x*y; i++ {
				buffer.Write(contents[i].Bytes)
			}
			result := buffer.String()

			assert.Contains(t, result, "In Progress")
			assert.NotContains(t, result, "To Do")
		})
	}
}

func Test_fjiraStatusChangeView_changeStatusForTicket(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()

	type args struct {
		issue  *jira.Issue
		status *jira.IssueTransition
	}

	tests := []struct {
		name string
		args args
	}{
		{"should send change status request", args{issue: &jira.Issue{Key: "ABC", Id: "123"}, status: &jira.IssueTransition{Name: "In Progress", Id: "333"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.CreateNewAppWithScreen(screen)
			CreateNewFjira(&fjiraWorkspaceSettings{})
			view := NewStatusChangeView(tt.args.issue)

			// when
			changeStatusRequestSent := make(chan bool)
			_ = SetApi(jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(`{}`))

				assert.Contains(t, r.RequestURI, tt.args.issue.Key)
				changeStatusRequestSent <- true
			}))
			go func() {
				view.changeStatusTo(tt.args.status)
			}()
			// wait for confirmation
			var confirmation *app.Confirmation
			for confirmation == nil {
				if c, ok := (app.GetApp().LastDrawable()).(*app.Confirmation); ok {
					confirmation = c
				}
				<-time.After(10 * time.Millisecond)
			}
			confirmation.HandleKeyEvent(tcell.NewEventKey(0, app.Yes, 0))

			// then
			select {
			case <-changeStatusRequestSent:
			case <-time.After(3 * time.Second):
				t.Fail()
			}
		})
	}
}
