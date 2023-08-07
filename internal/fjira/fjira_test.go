package fjira

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	os2 "github.com/mk-5/fjira/internal/os"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
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
		{"should switch to projects search by default", args{
			cliArgs: CliArgs{},
			viewPredicate: func() bool {
				<-time.After(500 * time.Millisecond)
				_, ok := app.GetApp().CurrentView().(*fjiraSearchProjectsView)
				return ok
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			tempDir := t.TempDir()
			_ = os2.SetUserHomeDir(tempDir)
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Write([]byte("{}")) //nolint:errcheck
			})
			a := app.CreateNewAppWithScreen(screen)
			fjira := CreateNewFjira(&fjiraWorkspaceSettings{})
			fjira.app = a
			_ = SetApi(api)
			go a.Start()

			// when
			go fjira.bootstrap(&tt.args.cliArgs)
			for a.CurrentView() == nil {
				<-time.After(10 * time.Millisecond)
			}

			// then
			ok := tt.args.viewPredicate()
			assert.New(t).True(ok, "Current view is invalid: ", app.GetApp().CurrentView())
		})
	}
}

func TestFjira_run_should_run_without_error(t *testing.T) {
	// given
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()
	api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("{}")) //nolint:errcheck
	})
	app.CreateNewAppWithScreen(screen)
	fjira := CreateNewFjira(&fjiraWorkspaceSettings{})
	_ = SetApi(api)

	// when
	go fjira.Run(&CliArgs{})
	<-time.After(300 * time.Millisecond)

	// then
	assert.False(t, app.GetApp().IsQuit())
}
