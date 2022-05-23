package fjira

import (
	"errors"
	"github.com/gdamore/tcell"
	"github.com/mk5/fjira/internal/app"
	"github.com/mk5/fjira/internal/jira"
	"os"
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
	ProjectId string
	IssueKey  string
}

var (
	fjiraInstance *Fjira
	fjiraOnce     sync.Once
)

func CreateNewFjira(api jira.JiraApi) *Fjira {
	fjiraOnce.Do(func() {
		url := os.Getenv(JiraRestUrlEnv)
		if api == nil {
			a, err := jira.NewJiraApi(url, os.Getenv(JiraUsernameEnv), os.Getenv(JiraTokenEnv))
			api = a
			if err != nil {
				app.Error(err.Error())
			}
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

func (f *Fjira) Install() []error {
	errs := make([]error, 0, 10)
	if err := checkJiraEnvironments(); err != nil {
		errs = append(errs, checkJiraEnvironments())
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
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
}

func (f *Fjira) bootstrap(args *CliArgs) {
	if args.IssueKey != "" {
		goIntoIssueView(args.IssueKey)
		return
	}
	time.Sleep(500 * time.Millisecond)
	f.app.RunOnAppRoutine(func() {
		goIntoProjectsSearch()
	})
}
