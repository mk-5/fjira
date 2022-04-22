package fjira

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/mk5/fjira/internal/app"
	"github.com/mk5/fjira/internal/jira"
	"log"
	"strings"
)

type fjiraSearchIssuesView struct {
	bottomBar *app.ActionBar
	fuzzyFind *app.FuzzyFind
	project   *jira.JiraProject
}

const (
	MessageSearchIssuesLoading = "Fetching"
	MessageSelectIssue         = "Select issue or ESC to cancel"
	JiraRecordsMax             = 100
	// TODO - improve query .. would be nice to search via summary/status/assignee
	JqlSummaryQuery = "AND summary~\"%s*\""
	JqlSearchQuery  = "project=%s %s ORDER BY status"
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
	latestRecords := view.searchProject("")
	lastQuery := ""
	// TODO - maybe we should have some additional condition here ..
	// TODO - there is a problem when there is no match from JQL but it's from fuzzy matcher
	issuesSupplier := func(query string) []string {
		if len(latestRecords) >= JiraRecordsMax || len(query) < len(lastQuery) {
			a.LoadingWithText(true, MessageSearchIssuesLoading)
			latestRecords = view.searchProject(query)
			a.Loading(false)
		}
		lastQuery = query
		return formatter.formatJiraIssues(latestRecords)
	}
	view.fuzzyFind = app.NewFuzzyFindWithSupplier(MessageSelectIssue, issuesSupplier)
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
		//go goIntoIssueView(&chosenIssue)
		go goIntoIssueViewFetchIssue(chosenIssue.Key)
	}
}

func (view *fjiraSearchIssuesView) search(query string) []jira.JiraIssue {
	api, err := GetApi()
	issues, _, err := api.Search(query)
	if err != nil {
		log.Fatalln(err)
	}
	return issues
}

func (view *fjiraSearchIssuesView) searchProject(query string) []jira.JiraIssue {
	// TODO query encoding
	summaryWhere := ""
	q := strings.TrimSpace(query)
	if q != "" {
		summaryWhere = fmt.Sprintf(JqlSummaryQuery, q)
	}
	api, err := GetApi()
	jql := fmt.Sprintf(JqlSearchQuery, view.project.Id, summaryWhere)
	issues, err := api.SearchJql(jql)
	if err != nil {
		log.Fatalln(err)
	}
	return issues
}
