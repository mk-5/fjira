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
			CreateNewFjira(&fjiraWorkspaceSettings{})
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, err := w.Write([]byte(`[{"id": "U1", "displayName": "Bob"}, {"id": "U2", "displayName": "John"}]`))
				println(err)
			})
			_ = SetApi(api)
			view := NewAssignChangeView(tt.args.issue)

			// when
			view.Init()
			for view.fuzzyFind == nil {
				<-time.After(10 * time.Millisecond)
			}
			query := "bob"
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

			assert.Contains(t, result, "Bob")
		})
	}
}

func Test_fjiraAssignChangeView_assignUserToTicket(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()

	type args struct {
		issue         *jira.Issue
		confirmAction rune
	}

	tests := []struct {
		name string
		args args
	}{
		{"should process assign user request", args{issue: &jira.Issue{Key: "ABC", Id: "123"}, confirmAction: app.Yes}},
		{"should process assign user request", args{issue: &jira.Issue{Key: "ABC", Id: "123"}, confirmAction: app.No}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.CreateNewAppWithScreen(screen)
			CreateNewFjira(&fjiraWorkspaceSettings{})
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
			go func() {
				view.startUsersSearching()
			}()
			for view.fuzzyFind == nil {
				<-time.After(10 * time.Millisecond)
			}
			view.fuzzyFind.HandleKeyEvent(tcell.NewEventKey(-1, 'B', tcell.ModNone))
			view.fuzzyFind.HandleKeyEvent(tcell.NewEventKey(tcell.KeyEnter, -1, tcell.ModNone))
			// wait for confirmation
			var confirmation *app.Confirmation
			for confirmation == nil {
				if c, ok := (app.GetApp().LastDrawable()).(*app.Confirmation); ok {
					confirmation = c
				}
				<-time.After(10 * time.Millisecond)
			}
			confirmation.HandleKeyEvent(tcell.NewEventKey(0, tt.args.confirmAction, 0))
			confirmation.Update()

			// then
			select {
			case <-assignUserRequestSent:
			case <-time.After(5 * time.Second):
				_, ok := app.GetApp().CurrentView().(*fjiraIssueView)
				if !ok {
					t.Fail()
				}
			}
		})
	}
}

func Test_fjiraAssignChangeView_assignUserToTicket_empty_user(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()

	tests := []struct {
		name string
	}{
		{"should send set issue view instead of assign user request when user is empty"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.CreateNewAppWithScreen(screen)
			CreateNewFjira(&fjiraWorkspaceSettings{})
			view := NewAssignChangeView(&jira.Issue{})

			// when
			view.assignUserToTicket(&jira.Issue{}, nil)

			// then
			_, ok := app.GetApp().CurrentView().(*fjiraIssueView)
			assert.True(t, ok)
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
			CreateNewFjira(&fjiraWorkspaceSettings{})
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
