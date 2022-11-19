package fjira

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mk5/fjira/internal/app"
	"github.com/mk5/fjira/internal/jira"
)

type fjiraAssignChangeView struct {
	app.View
	bottomBar *app.ActionBar
	topBar    *app.ActionBar
	fuzzyFind *app.FuzzyFind
	issue     *jira.Issue
}

func NewAssignChangeView(issue *jira.Issue) *fjiraAssignChangeView {
	return &fjiraAssignChangeView{
		issue:     issue,
		topBar:    CreateIssueTopBar(issue),
		bottomBar: CreateBottomLeftBar(),
	}
}

func (view *fjiraAssignChangeView) Init() {
	go view.startUsersSearching()
}

func (view *fjiraAssignChangeView) Destroy() {

}

func (view *fjiraAssignChangeView) Draw(screen tcell.Screen) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Draw(screen)
	}
	view.topBar.Draw(screen)
	view.bottomBar.Draw(screen)
}

func (view *fjiraAssignChangeView) Update() {
	view.bottomBar.Update()
	if view.fuzzyFind != nil {
		view.fuzzyFind.Update()
	}
}

func (view *fjiraAssignChangeView) Resize(screenX, screenY int) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.Resize(screenX, screenY)
	}
	view.topBar.Resize(screenX, screenY)
	view.bottomBar.Resize(screenX, screenY)
}

func (view *fjiraAssignChangeView) HandleKeyEvent(ev *tcell.EventKey) {
	if view.fuzzyFind != nil {
		view.fuzzyFind.HandleKeyEvent(ev)
	}
}

func (view *fjiraAssignChangeView) startUsersSearching() {
	app.GetApp().ClearNow()
	app.GetApp().Loading(true)
	formatter, _ := GetFormatter()
	users := view.findUser(view.issue.Fields.Project.Id)
	usersStrings := formatter.formatJiraUsers(users)
	view.fuzzyFind = app.NewFuzzyFind(MessageUsersFuzzyFind, usersStrings)
	view.fuzzyFind.MarginBottom = 0
	app.GetApp().Loading(false)
	if user := <-view.fuzzyFind.Complete; true {
		app.GetApp().ClearNow()
		if user.Index < 0 {
			app.GetApp().SetView(NewIssueView(view.issue))
			return
		}
		view.fuzzyFind = nil
		view.assignUserToTicket(view.issue, &users[user.Index])
	}
}

func (view *fjiraAssignChangeView) findUser(project string) []jira.User {
	api, _ := GetApi()
	users, err := api.FindUsers(project)
	if err != nil {
		app.Error(err.Error())
	}
	return users
}

func (view *fjiraAssignChangeView) assignUserToTicket(issue *jira.Issue, user *jira.User) {
	if user == nil {
		app.GetApp().SetView(NewIssueView(view.issue))
		return
	}
	message := fmt.Sprintf(MessageChangingAssigneeTo, issue.Key, user.DisplayName)
	app.GetApp().ClearNow()
	view.bottomBar.AddItem(NewYesBarItem())
	view.bottomBar.AddItem(NewCancelBarItem())
	// TODO - should confirm be also drawable? at the moment yes/no are rendered out of the confirm thingy..
	userAssign := app.Confirm(app.GetApp(), message)
	switch userAssign {
	case true:
		view.doAssignmentChange(issue, user)
		goIntoIssueView(view.issue.Key)
	case false:
		app.GetApp().SetView(NewIssueView(view.issue))
	}
}

func (view fjiraAssignChangeView) doAssignmentChange(issue *jira.Issue, user *jira.User) {
	app.GetApp().LoadingWithText(true, MessageAssigningUser)
	api, _ := GetApi()
	err := api.DoAssignee(issue.Key, user.AccountId)
	app.GetApp().Loading(false)
	if err != nil {
		app.Error(fmt.Sprintf(MessageCannotAssignUser, user.DisplayName, issue.Key, err))
	}
	app.Success(fmt.Sprintf(MessageAssignSuccess, user.DisplayName, issue.Key))
}
