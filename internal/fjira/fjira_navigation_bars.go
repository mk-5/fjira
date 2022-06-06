package fjira

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mk5/fjira/internal/app"
	"github.com/mk5/fjira/internal/jira"
)

const (
	ActionAssigneeChange app.ActionBarAction = iota
	ActionStatusChange
	ActionComment
	ActionEscape
	ActionOpen
	ActionYes
)

var (
	BottomBarActionBarItemBold = tcell.StyleDefault.Bold(true).Foreground(tcell.ColorDarkKhaki)
	BottomBarActionBarKeyBold  = tcell.StyleDefault.Bold(true).Foreground(tcell.ColorDarkCyan).Underline(true)
	TopBarItemBold             = tcell.StyleDefault.Bold(true).Foreground(tcell.ColorDarkKhaki)
	IssueBarActionBarItemBold  = tcell.StyleDefault.Bold(true).Foreground(tcell.ColorDarkKhaki)
)

func CreateProjectBottomBar() *app.ActionBar {
	actionBar := app.NewActionBar(app.Bottom, app.Left)
	actionBar.AddItemWithStyles(
		MessageProjectLabel,
		app.ActionBarLabel(""),
		tcell.StyleDefault, BottomBarActionBarItemBold,
	)
	return actionBar
}

func CreateIssueBottomBar(issue *jira.JiraIssue) *app.ActionBar {
	actionBar := app.NewActionBar(app.Bottom, app.Left)
	actionBar.AddItemWithStyles(
		MessageIssueLabel,
		app.ActionBarLabel(issue.Key),
		tcell.StyleDefault, BottomBarActionBarItemBold,
	)
	return actionBar
}

func CreateSearchIssuesBottomBar(project *jira.JiraProject) *app.ActionBar {
	actionBar := app.NewActionBar(app.Bottom, app.Left)
	actionBar.AddItemWithStyles(
		MessageProjectLabel,
		app.ActionBarLabel(fmt.Sprintf("[%s]%s", project.Key, project.Name)),
		tcell.StyleDefault, BottomBarActionBarItemBold,
	)
	actionBar.AddItem(NewByStatusBarItem())
	actionBar.AddItem(NewByAssigneeBarItem())
	return actionBar
}

func CreateSearchIssuesTopBar() *app.ActionBar {
	actionBar := app.NewActionBar(app.Top, app.Right)
	actionBar.AddItemWithStyles(
		"Status: ",
		MessageAll,
		tcell.StyleDefault, TopBarItemBold,
	)
	actionBar.AddItemWithStyles(
		"Assignee: ",
		MessageAll,
		tcell.StyleDefault, TopBarItemBold,
	)
	return actionBar
}

func CreateIssueTopBar(issue *jira.JiraIssue) *app.ActionBar {
	actionBar := app.NewActionBar(app.Top, app.Right)
	actionBar.AddItemWithStyles(
		MessageLabelReporter,
		issue.Fields.Reporter.DisplayName,
		tcell.StyleDefault,
		IssueBarActionBarItemBold,
	)
	actionBar.AddItemWithStyles(
		MessageLabelAssignee,
		issue.Fields.Assignee.DisplayName,
		tcell.StyleDefault,
		IssueBarActionBarItemBold,
	)
	actionBar.AddItemWithStyles(
		MessageTypeStatus,
		issue.Fields.Type.Name,
		tcell.StyleDefault,
		IssueBarActionBarItemBold,
	)
	actionBar.AddItemWithStyles(
		MessageLabelStatus,
		issue.Fields.Status.Name,
		tcell.StyleDefault,
		IssueBarActionBarItemBold,
	)
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
		TriggerKey:  tcell.KeyF1,
		TriggerRune: 's',
	}
}

func NewByStatusBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:         int(ActionStatusChange),
		Text1:      "F1",
		Text2:      " - by status",
		Text1Style: BottomBarActionBarKeyBold,
		Text2Style: tcell.StyleDefault,
		TriggerKey: tcell.KeyF1,
	}
}

func NewByAssigneeBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:         int(ActionAssigneeChange),
		Text1:      "F2",
		Text2:      " - by assignee",
		Text1Style: BottomBarActionBarKeyBold,
		Text2Style: tcell.StyleDefault,
		TriggerKey: tcell.KeyF2,
	}
}

func NewAssigneeChangeBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:          int(ActionAssigneeChange),
		Text1:       "a",
		Text2:       " - assign user",
		Text1Style:  BottomBarActionBarKeyBold,
		Text2Style:  tcell.StyleDefault,
		TriggerKey:  tcell.KeyF2,
		TriggerRune: 'a',
	}
}

func CreateCommentBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:          int(ActionComment),
		Text1:       "c",
		Text2:       " - comment",
		Text1Style:  BottomBarActionBarKeyBold,
		Text2Style:  tcell.StyleDefault,
		TriggerKey:  tcell.KeyF2,
		TriggerRune: 'c',
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
		Text2:      fmt.Sprintf("%s", newAssignee.DisplayName),
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

func NewOpenBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:          int(ActionOpen),
		Text1:       "o",
		Text2:       " - open",
		Text1Style:  BottomBarActionBarKeyBold,
		Text2Style:  tcell.StyleDefault,
		TriggerRune: 'o',
	}
}

func NewSaveBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:         int(ActionYes),
		Text1:      "F1",
		Text2:      " - save",
		Text1Style: BottomBarActionBarKeyBold,
		Text2Style: tcell.StyleDefault,
		TriggerKey: tcell.KeyF1,
	}
}
