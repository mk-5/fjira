package fjira

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/mk5/fjira/internal/app"
	"github.com/mk5/fjira/internal/jira"
	"log"
	"regexp"
	"strings"
)

type fjiraSearchIssuesView struct {
	bottomBar       *app.ActionBar
	fuzzyFind       *app.FuzzyFind
	project         *jira.JiraProject
	currentQuery    string
	searchForStatus *jira.JiraIssueStatus
}

const (
	MessageSearchIssuesLoading = "Fetching"
	MessageSelectIssue         = "Select issue or ESC to cancel"
	JiraRecordsMax             = 100
	StatusAll                  = "All"
	// TODO - improve query .. would be nice to search via summary/status/assignee
	JqlSummaryQuery = "AND summary~\"%s*\""
	JqlSearchQuery  = "project=%s %s ORDER BY status"
)

var (
	issueRegExp = regexp.MustCompile("^[A-Za-z0-9]{2,10}-[0-9]+$")
)

func NewIssuesSearchView(project *jira.JiraProject) *fjiraSearchIssuesView {
	bottomBar := CreateNewProjectBottomBar(project)
	return &fjiraSearchIssuesView{
		bottomBar: bottomBar,
		project:   project,
	}
}

func (view *fjiraSearchIssuesView) Init() {
	app.GetApp().LoadingWithText(true, MessageSearchIssuesLoading)
	go view.runIssuesFuzzyFind()
	go view.handleSearchActions()
}

func (view *fjiraSearchIssuesView) Destroy() {

}

func (view *fjiraSearchIssuesView) Draw(screen tcell.Screen) {
	view.bottomBar.Draw(screen)
	if view.fuzzyFind != nil {
		view.fuzzyFind.Draw(screen)
	}
}

func (view *fjiraSearchIssuesView) Update() {
	view.bottomBar.Update()
	if view.fuzzyFind != nil {
		view.fuzzyFind.Update()
	}
}

func (view *fjiraSearchIssuesView) Resize(screenX, screenY int) {
	view.bottomBar.Resize(screenX, screenY)
	if view.fuzzyFind != nil {
		view.fuzzyFind.Resize(screenX, screenY)
	}
}

func (view *fjiraSearchIssuesView) HandleKeyEvent(ev *tcell.EventKey) {
	view.bottomBar.HandleKeyEvent(ev)
	if view.fuzzyFind != nil {
		view.fuzzyFind.HandleKeyEvent(ev)
	}
}

func (view *fjiraSearchIssuesView) runIssuesFuzzyFind() {
	formatter, _ := GetFormatter()
	a := app.GetApp()
	latestRecords := view.searchForIssues("")
	// TODO - maybe we should have some additional condition here ..
	// TODO - there is a problem when there is no match from JQL but it's from fuzzy matcher
	issuesProvider := func(query string) []string {
		if len(latestRecords) >= JiraRecordsMax || len(query) < len(view.currentQuery) || view.queryHasIssueFormat() {
			a.LoadingWithText(true, MessageSearchIssuesLoading)
			latestRecords = view.searchForIssues(query)
			a.Loading(false)
		}
		view.currentQuery = query
		return formatter.formatJiraIssues(latestRecords)
	}
	view.fuzzyFind = app.NewFuzzyFindWithProvider(MessageSelectIssue, issuesProvider)
	a.Loading(false)
	a.ClearNow()
	select {
	case chosen := <-view.fuzzyFind.Complete:
		a.ClearNow()
		if chosen.Index < 0 {
			go goIntoProjectsSearch()
			return
		}
		chosenIssue := latestRecords[chosen.Index]
		go goIntoIssueViewFetchIssue(chosenIssue.Key)
	}
}

func (view *fjiraSearchIssuesView) queryHasIssueFormat() bool {
	return issueRegExp.MatchString(view.currentQuery)
}

func (view *fjiraSearchIssuesView) runSelectStatus() {
	app.GetApp().ClearNow()
	app.GetApp().Loading(true)
	formatter, _ := GetFormatter()
	statuses := view.statuses(view.project.Id)
	statuses = append(statuses, jira.JiraIssueStatus{Name: "All"})
	statusesStrings := formatter.formatJiraStatuses(statuses)
	view.fuzzyFind = app.NewFuzzyFind(MessageStatusFuzzyFind, statusesStrings)
	app.GetApp().Loading(false)
	select {
	case status := <-view.fuzzyFind.Complete:
		app.GetApp().ClearNow()
		if status.Index >= 0 {
			view.searchForStatus = &statuses[status.Index]
		}
		view.runIssuesFuzzyFind()
	}
}

func (view *fjiraSearchIssuesView) search(query string) []jira.JiraIssue {
	api, _ := GetApi()
	issues, _, err := api.Search(query)
	if err != nil {
		log.Fatalln(err)
	}
	return issues
}

func (view *fjiraSearchIssuesView) searchForIssues(query string) []jira.JiraIssue {
	q := strings.TrimSpace(query)
	api, _ := GetApi()
	jql := view.buildJql(q)
	issues, err := api.SearchJql(jql)
	if err != nil {
		app.Error(err.Error())
	}
	return issues
}

// TODO - jql builder?
func (view *fjiraSearchIssuesView) buildJql(query string) string {
	jql := fmt.Sprintf("project=%s", view.project.Id)
	orderBy := "ORDER BY status"
	query = strings.TrimSpace(query)
	if query != "" {
		jql = jql + fmt.Sprintf(" AND summary~\"%s*\"", query)
	}
	if view.searchForStatus != nil && view.searchForStatus.Name == StatusAll {
		view.searchForStatus = nil
	}
	if view.searchForStatus != nil {
		jql = jql + fmt.Sprintf(" AND status=%s", view.searchForStatus.Id)
	}
	if query != "" && issueRegExp.MatchString(query) {
		jql = jql + fmt.Sprintf(" OR issuekey=\"%s\"", query)
	}
	return fmt.Sprintf("%s %s", jql, orderBy)
}

func (view *fjiraSearchIssuesView) statuses(projectId string) []jira.JiraIssueStatus {
	api, _ := GetApi()
	statuses, _ := api.FindProjectStatuses(projectId)
	return statuses
}

func (view *fjiraSearchIssuesView) handleSearchActions() {
	select {
	case selectedAction := <-view.bottomBar.Action:
		switch selectedAction {
		case ActionStatusChange:
			view.runSelectStatus()
			return
		}
	}
}
