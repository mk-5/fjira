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

func TestNewAssignChangeView(t *testing.T) {
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
		{"should initialize & draw assign user view", args{issue: &jira.Issue{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.CreateNewAppWithScreen(screen)
			CreateNewFjira(&fjiraSettings{})
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, err := w.Write([]byte(`[{"id": "U1", "displayName": "Bob"}, {"id": "U2", "displayName": "John"}]`))
				println(err)
			})
			_ = SetApi(api)
			view := NewAssignChangeView(tt.args.issue)

			// when
			view.Init()
			<-time.NewTimer(100 * time.Millisecond).C
			query := "bob"
			for _, key := range query {
				view.HandleKeyEvent(tcell.NewEventKey(-1, key, tcell.ModNone))
			}
			view.Update()
			view.Resize(screen.Size())
			<-time.NewTimer(100 * time.Millisecond).C
			view.Update()
			view.Draw(screen)
			<-time.NewTimer(100 * time.Millisecond).C

			var buffer bytes.Buffer
			contents, x, y := screen.GetContents()
			screen.Show()
			for i := 0; i < x*y; i++ {
				buffer.Write(contents[i].Bytes)
			}
			result := buffer.String()

			// then
			assert.Contains(t, result, "Bob")
			assert.NotContains(t, result, "John")
		})
	}
}

func Test_fjiraAssignChangeView_doAssignmentChange(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()

	type args struct {
		issue *jira.Issue
		user  *jira.User
	}

	tests := []struct {
		name string
		args args
	}{
		{"should send assign user request", args{issue: &jira.Issue{Key: "ABC", Id: "123"}, user: &jira.User{DisplayName: "Bob", AccountId: "333"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.CreateNewAppWithScreen(screen)
			CreateNewFjira(&fjiraSettings{})
			view := NewAssignChangeView(tt.args.issue)

			// when
			assignUserRequestSent := make(chan bool)
			_ = SetApi(jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(`{}`))

				assert.Contains(t, r.RequestURI, tt.args.issue.Key)
				assignUserRequestSent <- true
			}))
			go view.doAssignmentChange(tt.args.issue, tt.args.user)

			// then
			select {
			case <-assignUserRequestSent:
			case <-time.After(3 * time.Second):
				t.Fail()
			}
		})
	}
}
