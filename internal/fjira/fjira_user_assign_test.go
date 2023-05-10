package fjira

import (
	"bytes"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
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
			<-time.NewTimer(300 * time.Millisecond).C
			query := "bob"
			for _, key := range query {
				view.HandleKeyEvent(tcell.NewEventKey(-1, key, tcell.ModNone))
			}
			view.Update()
			view.Update()
			view.Resize(screen.Size())
			<-time.NewTimer(300 * time.Millisecond).C
			view.Update()
			view.Draw(screen)
			<-time.NewTimer(300 * time.Millisecond).C

			var buffer bytes.Buffer
			contents, x, y := screen.GetContents()
			screen.Show()
			for i := 0; i < x*y; i++ {
				buffer.Write(contents[i].Bytes)
			}
			result := buffer.String()

			// then
			assert.Contains(t, result, "Bob")
		})
	}
}

func Test_fjiraAssignChangeView_assignUserToTicket(t *testing.T) {
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
		{"should send assign user request", args{issue: &jira.Issue{Key: "ABC", Id: "123"}}},
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
				if strings.Contains(r.RequestURI, tt.args.issue.Key) {
					_, _ = w.Write([]byte(`{}`))
					assignUserRequestSent <- true
				} else {
					_, _ = w.Write([]byte(`[{"id": "U1", "displayName": "Bob"}, {"id": "U2", "displayName": "John"}]`))
				}
			}))
			go view.startUsersSearching()
			<-time.NewTimer(500 * time.Millisecond).C
			view.HandleKeyEvent(tcell.NewEventKey(-1, 'B', tcell.ModNone))
			view.Update()
			view.Update()
			view.HandleKeyEvent(tcell.NewEventKey(tcell.KeyEnter, -1, tcell.ModNone))
			<-time.NewTimer(250 * time.Millisecond).C
			confirmation := app.GetApp().LastDrawable()
			if kl, ok := (confirmation).(app.KeyListener); ok {
				kl.HandleKeyEvent(tcell.NewEventKey(0, app.Yes, 0))
			}
			<-time.NewTimer(250 * time.Millisecond).C

			// then
			select {
			case <-assignUserRequestSent:
			case <-time.After(5 * time.Second):
				t.Fail()
			}
		})
	}
}

func Test_fjiraAssignChangeView_noUserFound(t *testing.T) {
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
		{"should open issue view again when no user found", args{issue: &jira.Issue{Key: "ABC", Id: "123"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.CreateNewAppWithScreen(screen)
			CreateNewFjira(&fjiraSettings{})
			view := NewAssignChangeView(tt.args.issue)

			// when
			_ = SetApi(jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(`[]`))
			}))
			go view.startUsersSearching()
			<-time.NewTimer(250 * time.Millisecond).C
			view.Update()
			view.Update()
			view.HandleKeyEvent(tcell.NewEventKey(tcell.KeyEnter, -1, tcell.ModNone))
			<-time.NewTimer(1100 * time.Millisecond).C

			// then
			_, ok := app.GetApp().CurrentView().(*fjiraIssueView)
			assert.True(t, ok || app.GetApp().CurrentView() == nil)
		})
	}
}
