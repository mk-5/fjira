package fjira

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestFjira_bootstrap(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()

	type args struct {
		cliArgs       CliArgs
		viewPredicate func() bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"should switch to workspace view", args{
			cliArgs: CliArgs{WorkspaceSwitch: true},
			viewPredicate: func() bool {
				_, ok := app.GetApp().CurrentView().(*fjiraSwitchWorkspaceView)
				return ok
			},
		}},
		{"should switch to project view", args{
			cliArgs: CliArgs{ProjectId: "test"},
			viewPredicate: func() bool {
				_, ok := app.GetApp().CurrentView().(*fjiraSearchIssuesView)
				return ok
			},
		}},
		{"should switch to issue view", args{
			cliArgs: CliArgs{IssueKey: "test"},
			viewPredicate: func() bool {
				_, ok := app.GetApp().CurrentView().(*fjiraIssueView)
				return ok
			},
		}},
		{"should switch to jql view", args{
			cliArgs: CliArgs{JqlMode: true},
			viewPredicate: func() bool {
				_, ok := app.GetApp().CurrentView().(*fjiraJqlSearchView)
				return ok
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Write([]byte("{}")) //nolint:errcheck
			})
			app.CreateNewAppWithScreen(screen)
			fjira := CreateNewFjira(&fjiraSettings{})
			_ = SetApi(api)

			// when
			fjira.bootstrap(&tt.args.cliArgs)

			// then
			ok := tt.args.viewPredicate()
			assert.New(t).True(ok, "Current view is invalid.")
		})
	}
}
