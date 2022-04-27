package fjira

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/mk5/fjira/internal/app"
	"github.com/mk5/fjira/internal/jira"
	"log"
)

type fjiraAssignChangeView struct {
	app.View
	bottomBar *app.ActionBar
	topBar    *app.ActionBar
	fuzzyFind *app.FuzzyFind
	issue     *jira.JiraIssue
}

const (
	MessageChangingAssigneeTo = "Changing %s assignee to %s [yn]: "
	MessageCannotAssignUser   = "Cannot assign user %s to ticket %s. Reason: %s"
	MessageAssignSuccess      = "User %s has been successfully assigned to issue %s."
	MessageUsersFuzzyFind     = "Select new assignee or ESC to cancel"
	MessageAssigningUser      = "Assigning user"
	Unassigned                = "Unassigned"
)

func NewAssignChangeView(issue *jira.JiraIssue) *fjiraAssignChangeView {
	return &fjiraAssignChangeView{
		issue:     issue,
		topBar:    CreateNewIssueTopBar(issue),
		bottomBar: CreateNewIssueBottomBar(issue),
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
	app.GetApp().Loading(false)
	select {
	case user := <-view.fuzzyFind.Complete:
		app.GetApp().ClearNow()
		if user.Index < 0 {
			app.GetApp().SetView(NewIssueView(view.issue))
			return
		}
		view.fuzzyFind = nil
		view.assignUserToTicket(view.issue, &users[user.Index])
	}
}

func (view *fjiraAssignChangeView) findUser(project string) []jira.JiraUser {
	api, _ := GetApi()
	users, err := api.FindUsers(project)
	if err != nil {
		log.Fatalln(err)
	}
	return users
}

func (view *fjiraAssignChangeView) assignUserToTicket(issue *jira.JiraIssue, user *jira.JiraUser) {
	if user == nil {
		app.GetApp().SetView(NewIssueView(view.issue))
		return
	}
	message := fmt.Sprintf(MessageChangingAssigneeTo, issue.Key, user.DisplayName)
	app.GetApp().ClearNow()
	view.bottomBar.AddItem(NewNewAssigneeBarItem(user))
	view.bottomBar.AddItem(NewYesBarItem())
	view.bottomBar.AddItem(NewCancelBarItem())
	userAssign := app.Confirm(app.GetApp(), message)
	switch userAssign {
	case true:
		view.doAssignmentChange(issue, user)
		goIntoIssueViewFetchIssue(view.issue.Key)
		break
	case false:
		app.GetApp().SetView(NewIssueView(view.issue))
	}
}

func (view fjiraAssignChangeView) doAssignmentChange(issue *jira.JiraIssue, user *jira.JiraUser) {
	app.GetApp().LoadingWithText(true, MessageAssigningUser)
	api, _ := GetApi()
	err := api.DoAssignee(issue.Key, user.AccountId)
	app.GetApp().Loading(false)
	if err != nil {
		log.Fatalln(fmt.Sprintf(MessageCannotAssignUser, user.DisplayName, issue.Key, err))
	}
	//fmt.Println(app.EmptyLine)
	//fmt.Println(color.GreenString(MessageAssignSuccess, user.DisplayName, issue.Key))
	//b := make([]byte, 1)
	//os.Stdin.Read(b)
}
