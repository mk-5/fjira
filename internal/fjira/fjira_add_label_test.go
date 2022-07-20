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

func TestNewAddLabelView(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init()
	defer screen.Fini()

	type args struct {
		issue *jira.JiraIssue
	}
	tests := []struct {
		name string
		args args
	}{
		{"should initialize & draw add label view", args{issue: &jira.JiraIssue{Key: "ABC-123"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.CreateNewAppWithScreen(screen)
			CreateNewFjira(&fjiraSettings{})
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				body := `
{
    "maxResults": 1000,
    "startAt": 0,
    "total": 3,
    "isLast": true,
    "values": [
        "TestLabel1", "TestLabel2"
    ]
}
`
				_, _ = w.Write([]byte(body))
			})
			_ = SetApi(api)
			view := NewAddLabelView(tt.args.issue)

			// when
			view.Init()
			<-time.NewTimer(100 * time.Millisecond).C
			query := "label1"
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
			assert.Contains(t, result, "TestLabel1")
			assert.NotContains(t, result, "TestLabel2")
		})
	}
}

func Test_fjiraAddLabelView_doAddLabel(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()

	type args struct {
		issue *jira.JiraIssue
		label string
	}

	tests := []struct {
		name string
		args args
	}{
		{"should send add label request", args{issue: &jira.JiraIssue{Key: "ABC", Id: "123"}, label: "testlabel1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.CreateNewAppWithScreen(screen)
			CreateNewFjira(&fjiraSettings{})
			view := NewAddLabelView(tt.args.issue)

			// when
			addLabelRequestSent := make(chan bool)
			_ = SetApi(jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(``))

				assert.Contains(t, r.RequestURI, tt.args.issue.Key)
				addLabelRequestSent <- true
			}))
			go view.doAddLabel(tt.args.issue, tt.args.label)

			// then
			select {
			case <-addLabelRequestSent:
			case <-time.After(3 * time.Second):
				t.Fail()
			}
		})
	}
}
