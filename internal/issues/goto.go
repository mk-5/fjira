package issues

import (
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	"time"
)

type SearchArgs struct {
	ProjectKey string
	GoBackFn   func()
	Api        jira.Api
}

const (
	issue           string = "issue"
	issuesSearch    string = "issues-search"
	issuesSearchJql string = "issues-search-jql"
	jql             string = "jql"
)

func RegisterGoTo() {
	app.RegisterGoto(issue, func(args ...interface{}) {
		issueKey := args[0].(string)
		var goBackFn func()
		if fn, ok := args[1].(func()); ok {
			goBackFn = fn
		}
		api := args[2].(jira.Api)

		defer app.GetApp().PanicRecover()
		app.GetApp().Loading(true)
		issue, err := api.GetIssueDetailed(issueKey)
		if err != nil {
			app.GetApp().Loading(false)
			app.Error(err.Error())
			return
		}
		app.GetApp().Loading(false)
		if goBackFn == nil {
			goBackFn = func() {
				app.GoTo(issuesSearch, issue.Fields.Project.Id, func() {
					app.GoTo("projects", api)
				}, api)
			}
		}
		issueView := NewIssueView(issue, goBackFn, api)
		app.GetApp().SetView(issueView)
	})
	app.RegisterGoto(issuesSearch, func(args ...interface{}) {
		projectKey := args[0].(string)
		var goBackFn func()
		if fn, ok := args[1].(func()); ok {
			goBackFn = fn
		}
		api := args[2].(jira.Api)

		defer app.GetApp().PanicRecover()
		app.GetApp().Loading(true)
		project, err := api.FindProject(projectKey)
		if err != nil {
			app.Error(err.Error())
			<-time.NewTimer(2 * time.Second).C
			app.GetApp().Quit()
			return
		}
		app.GetApp().Loading(false)
		projectsView := NewIssuesSearchView(project, goBackFn, api)
		app.GetApp().SetView(projectsView)
	})
	app.RegisterGoto(issuesSearchJql, func(args ...interface{}) {
		defer app.GetApp().PanicRecover()
		jql := args[0].(string)
		api := args[1].(jira.Api)
		issuesSearchView := NewIssuesSearchViewWithCustomJql(jql, func() {
			app.GoTo("jql", api)
		}, api)
		app.GetApp().SetView(issuesSearchView)
	})
	app.RegisterGoto(jql, func(args ...interface{}) {
		defer app.GetApp().PanicRecover()
		api := args[0].(jira.Api)
		jqlView := NewJqlSearchView(api)
		app.GetApp().SetView(jqlView)
	})
}
