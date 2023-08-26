package labels

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

func TestNewAddLabelView(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init()
	defer screen.Fini()

	type args struct {
		issue *jira.Issue
	}
	tests := []struct {
		name string
		args args
	}{
		{"should initialize & draw add label view", args{issue: &jira.Issue{Key: "ABC-123"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.InitTestApp(screen)
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				body := `
{
    "token": "",
    "suggestions": [
        {
            "label": "TestLabel1",
            "html": "<b></b>TestLabel1"
        },
        {
            "label": "TestLabel2",
            "html": "<b></b>TestLabel2"
        }
    ]
}
`
				_, _ = w.Write([]byte(body))
			})
			view := NewAddLabelView(tt.args.issue, func() {}, api).(*addLabelView)

			// when
			view.Init()
			for view.fuzzyFind == nil {
				<-time.After(10 * time.Millisecond)
			}
			query := "label1"
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
		issue *jira.Issue
		label string
	}

	tests := []struct {
		name string
		args args
	}{
		{"should send add label request", args{issue: &jira.Issue{Key: "ABC", Id: "123"}, label: "testlabel1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.InitTestApp(screen)
			addLabelRequestSent := make(chan bool)
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(``))

				assert.Contains(t, r.RequestURI, tt.args.issue.Key)
				addLabelRequestSent <- true
			})
			view := NewAddLabelView(tt.args.issue, func() {}, api).(*addLabelView)

			// when
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
