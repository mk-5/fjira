package issues

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/boards"
	"github.com/mk-5/fjira/internal/jira"
	"github.com/mk-5/fjira/internal/statuses"
	"github.com/mk-5/fjira/internal/ui"
	"github.com/mk-5/fjira/internal/users"
	"regexp"
	"strings"
)

type searchIssuesView struct {
	api          jira.Api
	bottomBar    *app.ActionBar
	topBar       *app.ActionBar
	fuzzyFind    *app.FuzzyFind
	project      *jira.Project
	goBackFn     func()
	currentQuery string
	customJql    string
	screenX      int
	screenY      int
	issues       []jira.Issue
	labels       []string
	dirty        bool // refetch jira issues from api if dirty
}

const (
	JiraFetchRecordsThreshold = 100
	topBarStatus              = 1
	topBarAssignee            = 2
	topBarLabel               = 3
)

var (
	issueRegExp     = regexp.MustCompile("^[A-Za-z0-9]{2,10}-[0-9]+$")
	searchForStatus *jira.IssueStatus // global in order to keep status&user between views
	searchForUser   *jira.User
	searchForLabel  string
	searchNavItems  = []ui.NavItemConfig{
		ui.NavItemConfig{Action: ui.ActionSearchByStatus, Text1: ui.MessageByStatus, Text2: "[F1]", Key: tcell.KeyF1},
		ui.NavItemConfig{Action: ui.ActionSearchByAssignee, Text1: ui.MessageByAssignee, Text2: "[F2]", Key: tcell.KeyF2},
		ui.NavItemConfig{Action: ui.ActionSearchByLabel, Text1: ui.MessageByLabel, Text2: "[F3]", Key: tcell.KeyF3},
		ui.NavItemConfig{Action: ui.ActionBoards, Text1: ui.MessageBoards, Text2: "[F4]", Key: tcell.KeyF4},
	}
)

func NewIssuesSearchView(project *jira.Project, goBackFn func(), api jira.Api) app.View {
	bottomBar := ui.CreateBottomActionBarWithItems(searchNavItems)
	topBarItems := []ui.NavItemConfig{
		ui.NavItemConfig{Text1: ui.MessageProjectLabel, Text2: app.ActionBarLabel(fmt.Sprintf("[%s]%s", project.Key, project.Name))},
		ui.NavItemConfig{Text1: ui.MessageLabelStatus, Text2: ui.MessageAll},
		ui.NavItemConfig{Text1: ui.MessageLabelAssignee, Text2: ui.MessageAll},
		ui.NavItemConfig{Text1: ui.MessageLabelLabel, Text2: ui.MessageAll},
	}
	topBar := ui.CreateTopActionBarWithItems(topBarItems)
	return &searchIssuesView{
		api:       api,
		goBackFn:  goBackFn,
		bottomBar: bottomBar,
		topBar:    topBar,
		project:   project,
	}
}

func NewIssuesSearchViewWithCustomJql(jql string, goBackFn func(), api jira.Api) app.View {
	project := &jira.Project{Id: "", Key: ui.MessageCustomJql, Name: ""}
	var jqlTopBar string
	jqlTopBar = jql
	if len(jqlTopBar) > 250 {
		jqlTopBar = jqlTopBar[:250]
	}
	topBar := ui.CreateTopActionBarWithItems([]ui.NavItemConfig{
		ui.NavItemConfig{Text1: ui.MessageJqlLabel, Text2: app.ActionBarLabel(jqlTopBar)},
	})
	return &searchIssuesView{
		api:       api,
		goBackFn:  goBackFn,
		bottomBar: app.NewActionBar(app.Bottom, app.Left),
		topBar:    topBar,
		project:   project,
		customJql: jql,
	}
}

func (view *searchIssuesView) Init() {
	app.GetApp().LoadingWithText(true, ui.MessageSearchIssuesLoading)
	if view.project.Id == ui.MessageAll {
		view.bottomBar.RemoveItem(int(ui.ActionSearchByStatus))
		view.bottomBar.RemoveItem(int(ui.ActionSearchByAssignee))
	}
	go view.runIssuesFuzzyFind()
	go view.handleSearchActions()
}

func (view *searchIssuesView) Destroy() {
	// do nothing
}

func (view *searchIssuesView) Draw(screen tcell.Screen) {
	view.bottomBar.Draw(screen)
	if view.fuzzyFind != nil {
		view.fuzzyFind.Draw(screen)
	}
	view.topBar.Draw(screen)
}

