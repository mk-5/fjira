package users

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	"github.com/mk-5/fjira/internal/ui"
)

type userAssignChangeView struct {
	app.View
	api       jira.Api
	bottomBar *app.ActionBar
	topBar    *app.ActionBar
	fuzzyFind *app.FuzzyFind
	issue     *jira.Issue
	goBackFn  func()
}

func NewAssignChangeView(issue *jira.Issue, goBackFn func(), api jira.Api) app.View {
	return &userAssignChangeView{
		api:       api,
		issue:     issue,
		topBar:    ui.CreateIssueTopBar(issue),
		bottomBar: ui.CreateBottomLeftBar(),
		goBackFn:  goBackFn,
	}
}

func (view *userAssignChangeView) Init() {
	go view.startUsersSearching()
}

func (view *userAssignChangeView) Destroy() {
	// do nothing
}

func (view *userAssignChangeView) Draw(screen tcell.Screen) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Draw(screen)
	}
	view.topBar.Draw(screen)
	view.bottomBar.Draw(screen)
}

func (view *userAssignChangeView) Update() {
	view.bottomBar.Update()
	if view.fuzzyFind != nil {
		view.fuzzyFind.Update()
	}
}

func (view *userAssignChangeView) Resize(screenX, screenY int) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Resize(screenX, screenY)
	}
	view.topBar.Resize(screenX, screenY)
	view.bottomBar.Resize(screenX, screenY)
}

func (view *userAssignChangeView) HandleKeyEvent(ev *tcell.EventKey) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.HandleKeyEvent(ev)
	}
}

func (view *userAssignChangeView) startUsersSearching() {
	app.GetApp().ClearNow()
	app.GetApp().Loading(true)
	var us *[]jira.User
	view.fuzzyFind, us = NewFuzzyFind(view.issue.Fields.Project.Key, view.api)
	view.fuzzyFind.MarginBottom = 0
	app.GetApp().Loading(false)
	if user := <-view.fuzzyFind.Complete; true {
		app.GetApp().ClearNow()
		if user.Index < 0 {
			if view.goBackFn != nil {
				view.goBackFn()
			}
			return
		}
		view.fuzzyFind = nil
		view.assignUserToTicket(view.issue, &(*us)[user.Index])
	}
}

func (view *userAssignChangeView) assignUserToTicket(issue *jira.Issue, user *jira.User) {
	if user == nil {
		view.goBackFn()
		return
	}
	message := fmt.Sprintf(ui.MessageChangingAssigneeTo, issue.Key, user.DisplayName)
	app.GetApp().ClearNow()
	view.bottomBar.AddItem(ui.NewYesBarItem())
	view.bottomBar.AddItem(ui.NewCancelBarItem())
	// TODO - should confirm be also drawable? at the moment yes/no are rendered out of the confirm thingy..
	userAssign := app.Confirm(app.GetApp(), message)
	if userAssign {
		view.doAssignmentChange(issue, user)
	}
	if view.goBackFn != nil {
		view.goBackFn()
	}
}

func (view *userAssignChangeView) doAssignmentChange(issue *jira.Issue, user *jira.User) {
	app.GetApp().LoadingWithText(true, ui.MessageAssigningUser)
	err := view.api.DoAssignee(issue.Key, user)
	app.GetApp().Loading(false)
	if err != nil {
		app.Error(fmt.Sprintf(ui.MessageCannotAssignUser, user.DisplayName, issue.Key, err.Error(), user.AccountId))
		return
	}
	app.Success(fmt.Sprintf(ui.MessageAssignSuccess, user.DisplayName, issue.Key))
}
