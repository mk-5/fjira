package fjira

import (
	"errors"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/boards"
	"github.com/mk-5/fjira/internal/issues"
	"github.com/mk-5/fjira/internal/jira"
	"github.com/mk-5/fjira/internal/labels"
	"github.com/mk-5/fjira/internal/projects"
	"github.com/mk-5/fjira/internal/statuses"
	"github.com/mk-5/fjira/internal/ui"
	"github.com/mk-5/fjira/internal/users"
	"github.com/mk-5/fjira/internal/workspaces"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	WelcomeMessage = `
    ____    __________  ___ 
   / __/   / /  _/ __ \/   |
  / /___  / // // /_/ / /| |
 / __/ /_/ // // _, _/ ___ |
/_/  \____/___/_/ |_/_/  |_|
                            
The command line tool for Jira.
`
)

var InstallFailedErr = errors.New("cannot use fjira. Please check error logs in order to install missing packages")
var NotInitializedErr = errors.New("cannot use fjira. You need to call CreateNewFjira first")

type Fjira struct {
	app       *app.App
	api       jira.Api
	jiraUrl   string
	workspace string
}

type CliArgs struct {
	ProjectId       string
	IssueKey        string
	Workspace       string
	WorkspaceSwitch bool
	WorkspaceEdit   bool
	JqlMode         bool
}

var (
	fjiraInstance *Fjira
	fjiraOnce     sync.Once
)

func CreateNewFjira(settings *workspaces.WorkspaceSettings) *Fjira {
	if settings == nil {
		panic("Cannot find appropriate fjira settings!")
	}
	fjiraOnce.Do(func() {
		url := strings.TrimSuffix(settings.JiraRestUrl, "/")
		api, err := jira.NewApi(url, settings.JiraUsername, settings.JiraToken, settings.JiraTokenType)
		if err != nil {
			app.Error(err.Error())
		}
		fjiraInstance = &Fjira{
			app:       app.CreateNewApp(),
			api:       api,
			jiraUrl:   url,
			workspace: settings.Workspace,
		}
	})
	return fjiraInstance
}

func SetApi(api jira.Api) error {
	if fjiraInstance == nil {
		return NotInitializedErr
	}
	fjiraInstance.api = api
	return nil
}

func Install(args CliArgs) (*workspaces.WorkspaceSettings, error) {
	// it will be removed after a few version
	u := workspaces.NewDeprecatedUserHomeWorkspaces()
	_ = u.MigrateFromGlobWorkspacesToYaml()

	err := validateWorkspaceName(args.Workspace)
	if err != nil {
		return nil, err
	}
	if args.WorkspaceEdit {
		s, err := readFromWorkspaceEdit(os.Stdin, args.Workspace)
		if err != nil {
			panic(err)
		}
		return s, nil
	}
	s, err := readFromEnvironments()
	if err == nil {
		return s, nil // envs found
	}
	if err != EnvironmentsMissingErr {
		return nil, err
	}
	s2, err := readFromUserSettings(args.Workspace)
	if err == workspaces.WorkspaceNotFoundErr || errors.Unwrap(err) == workspaces.WorkspaceNotFoundErr {
		return readFromUserInputAndStore(os.Stdin, args.Workspace, nil)
	}
	if err != nil {
		return nil, err
	}
	return s2, nil
}

func (f *Fjira) SetApi(api jira.Api) {
	f.api = api
}

func (f *Fjira) GetApi() jira.Api {
	return f.api
}

func (f *Fjira) Run(args *CliArgs) {
	x := app.ClampInt(f.app.ScreenX/2-18, 0, f.app.ScreenX)
	y := app.ClampInt(f.app.ScreenY/2-4, 0, f.app.ScreenY)
	welcomeText := app.NewText(x, y, app.DefaultStyle, WelcomeMessage)
	f.app.AddDrawable(welcomeText)
	f.registerGoTos()
	go f.bootstrap(args)
	f.app.Start()
}

func (f *Fjira) Close() {
	f.api.Close()
	if f.app != nil {
		f.app.PanicRecover()
	}
}

func (f *Fjira) registerGoTos() {
	projects.RegisterGoto()
	issues.RegisterGoTo()
	users.RegisterGoTo()
	statuses.RegisterGoTo()
	labels.RegisterGoTo()
	workspaces.RegisterGoTo()
	boards.RegisterGoTo()
	ui.RegisterGoTo()
}

func (f *Fjira) bootstrap(args *CliArgs) {
	defer f.app.PanicRecover()
	if args.WorkspaceSwitch {
		app.GoTo("workspaces-switch")
		return
	}
	if args.ProjectId != "" {
		app.GoTo("issues-search", args.ProjectId, func() {
			app.GoTo("projects", f.api)
		}, f.api)
		return
	}
	if args.IssueKey != "" {
		app.GoTo("issue", args.IssueKey, nil, f.api)
		return
	}
	if args.JqlMode {
		app.GoTo("jql", f.api)
		return
	}
	time.Sleep(350 * time.Millisecond)
	f.app.RunOnAppRoutine(func() {
		app.GoTo("projects", f.api)
	})
}
