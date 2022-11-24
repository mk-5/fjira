package fjira

import (
	"errors"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
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
	api       jira.Api
	formatter FjiraFormatter
	jiraUrl   string
}

type CliArgs struct {
	ProjectId       string
	IssueKey        string
	Workspace       string
	WorkspaceSwitch bool
	WorkspaceEdit   bool
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
		api, err := jira.NewApi(url, settings.JiraUsername, settings.JiraToken)
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

func GetApi() (jira.Api, error) {
	if fjiraInstance == nil {
		return nil, FjiraNotInitalizedErr
	}
	return fjiraInstance.api, nil
}

func SetApi(api jira.Api) error {
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

func Install(args CliArgs) (*fjiraSettings, error) {
	err := validateWorkspaceName(args.Workspace)
	if err != nil {
		return nil, err
	}
	if args.WorkspaceEdit {
		settings, err := readFromWorkspaceEdit(args.Workspace)
		if err != nil {
			panic(err)
		}
		return settings, nil
	}
	settings, err := readFromEnvironments()
	if err == nil {
		return settings, nil // envs found
	}
	if err != EnvironmentsMissingErr {
		return nil, err
	}
	settings2, err := readFromUserSettings(args.Workspace)
	if err == WorkspaceNotFoundErr {
		return readFromUserInputAndStore(args.Workspace, nil)
	}
	if err != nil {
		return nil, err
	}
	return settings2, nil
}

func (f *Fjira) SetApi(api jira.Api) {
	f.api = api
}

func (f *Fjira) Run(args *CliArgs) {
	x := app.ClampInt(f.app.ScreenX/2-18, 0, f.app.ScreenX)
	y := app.ClampInt(f.app.ScreenY/2-4, 0, f.app.ScreenY)
	welcomeText := app.NewText(x, y, app.DefaultStyle, WelcomeMessage)
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
	if args.WorkspaceSwitch {
		goIntoSwitchWorkspaceView()
		return
	}
	if args.ProjectId != "" {
		goIntoIssuesSearchForProject(args.ProjectId)
		return
	}
	if args.IssueKey != "" {
		goIntoIssueView(args.IssueKey)
		return
	}
	time.Sleep(350 * time.Millisecond)
	f.app.RunOnAppRoutine(func() {
		goIntoProjectsSearch()
	})
}
