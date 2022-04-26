package fjira

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/mk5/fjira/internal/app"
	"github.com/mk5/fjira/internal/jira"
)

const (
	MessageProjectLabel = "Project: "
	MessageIssueLabel   = "Issue: "
)

const (
	ActionAssigneeChange app.ActionBarAction = iota
	ActionStatusChange
	ActionEscape
	ActionYes
	ActionNo
)

var (
	BottomBarActionBarItemBold = tcell.StyleDefault.Bold(true).Foreground(tcell.ColorDarkKhaki)
	BottomBarActionBarKeyBold  = tcell.StyleDefault.Bold(true).Foreground(tcell.ColorDarkCyan).Underline(true)
)

// TODO - concrete "bottomBar" struct with helper methods like "setCurrentProject(..)"
func CreateNewEmptyProjectBottomBar() *app.ActionBar {
	actionBar := app.NewActionBar(app.Bottom, app.Left)
	actionBar.AddItemWithStyles(
		MessageProjectLabel,
		app.ActionBarLabel(""),
		tcell.StyleDefault, BottomBarActionBarItemBold,
	)
	return actionBar
}

func CreateNewIssueBottomBar(issue *jira.JiraIssue) *app.ActionBar {
	actionBar := app.NewActionBar(app.Bottom, app.Left)
	actionBar.AddItemWithStyles(
		MessageIssueLabel,
		app.ActionBarLabel(issue.Key),
		tcell.StyleDefault, BottomBarActionBarItemBold,
	)
	return actionBar
}

func CreateNewProjectBottomBar(project *jira.JiraProject) *app.ActionBar {
	actionBar := app.NewActionBar(app.Bottom, app.Left)
	actionBar.AddItemWithStyles(
		MessageProjectLabel,
		app.ActionBarLabel(fmt.Sprintf("[%s]%s", project.Key, project.Name)),
		tcell.StyleDefault, BottomBarActionBarItemBold,
	)
	actionBar.AddItem(NewByStatusBarItem())
	return actionBar
}

func NewCancelBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:         int(ActionEscape),
		Text1:      "ESC",
		Text2:      " - cancel",
		Text1Style: BottomBarActionBarKeyBold,
		Text2Style: tcell.StyleDefault,
		TriggerKey: tcell.KeyEscape,
	}
}

func NewStatusChangeBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:          int(ActionStatusChange),
		Text1:       "s",
		Text2:       " - change status",
		Text1Style:  BottomBarActionBarKeyBold,
		Text2Style:  tcell.StyleDefault,
		TriggerRune: 's',
	}
}

func NewByStatusBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:          int(ActionStatusChange),
		Text1:       "s",
		Text2:       " - by status",
		Text1Style:  BottomBarActionBarKeyBold,
		Text2Style:  tcell.StyleDefault,
		TriggerRune: 's',
	}
}

func NewAssigneeChangeBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:          int(ActionAssigneeChange),
		Text1:       "a",
		Text2:       " - assign user",
		Text1Style:  BottomBarActionBarKeyBold,
		Text2Style:  tcell.StyleDefault,
		TriggerRune: 'a',
	}
}

func NewNewStatusBarItem(newStatus string) *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:         -1,
		Text1:      "New status: ",
		Text2:      newStatus,
		Text1Style: tcell.StyleDefault,
		Text2Style: BottomBarActionBarKeyBold,
	}
}

func NewNewAssigneeBarItem(newAssignee *jira.JiraUser) *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:         -1,
		Text1:      "New assignee: ",
		Text2:      fmt.Sprintf("%s <%s>", newAssignee.DisplayName, newAssignee.EmailAddress),
		Text1Style: tcell.StyleDefault,
		Text2Style: BottomBarActionBarKeyBold,
	}
}

func NewYesBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:          int(ActionYes),
		Text1:       "y",
		Text2:       " - yes",
		Text1Style:  BottomBarActionBarKeyBold,
		Text2Style:  tcell.StyleDefault,
		TriggerRune: 'y',
	}
}
