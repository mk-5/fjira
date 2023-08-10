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

func TestNewProjectsSearchView(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()
	tests := []struct {
		name string
	}{
		{"should initialize & draw projects search view"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.CreateNewAppWithScreen(screen)
			CreateNewFjira(&fjiraWorkspaceSettings{})
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(`[{"id": "1", "name": "Test", "key": "TEST"}, {"id": "2", "name": "Fjira", "key":"FJIR"}]`))
			})
			_ = SetApi(api)
			view := NewProjectsSearchView()

			// when
			view.Init()
			for view.fuzzyFind == nil {
				<-time.After(10 * time.Millisecond)
			}
			query := "FJIR"
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

			assert.Contains(t, result, "Fjira")
			assert.NotContains(t, result, "TEST")
		})
	}
}
