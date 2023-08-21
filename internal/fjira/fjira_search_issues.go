package fjira

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	"regexp"
	"strings"
)

type fjiraSearchIssuesView struct {
	bottomBar    *app.ActionBar
	topBar       *app.ActionBar
	fuzzyFind    *app.FuzzyFind
	project      *jira.Project
	currentQuery string
	customJql    string
	screenX      int
	screenY      int
	issues       []jira.Issue
	labels       []string
	queryDirty   bool
}

const (
	JiraRecordsMax = 100
	topBarStatus   = 1
	topBarAssignee = 2
	topBarLabel    = 3
)

var (
	issueRegExp     = regexp.MustCompile("^[A-Za-z0-9]{2,10}-[0-9]+$")
	searchForStatus *jira.IssueStatus // global in order to keep status&user between views
	searchForUser   *jira.User
	searchForLabel  string
)

func NewIssuesSearchView(project *jira.Project) *fjiraSearchIssuesView {
	bottomBar := CreateSearchIssuesBottomBar()
	topBar := CreateSearchIssuesTopBar(project)
	return &fjiraSearchIssuesView{
		bottomBar: bottomBar,
		topBar:    topBar,
		project:   project,
	}
}

func NewIssuesSearchViewWithCustomJql(jql string) *fjiraSearchIssuesView {
	project := &jira.Project{Id: "", Key: MessageCustomJql, Name: ""}
	topBar := CreateCustomJqlTopBar(jql)
	return &fjiraSearchIssuesView{
		bottomBar: app.NewActionBar(app.Bottom, app.Left),
		topBar:    topBar,
		project:   project,
		customJql: jql,
	}
}

func (view *fjiraSearchIssuesView) Init() {
	app.GetApp().LoadingWithText(true, MessageSearchIssuesLoading)
	if view.project.Id == MessageAll {
		view.bottomBar.RemoveItem(int(ActionSearchByStatus))
		view.bottomBar.RemoveItem(int(ActionSearchByAssignee))
	}
	go view.runIssuesFuzzyFind()
	go view.handleSearchActions()
}

func (view *fjiraSearchIssuesView) Destroy() {
	// do nothing
}

func (view *fjiraSearchIssuesView) Draw(screen tcell.Screen) {
	view.bottomBar.Draw(screen)
	if view.fuzzyFind != nil {
		view.fuzzyFind.Draw(screen)
	}
	view.topBar.Draw(screen)
}

func (view *fjiraSearchIssuesView) Update() {
	view.bottomBar.Update()
	view.topBar.Update()
	if view.fuzzyFind != nil {
		view.fuzzyFind.Update()
	}
	if searchForStatus != nil && view.topBar.GetItem(topBarStatus).Text2 != searchForStatus.Name {
		view.topBar.GetItem(topBarStatus).ChangeText(MessageLabelStatus, searchForStatus.Name)
		view.topBar.Resize(view.screenX, view.screenY)
	}
	if searchForUser != nil && view.topBar.GetItem(topBarAssignee).Text2 != searchForUser.DisplayName {
		view.topBar.GetItem(topBarAssignee).ChangeText(MessageLabelAssignee, searchForUser.DisplayName)
		view.topBar.Resize(view.screenX, view.screenY)
	}
	if searchForLabel != "" && view.topBar.GetItem(topBarLabel).Text2 != searchForLabel {
		view.topBar.GetItem(topBarLabel).ChangeText(MessageLabelLabel, searchForLabel)
		view.topBar.Resize(view.screenX, view.screenY)
	}
}

func (view *fjiraSearchIssuesView) Resize(screenX, screenY int) {
	view.bottomBar.Resize(screenX, screenY)
	view.topBar.Resize(screenX, screenY)
	if view.fuzzyFind != nil {
		view.fuzzyFind.Resize(screenX, screenY)
	}
	view.screenX = screenX
	view.screenY = screenY
}

func (view *fjiraSearchIssuesView) HandleKeyEvent(ev *tcell.EventKey) {
	go view.bottomBar.HandleKeyEvent(ev) // TODO - do not trigger new routine
	if view.fuzzyFind != nil {
		view.fuzzyFind.HandleKeyEvent(ev)
	}
}

