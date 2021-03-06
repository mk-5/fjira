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

func TestNewProjectsSearchView(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	screen.Init() //nolint:errcheck
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
			CreateNewFjira(&fjiraSettings{})
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(`[{"id": "1", "name": "Test", "key": "TEST"}, {"id": "2", "name": "Fjira", "key":"FJIR"}]`))
			})
			_ = SetApi(api)
			view := NewProjectsSearchView()

			// when
			view.Init()
			<-time.NewTimer(100 * time.Millisecond).C
			query := "FJIR"
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
			contents, x, y := screen.(tcell.SimulationScreen).GetContents()
			screen.Show()
			for i := 0; i < x*y; i++ {
				buffer.Write(contents[i].Bytes)
			}
			result := buffer.String()

			// then
			assert.Contains(t, result, "Fjira")
			assert.NotContains(t, result, "TEST")
		})
	}
}
