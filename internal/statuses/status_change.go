package statuses

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	"github.com/mk-5/fjira/internal/ui"
)

type statusChangeView struct {
	app.View
	api       jira.Api
	topBar    *app.ActionBar
	bottomBar *app.ActionBar
	fuzzyFind *app.FuzzyFind
	issue     *jira.Issue
	goBackFn  func()
}

func NewStatusChangeView(issue *jira.Issue, goBackFn func(), api jira.Api) app.View {
	return &statusChangeView{
		api:       api,
		goBackFn:  goBackFn,
		issue:     issue,
		topBar:    ui.CreateIssueTopBar(issue),
		bottomBar: ui.CreateBottomLeftBar(),
	}
}

func (view *statusChangeView) Init() {
	go view.startStatusSearching()
}

func (view *statusChangeView) Destroy() {
	// do nothing
}

func (view *statusChangeView) Draw(screen tcell.Screen) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Draw(screen)
	}
	view.topBar.Draw(screen)
	view.bottomBar.Draw(screen)
}

func (view *statusChangeView) Update() {
	view.topBar.Update()
	view.bottomBar.Update()
	if view.fuzzyFind != nil {
		view.fuzzyFind.Update()
	}
}

func (view *statusChangeView) Resize(screenX, screenY int) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Resize(screenX, screenY)
	}
	view.topBar.Resize(screenX, screenY)
	view.bottomBar.Resize(screenX, screenY)
}

func (view *statusChangeView) HandleKeyEvent(ev *tcell.EventKey) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.HandleKeyEvent(ev)
	}
}

func (view *statusChangeView) startStatusSearching() {
	app.GetApp().ClearNow()
	app.GetApp().Loading(true)
	statuses := view.transitions(view.issue.Id)
	statusesStrings := FormatJiraTransitions(statuses)
	view.fuzzyFind = app.NewFuzzyFind(ui.MessageStatusFuzzyFind, statusesStrings)
	view.fuzzyFind.MarginBottom = 0
	app.GetApp().Loading(false)
	if status := <-view.fuzzyFind.Complete; true {
		app.GetApp().ClearNow()
		if status.Index < 0 {
			if view.goBackFn != nil {
				view.goBackFn()
			}
			return
		}
		view.fuzzyFind = nil
		view.changeStatusTo(&statuses[status.Index])
	}
}

func (view *statusChangeView) changeStatusTo(status *jira.IssueTransition) {
	message := fmt.Sprintf(ui.MessageChangingStatusTo, view.issue.Key, status.Name)
	app.GetApp().ClearNow()
	//view.bottomBar.AddItem(NewNewStatusBarItem(status.Name))
	view.bottomBar.AddItem(ui.NewYesBarItem())
	view.bottomBar.AddItem(ui.NewCancelBarItem())
	changeStatus := app.Confirm(app.GetApp(), message)
	if changeStatus {
		view.changeStatusForTicket(view.issue, status)
	}
	if view.goBackFn != nil {
		view.goBackFn()
	}
}

func (view *statusChangeView) transitions(issueId string) []jira.IssueTransition {
	transitions, _ := view.api.FindTransitions(issueId)
	return transitions
}

func (view *statusChangeView) changeStatusForTicket(issue *jira.Issue, status *jira.IssueTransition) {
	app.GetApp().ClearNow()
	app.GetApp().LoadingWithText(true, ui.MessageChangingStatus)
	err := view.api.DoTransition(issue.Key, status)
	app.GetApp().Loading(false)
	if err != nil {
		app.Error(err.Error())
		return
	}
	app.Success(fmt.Sprintf(ui.MessageChangeStatusSuccess, issue.Key, status.Name))
}
