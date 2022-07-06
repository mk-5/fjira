package fjira

import (
	"errors"
	"github.com/gdamore/tcell/v2"
	"github.com/mk5/fjira/internal/app"
	"github.com/mk5/fjira/internal/jira"
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

var InstallFailedErr = errors.New("Cannot use fjira. Please check error logs in order to install missing packages.")
var FjiraNotInitalizedErr = errors.New("Cannot use fjira. You need to call CreateNewFjira first.")

type Fjira struct {
	app       *app.App
	api       jira.JiraApi
	formatter FjiraFormatter
	jiraUrl   string
}

type CliArgs struct {
	ProjectId       string
	IssueKey        string
	Workspace       string
	SwitchWorkspace bool
}

var (
	fjiraInstance *Fjira
	fjiraOnce     sync.Once
)

func CreateNewFjira(settings *fjiraSettings) *Fjira {
	if settings == nil {
		panic("Cannot find appropriate fjira settings!")
	}
	fjiraOnce.Do(func() {
		url := strings.TrimSuffix(settings.JiraRestUrl, "/")
		api, err := jira.NewJiraApi(url, settings.JiraUsername, settings.JiraToken)
		if err != nil {
			app.Error(err.Error())
		}
		fjiraInstance = &Fjira{
			app:       app.CreateNewApp(),
			api:       api,
			formatter: &defaultFormatter{},
			jiraUrl:   url,
		}
	})
	return fjiraInstance
}

func GetApi() (jira.JiraApi, error) {
	if fjiraInstance == nil {
		return nil, FjiraNotInitalizedErr
	}
	return fjiraInstance.api, nil
}

func SetApi(api jira.JiraApi) error {
	if fjiraInstance == nil {
		return FjiraNotInitalizedErr
	}
	fjiraInstance.api = api
	return nil
}

func GetFormatter() (FjiraFormatter, error) {
	if fjiraInstance == nil {
		return nil, FjiraNotInitalizedErr
	}
	return fjiraInstance.formatter, nil
}

func GetJiraUrl() (string, error) {
	if fjiraInstance == nil {
		return "", FjiraNotInitalizedErr
	}
	return fjiraInstance.jiraUrl, nil
}

func Install(workspace string) (*fjiraSettings, error) {
	err := validateWorkspaceName(workspace)
	if err != nil {
		return nil, err
	}
	settings, err := readFromEnvironments()
	if err == nil {
		return settings, nil // envs found
	}
	if err != EnvironmentsMissingErr {
		return nil, err
	}
	settings2, err := readFromUserSettings(workspace)
	if err == WorkspaceNotFoundErr {
		return readFromUserInputAndStore(workspace)
	}
	if err != nil {
		return nil, err
	}
	return settings2, nil
}

func (f *Fjira) SetApi(api jira.JiraApi) {
	f.api = api
}

func (f *Fjira) Run(args *CliArgs) {
	x := app.ClampInt(f.app.ScreenX/2-18, 0, f.app.ScreenX)
	y := app.ClampInt(f.app.ScreenY/2-4, 0, f.app.ScreenY)
	welcomeText := app.NewText(x, y, tcell.StyleDefault, WelcomeMessage)
	f.app.AddDrawable(welcomeText)
	go f.bootstrap(args)
	f.app.Start()
}

func (f *Fjira) Close() {
	f.api.Close()
	if f.app != nil {
		f.app.PanicRecover()
	}
}

func (f *Fjira) bootstrap(args *CliArgs) {
	defer f.app.PanicRecover()
	if args.ProjectId != "" {
		goIntoIssuesSearchForProject(args.ProjectId)
		return
	}
	if args.IssueKey != "" {
		goIntoIssueView(args.IssueKey)
		return
	}
	if args.SwitchWorkspace {
		goIntoSwitchWorkspaceView()
		return
	}
	time.Sleep(500 * time.Millisecond)
	f.app.RunOnAppRoutine(func() {
		goIntoProjectsSearch()
	})
}
