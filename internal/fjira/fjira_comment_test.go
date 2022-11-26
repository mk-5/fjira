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

func TestNewCommentView(t *testing.T) {
	type args struct {
		issue *jira.Issue
	}
	tests := []struct {
		name string
		args args
	}{
		{"should create new comment view", args{issue: &jira.Issue{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, newCommentView(tt.args.issue), "newCommentView(%v)", tt.args.issue)
		})
	}
}

func Test_fjiraCommentView_Destroy(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should run Destroy without problem"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newCommentView(&jira.Issue{})
			view.Destroy()
		})
	}
}

func Test_fjiraCommentView_Draw(t *testing.T) {
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
		{"should draw comment view", args{screen: screen}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newCommentView(&jira.Issue{})
			view.text = "Comment text"

			// when
			view.Draw(tt.args.screen)
			var buffer bytes.Buffer
			contents, x, y := tt.args.screen.(tcell.SimulationScreen).GetContents()
			tt.args.screen.Show()
			for i := 0; i < x*y; i++ {
				if len(contents[i].Bytes) != 0 {
					buffer.Write(contents[i].Bytes)
				}
			}
			result := buffer.String()

			// then
			assert.Contains(t, result, view.text)
		})
	}
}

func Test_fjiraCommentView_HandleKeyEvent(t *testing.T) {
	type args struct {
		ev []*tcell.EventKey
	}
	tests := []struct {
		name            string
		args            args
		expectedComment string
	}{
		{"should handle key events and write comment", args{ev: []*tcell.EventKey{
			tcell.NewEventKey(0, 'a', tcell.ModNone),
			tcell.NewEventKey(0, 'b', tcell.ModNone),
			tcell.NewEventKey(0, 'c', tcell.ModNone),
		}}, "abc"},
		{"should handle key events with backspace", args{ev: []*tcell.EventKey{
			tcell.NewEventKey(0, 'a', tcell.ModNone),
			tcell.NewEventKey(0, 'b', tcell.ModNone),
			tcell.NewEventKey(0, 'c', tcell.ModNone),
			tcell.NewEventKey(tcell.KeyBackspace, '-', tcell.ModNone),
		}}, "ab"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newCommentView(&jira.Issue{})

			// when
			for _, key := range tt.args.ev {
				view.HandleKeyEvent(key)
			}

			// then
			assert.Equal(t, tt.expectedComment, view.text)
		})
	}
}

func Test_fjiraCommentView_Init(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()

	tests := []struct {
		name string
	}{
		{"should initialize doComment handling"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.CreateNewAppWithScreen(screen)
			CreateNewFjira(&fjiraSettings{})
			done := make(chan bool)
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(``)) //nolint:errcheck
				done <- true
			})
			_ = SetApi(api) //nolint:errcheck
			view := newCommentView(&jira.Issue{})

			// when
			view.Init()
			view.HandleKeyEvent(tcell.NewEventKey(tcell.KeyF1, 'F', tcell.ModNone))

			// then
			select {
			case <-done:
			case <-time.After(3 * time.Second):
				t.Fail()
			}
		})
	}
}

func Test_fjiraCommentView_Resize(t *testing.T) {
	type args struct {
		screenX int
		screenY int
	}
	tests := []struct {
		name string
		args args
	}{
		{"should resize without problems", args{screenY: 10, screenX: 10}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newCommentView(&jira.Issue{})
			view.Resize(tt.args.screenX, tt.args.screenY)
		})
	}
}

func Test_fjiraCommentView_Update(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should update without problems"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newCommentView(&jira.Issue{})
			view.Update()
		})
	}
}
