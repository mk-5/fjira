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

func TestNewStatusChangeView(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	screen.Init() //nolint:errcheck
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
			CreateNewFjira(&fjiraSettings{})
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
			<-time.After(1 * time.Second)
			query := "in progress"
			for _, key := range query {
				view.HandleKeyEvent(tcell.NewEventKey(-1, key, tcell.ModNone))
			}
			view.Update()
			view.Resize(screen.Size())
			<-time.After(1 * time.Second)
			view.Update()
			view.Draw(screen)
			<-time.After(1 * time.Second)

			var buffer bytes.Buffer
			contents, x, y := screen.GetContents()
			screen.Show()
			for i := 0; i < x*y; i++ {
				buffer.Write(contents[i].Bytes)
			}
			result := buffer.String()

			// then
			assert.Contains(t, result, "In Progress")
			assert.NotContains(t, result, "To Do")
		})
	}
}

func Test_fjiraStatusChangeView_changeStatusForTicket(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	screen.Init() //nolint:errcheck
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
			CreateNewFjira(&fjiraSettings{})
			view := NewStatusChangeView(tt.args.issue)

			// when
			changeStatusRequestSent := make(chan bool)
			_ = SetApi(jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(`{}`))

				assert.Contains(t, r.RequestURI, tt.args.issue.Key)
				changeStatusRequestSent <- true
			}))
			go view.changeStatusForTicket(tt.args.issue, tt.args.status)

			// then
			select {
			case <-changeStatusRequestSent:
			case <-time.After(3 * time.Second):
				t.Fail()
			}
		})
	}
}