func (view *searchIssuesView) Update() {
	view.bottomBar.Update()
	view.topBar.Update()
	if view.fuzzyFind != nil {
		view.fuzzyFind.Update()
	}
	if searchForStatus != nil && view.topBar.GetItem(topBarStatus).Text2 != searchForStatus.Name {
		view.topBar.GetItem(topBarStatus).ChangeText(ui.MessageLabelStatus, searchForStatus.Name)
		view.topBar.Resize(view.screenX, view.screenY)
	}
	if searchForUser != nil && view.topBar.GetItem(topBarAssignee).Text2 != searchForUser.DisplayName {
		view.topBar.GetItem(topBarAssignee).ChangeText(ui.MessageLabelAssignee, searchForUser.DisplayName)
		view.topBar.Resize(view.screenX, view.screenY)
	}
	if searchForLabel != "" && view.topBar.GetItem(topBarLabel).Text2 != searchForLabel {
		view.topBar.GetItem(topBarLabel).ChangeText(ui.MessageLabelLabel, searchForLabel)
		view.topBar.Resize(view.screenX, view.screenY)
	}
}

func (view *searchIssuesView) Resize(screenX, screenY int) {
	view.bottomBar.Resize(screenX, screenY)
	view.topBar.Resize(screenX, screenY)
	if view.fuzzyFind != nil {
		view.fuzzyFind.Resize(screenX, screenY)
	}
	view.screenX = screenX
	view.screenY = screenY
}

func (view *searchIssuesView) HandleKeyEvent(ev *tcell.EventKey) {
	go view.bottomBar.HandleKeyEvent(ev) // TODO - do not trigger new routine
	if view.fuzzyFind != nil {
		view.fuzzyFind.HandleKeyEvent(ev)
	}
}

func (view *searchIssuesView) runIssuesFuzzyFind() {
	a := app.GetApp()
	view.fuzzyFind = app.NewFuzzyFindWithProvider(ui.MessageSelectIssue, view.findIssues)
	view.fuzzyFind.MarginBottom = 1
	if view.customJql != "" {
		view.fuzzyFind.MarginBottom = 0
	}
	a.Loading(false)
	a.ClearNow()
	if chosen := <-view.fuzzyFind.Complete; true {
		a.ClearNow()
		if chosen.Index < 0 {
			view.goBack()
			searchForStatus = nil
			searchForUser = nil
			return
		}
		chosenIssue := view.issues[chosen.Index]
		go view.goToIssueView(chosenIssue.Key)
	}
}

func (view *searchIssuesView) goToIssueView(issueKey string) {
	app.GoTo("issue", issueKey, view.reopen, view.api)
}

func (view *searchIssuesView) findIssues(query string) []string {
	a := app.GetApp()
	query = strings.TrimSpace(query)

	// when no custom jql set
	// when manual set dirty=true
	// when there is more records than max
	// when query has issue format
	// when there is no results
	if view.customJql == "" || len(view.issues) >= JiraFetchRecordsThreshold || len(view.issues) == 0 || view.dirty || view.queryHasIssueFormat() || query == "" {
		a.LoadingWithText(true, ui.MessageSearchIssuesLoading)
		view.issues = view.searchForIssues(query)
		a.Loading(false)
		view.dirty = false
	}

	view.currentQuery = query
	return FormatJiraIssues(view.issues)
}

func (view *searchIssuesView) handleSearchActions() {
	if selectedAction := <-view.bottomBar.Action; true {
		switch selectedAction {
		case ui.ActionSearchByStatus:
			view.runSelectStatus()
		case ui.ActionSearchByAssignee:
			view.runSelectUser()
		case ui.ActionSearchByLabel:
			view.runSelectLabel()
		case ui.ActionBoards:
			view.runSelectBoard()
		}
	}
}

func (view *searchIssuesView) runSelectStatus() {
	app.GetApp().ClearNow()
	app.GetApp().Loading(true)
	ss := view.fetchStatuses(view.project.Id)
	ss = append(ss, jira.IssueStatus{Name: ui.MessageAll})
	statusesStrings := statuses.FormatJiraStatuses(ss)
	view.fuzzyFind = app.NewFuzzyFind(ui.MessageStatusFuzzyFind, statusesStrings)
	app.GetApp().Loading(false)
	if status := <-view.fuzzyFind.Complete; true {
		app.GetApp().ClearNow()
		if status.Index >= 0 && len(ss) > 0 {
			searchForStatus = &ss[status.Index]
			view.dirty = true
		}
		go view.runIssuesFuzzyFind()
		go view.handleSearchActions()
	}
}