func (view *fjiraSearchIssuesView) runIssuesFuzzyFind() {
	a := app.GetApp()
	view.fuzzyFind = app.NewFuzzyFindWithProvider(MessageSelectIssue, view.findIssues)
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

func (view *fjiraSearchIssuesView) goToIssueView(issueKey string) {
	goIntoIssueView(issueKey)
	if view.customJql != "" {
		if v, ok := app.GetApp().CurrentView().(*fjiraIssueView); ok {
			v.SetGoBackJql(view.customJql)
		}
	}
}

func (view *fjiraSearchIssuesView) findIssues(query string) []string {
	formatter, _ := GetFormatter()
	a := app.GetApp()

	// when no custom jql set
	// when manual set queryDirty=true
	// when there is more records than max
	// when backspace
	// when query has issue format
	// when there is no results
	if !(view.customJql != "" && len(view.issues) > 0) && view.queryDirty || len(view.issues) >= JiraRecordsMax || len(query) < len(view.currentQuery) || view.queryHasIssueFormat() || len(view.issues) == 0 {
		view.queryDirty = false
		a.LoadingWithText(true, MessageSearchIssuesLoading)
		view.issues = view.searchForIssues(query)
		a.Loading(false)
	}

	view.currentQuery = query
	return formatter.formatJiraIssues(view.issues)
}

func (view *fjiraSearchIssuesView) handleSearchActions() {
	if selectedAction := <-view.bottomBar.Action; true {
		switch selectedAction {
		case ActionSearchByStatus:
			view.runSelectStatus()
		case ActionSearchByAssignee:
			view.runSelectUser()
		case ActionSearchByLabel:
			view.runSelectLabel()
		case ActionBoards:
			view.runSelectBoard()
		}
	}
}

func (view *fjiraSearchIssuesView) runSelectStatus() {
	app.GetApp().ClearNow()
	app.GetApp().Loading(true)
	formatter, _ := GetFormatter()
	statuses := view.fetchStatuses(view.project.Id)
	statuses = append(statuses, jira.IssueStatus{Name: MessageAll})
	statusesStrings := formatter.formatJiraStatuses(statuses)
	view.fuzzyFind = app.NewFuzzyFind(MessageStatusFuzzyFind, statusesStrings)
	app.GetApp().Loading(false)
	if status := <-view.fuzzyFind.Complete; true {
		app.GetApp().ClearNow()
		if status.Index >= 0 && len(statuses) > 0 {
			searchForStatus = &statuses[status.Index]
			view.queryDirty = true
		}
		go view.runIssuesFuzzyFind()
		go view.handleSearchActions()
	}
}

func (view *fjiraSearchIssuesView) runSelectUser() {
	app.GetApp().ClearNow()
	app.GetApp().Loading(true)
	formatter, _ := GetFormatter()
	users := view.fetchUsers(view.project.Key)
	users = append(users, jira.User{DisplayName: MessageAll})
	usersStrings := formatter.formatJiraUsers(users)
	view.fuzzyFind = app.NewFuzzyFind(MessageSelectUser, usersStrings)
	app.GetApp().Loading(false)
	if user := <-view.fuzzyFind.Complete; true {
		app.GetApp().ClearNow()
		if user.Index >= 0 && len(users) > 0 {
			searchForUser = &users[user.Index]
			view.queryDirty = true
		}
		go view.runIssuesFuzzyFind()
		go view.handleSearchActions()
	}
}

func (view *fjiraSearchIssuesView) runSelectLabel() {
	app.GetApp().ClearNow()
	app.GetApp().Loading(true)
	view.fuzzyFind = app.NewFuzzyFindWithProvider(MessageSelectLabel, view.findLabels)
	app.GetApp().Loading(false)
	if label := <-view.fuzzyFind.Complete; true {
		app.GetApp().ClearNow()
		if label.Index >= 0 && len(view.labels) > 0 {
			searchForLabel = view.labels[label.Index]
			view.queryDirty = true
		}
		go view.runIssuesFuzzyFind()
		go view.handleSearchActions()
	}
}

func (view *fjiraSearchIssuesView) runSelectBoard() {
	app.GetApp().ClearNow()
	app.GetApp().Loading(true)
	formatter, _ := GetFormatter()
	boards := view.findBoards()
	boardsString := formatter.formatJiraBoards(boards)
	view.fuzzyFind = app.NewFuzzyFind(MessageSelectBoard, boardsString)
	app.GetApp().Loading(false)
	if board := <-view.fuzzyFind.Complete; true {
		app.GetApp().ClearNow()
		if board.Index >= 0 && len(boardsString) > 0 {
			goIntoBoardView(view.project, &boards[board.Index])
			return
		}
		go view.runIssuesFuzzyFind()
		go view.handleSearchActions()
	}
}

func (view *fjiraSearchIssuesView) searchForIssues(query string) []jira.Issue {
	q := strings.TrimSpace(query)
	api, _ := GetApi()
	jql := buildSearchIssuesJql(view.project, q, searchForStatus, searchForUser, searchForLabel)
	// when custom JQL - use it instead of fuzzy query
	if view.customJql != "" {
		jql = view.customJql
	}
	issues, err := api.SearchJql(jql)
	if err != nil {
		app.Error(err.Error())
	}
	return issues
}

func (view *fjiraSearchIssuesView) fetchStatuses(projectId string) []jira.IssueStatus {
	api, _ := GetApi()
	app.GetApp().Loading(true)
	statuses, err := api.FindProjectStatuses(projectId)
	if err != nil {
		app.Error(err.Error())
	}
	app.GetApp().Loading(false)
	return statuses
}

func (view *fjiraSearchIssuesView) fetchUsers(projectId string) []jira.User {
	api, _ := GetApi()
	users, err := api.FindUsers(projectId)
	if err != nil {
		app.Error(err.Error())
	}
	return users
}

func (view *fjiraSearchIssuesView) findLabels(query string) []string {
	api, _ := GetApi()
	app.GetApp().LoadingWithText(true, MessageSearchLabelsLoading)
	labels, err := api.FindLabels(nil, query)
	labels = append(labels, MessageAll)
	app.GetApp().Loading(false)
	if err != nil {
		app.Error(err.Error())
	}
	view.labels = labels
	return labels
}

func (view *fjiraSearchIssuesView) findBoards() []jira.BoardItem {
	api, _ := GetApi()
	app.GetApp().LoadingWithText(true, MessageSearchBoardsLoading)
	boards, err := api.FindBoards(view.project.Id)
	app.GetApp().Loading(false)
	if err != nil {
		app.Error(err.Error())
	}
	return boards
}

func (view *fjiraSearchIssuesView) queryHasIssueFormat() bool {
	return issueRegExp.MatchString(view.currentQuery)
}

func (view *fjiraSearchIssuesView) goBack() {
	if view.customJql == "" {
		go goIntoProjectsSearch()
		return
	}
	if view.customJql != "" {
		go goIntoJqlView()
	}
}
