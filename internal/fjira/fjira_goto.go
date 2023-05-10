package fjira

import (
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
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

func goIntoIssuesSearch(project *jira.Project) {
	issuesSearchView := NewIssuesSearchView(project)
	app.GetApp().SetView(issuesSearchView)
}

func goIntoIssuesSearchForJql(jql string) {
	issuesSearchView := NewIssuesSearchViewWithCustomJql(jql)
	app.GetApp().SetView(issuesSearchView)
}

func goIntoJqlView() {
	jqlView := NewJqlSearchView()
	app.GetApp().SetView(jqlView)
}

func goIntoIssueView(issueKey string) {
	defer app.GetApp().PanicRecover()
	app.GetApp().Loading(true)
	api, _ := GetApi()
	issue, err := api.GetIssueDetailed(issueKey)
	if err != nil {
		app.GetApp().Loading(false)
		app.Error(err.Error())
		return
	}
	app.GetApp().Loading(false)
	issueView := newIssueView(issue)
	app.GetApp().SetView(issueView)
}

func goIntoChangeStatus(issue *jira.Issue) {
	statusChangeView := NewStatusChangeView(issue)
	app.GetApp().SetView(statusChangeView)
}

func goIntoChangeAssignment(issue *jira.Issue) {
	assignChangeView := NewAssignChangeView(issue)
	app.GetApp().SetView(assignChangeView)
}
func goIntoTextWriterView(args *textWriterArgs) {
	view := newTextWriterView(args)
	app.GetApp().SetView(view)
}

func goIntoAddLabelView(issue *jira.Issue) {
	commentView := newAddLabelView(issue)
	app.GetApp().SetView(commentView)
}

func goIntoSwitchWorkspaceView() {
	switchWorkspaceView := newSwitchWorkspaceView()
	app.GetApp().SetView(switchWorkspaceView)
}

func goIntoBoardView(project *jira.Project, board *jira.BoardItem) {
	defer app.GetApp().PanicRecover()
	app.GetApp().Loading(true)
	api, _ := GetApi()
	boardConfig, err := api.GetBoardConfiguration(board.Id)
	if err != nil {
		app.GetApp().Loading(false)
		app.Error(err.Error())
		return
	}
	app.GetApp().Loading(false)
	boardView := newBoardView(project, boardConfig)
	app.GetApp().SetView(boardView)
}
