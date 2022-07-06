package fjira

import (
	"github.com/mk5/fjira/internal/app"
	"github.com/mk5/fjira/internal/jira"
	"time"
)

func goIntoProjectsSearch() {
	projectsView := NewProjectsSearchView()
	app.GetApp().SetView(projectsView)
}

func goIntoIssuesSearchForProject(projectKey string) {
	app.GetApp().Loading(true)
	api, _ := GetApi()
	project, err := api.FindProject(projectKey)
	if err != nil {
		app.Error(err.Error())
		<-time.NewTimer(2 * time.Second).C
		app.GetApp().Quit()
		return
	}
	app.GetApp().Loading(false)
	projectsView := NewIssuesSearchView(project)
	app.GetApp().SetView(projectsView)
}

func goIntoIssuesSearch(project *jira.JiraProject) {
	issuesSearchView := NewIssuesSearchView(project)
	app.GetApp().SetView(issuesSearchView)
}

func goIntoIssueView(issueKey string) {
	defer app.GetApp().PanicRecover()
	app.GetApp().Loading(true)
	api, _ := GetApi()
	issue, err := api.GetIssueDetailed(issueKey)
	if err != nil {
		app.Error(err.Error())
	}
	app.GetApp().Loading(false)
	issueView := NewIssueView(issue)
	app.GetApp().SetView(issueView)
}

func goIntoChangeStatus(issue *jira.JiraIssue) {
	statusChangeView := NewStatusChangeView(issue)
	app.GetApp().SetView(statusChangeView)
}

func goIntoChangeAssignment(issue *jira.JiraIssue) {
	assignChangeView := NewAssignChangeView(issue)
	app.GetApp().SetView(assignChangeView)
}

func goIntoCommentView(issue *jira.JiraIssue) {
	commentView := NewCommentView(issue)
	app.GetApp().SetView(commentView)
}

func goIntoSwitchWorkspaceView() {
	switchWorkspaceView := NewSwitchWorkspaceView()
	app.GetApp().SetView(switchWorkspaceView)
}
