package issues

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	"github.com/mk-5/fjira/internal/ui"
	"strings"
	"time"
)

const (
	DefaultJqlQuery = "created >= -30d order by created DESC"
	BadRequest      = "Bad Request"
)

type jqlSearchView struct {
	app.View
	api       jira.Api
	fuzzyFind *app.FuzzyFind
	issues    []jira.Issue
	jql       string
}

func NewJqlSearchView(api jira.Api) app.View {
	return &jqlSearchView{
		api: api,
		jql: DefaultJqlQuery,
	}
}

func (view *jqlSearchView) Init() {
	go view.startJqlFuzzyFind()
}

func (view *jqlSearchView) Destroy() {
	// do nothing
}

func (view *jqlSearchView) Draw(screen tcell.Screen) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Draw(screen)
	}
}

func (view *jqlSearchView) Update() {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Update()
	}
}

func (view *jqlSearchView) Resize(screenX, screenY int) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Resize(screenX, screenY)
	}
}

func (view *jqlSearchView) HandleKeyEvent(ev *tcell.EventKey) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.HandleKeyEvent(ev)
	}
}

func (view *jqlSearchView) startJqlFuzzyFind() {
	app.GetApp().ClearNow()
	app.GetApp().Loading(true)
	view.fuzzyFind = app.NewFuzzyFindWithProvider(ui.MessageJqlFuzzyFind, view.findIssues)
	view.fuzzyFind.MarginBottom = 0
	view.fuzzyFind.SetQuery(DefaultJqlQuery)
	view.fuzzyFind.AlwaysShowAllResults()
	// higher debounce in order to give more time to change jql
	view.fuzzyFind.SetDebounceMs(500 * time.Millisecond)
	app.GetApp().Loading(false)
	if chosen := <-view.fuzzyFind.Complete; true {
		app.GetApp().ClearNow()
		query := view.fuzzyFind.GetQuery()
		if chosen.Index < 0 && strings.TrimSpace(query) == "" {
			// do nothing
			return
		}
		if chosen.Index >= 0 {
			chosenIssue := view.issues[chosen.Index]
			app.GoTo("issue", chosenIssue.Key, view.reopen, view.api)
			return
		}
	}
}

func (view *jqlSearchView) reopen() {
	app.GoTo("jql", view.api)
}

func (view *jqlSearchView) findIssues(query string) []string {
	app.GetApp().LoadingWithText(true, ui.MessageSearchIssuesLoading)
	issues, err := view.api.SearchJql(query)
	app.GetApp().Loading(false)
	if err != nil && strings.Contains(err.Error(), BadRequest) {
		// do nothing, invalid JQL query
		return FormatJiraIssues(view.issues)
	}
	if err != nil {
		app.Error(err.Error())
		return FormatJiraIssues(view.issues)
	}
	view.issues = issues
	return FormatJiraIssues(view.issues)
}
