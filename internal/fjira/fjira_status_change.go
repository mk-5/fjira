package fjira

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
)

type fjiraStatusChangeView struct {
	app.View
	topBar    *app.ActionBar
	bottomBar *app.ActionBar
	fuzzyFind *app.FuzzyFind
	issue     *jira.Issue
}

func NewStatusChangeView(issue *jira.Issue) *fjiraStatusChangeView {
	return &fjiraStatusChangeView{
		issue:     issue,
		topBar:    CreateIssueTopBar(issue),
		bottomBar: CreateBottomLeftBar(),
	}
}

func (view *fjiraStatusChangeView) Init() {
	go view.startStatusSearching()
}

func (view *fjiraStatusChangeView) Destroy() {
	// do nothing
}

func (view *fjiraStatusChangeView) Draw(screen tcell.Screen) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Draw(screen)
	}
	view.topBar.Draw(screen)
	view.bottomBar.Draw(screen)
}

func (view *fjiraStatusChangeView) Update() {
	view.topBar.Update()
	view.bottomBar.Update()
	if view.fuzzyFind != nil {
		view.fuzzyFind.Update()
	}
}

func (view *fjiraStatusChangeView) Resize(screenX, screenY int) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Resize(screenX, screenY)
	}
	view.topBar.Resize(screenX, screenY)
	view.bottomBar.Resize(screenX, screenY)
}

func (view *fjiraStatusChangeView) HandleKeyEvent(ev *tcell.EventKey) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.HandleKeyEvent(ev)
	}
}

func (view *fjiraStatusChangeView) startStatusSearching() {
	app.GetApp().ClearNow()
	app.GetApp().Loading(true)
	formatter, _ := GetFormatter()
	statuses := view.transitions(view.issue.Id)
	statusesStrings := formatter.formatJiraTransitions(statuses)
	view.fuzzyFind = app.NewFuzzyFind(MessageStatusFuzzyFind, statusesStrings)
	view.fuzzyFind.MarginBottom = 0
	app.GetApp().Loading(false)
	if status := <-view.fuzzyFind.Complete; true {
		app.GetApp().ClearNow()
		if status.Index < 0 {
			app.GetApp().SetView(newIssueView(view.issue))
			return
		}
		view.fuzzyFind = nil
		view.changeStatusTo(&statuses[status.Index])
	}
}

func (view *fjiraStatusChangeView) changeStatusTo(status *jira.IssueTransition) {
	message := fmt.Sprintf(MessageChangingStatusTo, view.issue.Key, status.Name)
	app.GetApp().ClearNow()
	//view.bottomBar.AddItem(NewNewStatusBarItem(status.Name))
	view.bottomBar.AddItem(NewYesBarItem())
	view.bottomBar.AddItem(NewCancelBarItem())
	changeStatus := app.Confirm(app.GetApp(), message)
	switch changeStatus {
	case true:
		view.changeStatusForTicket(view.issue, status)
		goIntoIssueView(view.issue.Key)
	case false:
		app.GetApp().SetView(newIssueView(view.issue))
	}
}

func (view *fjiraStatusChangeView) transitions(issueId string) []jira.IssueTransition {
	api, _ := GetApi()
	transitions, _ := api.FindTransitions(issueId)
	return transitions
}

func (view *fjiraStatusChangeView) changeStatusForTicket(issue *jira.Issue, status *jira.IssueTransition) {
	app.GetApp().ClearNow()
	app.GetApp().LoadingWithText(true, MessageChangingStatus)
	api, _ := GetApi()
	err := api.DoTransition(issue.Key, status)
	app.GetApp().Loading(false)
	if err != nil {
		app.Error(err.Error())
		return
	}
	app.Success(fmt.Sprintf(MessageChangeStatusSuccess, issue.Key, status.Name))
}