func (view *searchIssuesView) runSelectUser() {
	app.GetApp().ClearNow()
	app.GetApp().Loading(true)
	us := view.fetchUsers(view.project.Key)
	us = append(us, jira.User{DisplayName: ui.MessageAll})
	usersStrings := users.FormatJiraUsers(us)
	view.fuzzyFind = app.NewFuzzyFind(ui.MessageSelectUser, usersStrings)
	app.GetApp().Loading(false)
	if user := <-view.fuzzyFind.Complete; true {
		app.GetApp().ClearNow()
		if user.Index >= 0 && len(us) > 0 {
			searchForUser = &us[user.Index]
			view.dirty = true
		}
		go view.runIssuesFuzzyFind()
		go view.handleSearchActions()
	}
}

func (view *searchIssuesView) runSelectLabel() {
	app.GetApp().ClearNow()
	app.GetApp().Loading(true)
	view.fuzzyFind = app.NewFuzzyFindWithProvider(ui.MessageSelectLabel, view.findLabels)
	app.GetApp().Loading(false)
	if label := <-view.fuzzyFind.Complete; true {
		app.GetApp().ClearNow()
		if label.Index >= 0 && len(view.labels) > 0 {
			searchForLabel = view.labels[label.Index]
			view.dirty = true
		}
		go view.runIssuesFuzzyFind()
		go view.handleSearchActions()
	}
}

func (view *searchIssuesView) runSelectBoard() {
	app.GetApp().ClearNow()
	app.GetApp().Loading(true)
	bs := view.findBoards()
	boardsString := boards.FormatJiraBoards(bs)
	view.fuzzyFind = app.NewFuzzyFind(ui.MessageSelectBoard, boardsString)
	app.GetApp().Loading(false)
	if board := <-view.fuzzyFind.Complete; true {
		app.GetApp().ClearNow()
		if board.Index >= 0 && len(boardsString) > 0 {
			app.GoTo("boards", view.project, &bs[board.Index], view.reopen, view.api)
			return
		}
		go view.runIssuesFuzzyFind()
		go view.handleSearchActions()
	}
}

func (view *searchIssuesView) reopen() {
	app.GoTo("issues-search", view.project.Id, view.goBackFn, view.api)
}

func (view *searchIssuesView) searchForIssues(query string) []jira.Issue {
	q := strings.TrimSpace(query)
	jql := BuildSearchIssuesJql(view.project, q, searchForStatus, searchForUser, searchForLabel)
	// when custom JQL - use it instead of fuzzy query
	if view.customJql != "" {
		jql = view.customJql
	}
	issues, err := view.api.SearchJql(jql)
	if err != nil {
		app.Error(err.Error())
	}
	return issues
}

func (view *searchIssuesView) fetchStatuses(projectId string) []jira.IssueStatus {
	app.GetApp().Loading(true)
	ss, err := view.api.FindProjectStatuses(projectId)
	if err != nil {
		app.Error(err.Error())
	}
	app.GetApp().Loading(false)
	return ss
}

func (view *searchIssuesView) fetchUsers(projectId string) []jira.User {
	us, err := view.api.FindUsers(projectId)
	if err != nil {
		app.Error(err.Error())
	}
	return us
}

func (view *searchIssuesView) findLabels(query string) []string {
	app.GetApp().LoadingWithText(true, ui.MessageSearchLabelsLoading)
	labels, err := view.api.FindLabels(nil, query)
	labels = append(labels, ui.MessageAll)
	app.GetApp().Loading(false)
	if err != nil {
		app.Error(err.Error())
	}
	view.labels = labels
	return labels
}

func (view *searchIssuesView) findBoards() []jira.BoardItem {
	app.GetApp().LoadingWithText(true, ui.MessageSearchBoardsLoading)
	bs, err := view.api.FindBoards(view.project.Id)
	app.GetApp().Loading(false)
	if err != nil {
		app.Error(err.Error())
	}
	return bs
}

func (view *searchIssuesView) queryHasIssueFormat() bool {
	return issueRegExp.MatchString(view.currentQuery)
}

func (view *searchIssuesView) goBack() {
	if view.goBackFn != nil {
		view.goBackFn()
	}
}
