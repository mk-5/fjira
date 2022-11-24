package fjira

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
)

const (
	ActionAssigneeChange app.ActionBarAction = iota
	ActionStatusChange
	ActionSearchByStatus
	ActionSearchByAssignee
	ActionSearchByLabel
	ActionBoards
	ActionComment
	ActionCancel
	ActionOpen
	ActionYes
	ActionAddLabel
	ActionSelect
	ActionUnselect
)

var (
	BottomBarItemDefaultStyle  = app.DefaultStyle.Background(tcell.NewRGBColor(95, 135, 175)).Foreground(tcell.ColorWhite)
	BottomBarActionBarItemBold = app.DefaultStyle.Bold(true).Foreground(tcell.ColorDarkKhaki)
	BottomBarActionBarKeyBold  = BottomBarItemDefaultStyle.Foreground(tcell.NewRGBColor(21, 21, 21))
	TopBarItemDefault          = app.DefaultStyle.Background(tcell.NewRGBColor(95, 135, 95)).Foreground(tcell.ColorWhite)
	TopBarItemBold             = TopBarItemDefault.Foreground(app.AppBackground) // DarkOrange looks good here as well
)

func CreateProjectBottomBar() *app.ActionBar {
	actionBar := app.NewActionBar(app.Bottom, app.Left)
	actionBar.AddItemWithStyles(
		MessageProjectLabel,
		app.ActionBarLabel(""),
		BottomBarItemDefaultStyle, BottomBarActionBarItemBold,
	)
	return actionBar
}

func CreateProjectsTopBar() *app.ActionBar {
	actionBar := app.NewActionBar(app.Top, app.Right)
	actionBar.AddItemWithStyles(
		MessageProjectLabel,
		app.ActionBarLabel(""),
		TopBarItemDefault, TopBarItemBold,
	)
	return actionBar
}

func CreateSearchIssuesBottomBar() *app.ActionBar {
	actionBar := app.NewActionBar(app.Bottom, app.Left)
	actionBar.AddItem(NewByStatusBarItem())
	actionBar.AddItem(NewByAssigneeBarItem())
	actionBar.AddItem(NewByLabelBarItem())
	actionBar.AddItem(NewBoardsViewItem())
	return actionBar
}

func CreateSearchIssuesTopBar(project *jira.Project) *app.ActionBar {
	actionBar := app.NewActionBar(app.Top, app.Left)
	actionBar.AddItemWithStyles(
		MessageProjectLabel,
		app.ActionBarLabel(fmt.Sprintf("[%s]%s", project.Key, project.Name)),
		TopBarItemDefault, TopBarItemBold,
	)
	actionBar.AddItemWithStyles(
		MessageLabelStatus,
		MessageAll,
		TopBarItemDefault, TopBarItemBold,
	)
	actionBar.AddItemWithStyles(
		MessageLabelAssignee,
		MessageAll,
		TopBarItemDefault, TopBarItemBold,
	)
	actionBar.AddItemWithStyles(
		MessageLabelLabel,
		MessageAll,
		TopBarItemDefault, TopBarItemBold,
	)
	return actionBar
}

func CreateBottomLeftBar() *app.ActionBar {
	actionBar := app.NewActionBar(app.Bottom, app.Left)
	return actionBar
}

func CreateIssueTopBar(issue *jira.Issue) *app.ActionBar {
	actionBar := app.NewActionBar(app.Top, app.Left)
	actionBar.AddItemWithStyles(
		MessageIssueLabel,
		app.ActionBarLabel(issue.Key),
		TopBarItemDefault,
		TopBarItemBold,
	)
	actionBar.AddItemWithStyles(
		MessageLabelReporter,
		issue.Fields.Reporter.DisplayName,
		TopBarItemDefault,
		TopBarItemBold,
	)
	actionBar.AddItemWithStyles(
		MessageLabelAssignee,
		issue.Fields.Assignee.DisplayName,
		TopBarItemDefault,
		TopBarItemBold,
	)
	actionBar.AddItemWithStyles(
		MessageTypeStatus,
		issue.Fields.Type.Name,
		TopBarItemDefault,
		TopBarItemBold,
	)
	actionBar.AddItemWithStyles(
		MessageLabelStatus,
		issue.Fields.Status.Name,
		TopBarItemDefault,
		TopBarItemBold,
	)
	return actionBar
}

func NewCancelBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:         int(ActionCancel),
		Text1:      "Cancel ",
		Text2:      "[ESC]",
		Text1Style: BottomBarItemDefaultStyle,
		Text2Style: BottomBarActionBarKeyBold,
		TriggerKey: tcell.KeyEscape,
	}
}

func NewStatusChangeBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:          int(ActionStatusChange),
		Text1:       MessageChangeStatus,
		Text2:       "[s]",
		Text1Style:  BottomBarItemDefaultStyle,
		Text2Style:  BottomBarActionBarKeyBold,
		TriggerKey:  tcell.KeyF1,
		TriggerRune: 's',
	}
}

func NewByStatusBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:         int(ActionSearchByStatus),
		Text1:      MessageByStatus,
		Text2:      "[F1]",
		Text1Style: BottomBarItemDefaultStyle,
		Text2Style: BottomBarActionBarKeyBold,
		TriggerKey: tcell.KeyF1,
	}
}

func NewByAssigneeBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:         int(ActionSearchByAssignee),
		Text1:      MessageByAssignee,
		Text2:      "[F2]",
		Text1Style: BottomBarItemDefaultStyle,
		Text2Style: BottomBarActionBarKeyBold,
		TriggerKey: tcell.KeyF2,
	}
}

func NewByLabelBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:         int(ActionSearchByLabel),
		Text1:      MessageByLabel,
		Text2:      "[F3]",
		Text1Style: BottomBarItemDefaultStyle,
		Text2Style: BottomBarActionBarKeyBold,
		TriggerKey: tcell.KeyF3,
	}
}

func NewBoardsViewItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:         int(ActionBoards),
		Text1:      MessageBoards,
		Text2:      "[F4]",
		Text1Style: BottomBarItemDefaultStyle,
		Text2Style: BottomBarActionBarKeyBold,
		TriggerKey: tcell.KeyF4,
	}
}

func NewAssigneeChangeBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:          int(ActionAssigneeChange),
		Text1:       MessageAssignUser,
		Text2:       "[a]",
		Text1Style:  BottomBarItemDefaultStyle,
		Text2Style:  BottomBarActionBarKeyBold,
		TriggerKey:  -1,
		TriggerRune: 'a',
	}
}

func CreateCommentBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:          int(ActionComment),
		Text1:       MessageComment,
		Text2:       "[c]",
		Text1Style:  BottomBarItemDefaultStyle,
		Text2Style:  BottomBarActionBarKeyBold,
		TriggerKey:  -1,
		TriggerRune: 'c',
	}
}

func CreateAddLabelBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:          int(ActionAddLabel),
		Text1:       MessageLabel,
		Text2:       "[l]",
		Text1Style:  BottomBarItemDefaultStyle,
		Text2Style:  BottomBarActionBarKeyBold,
		TriggerKey:  tcell.KeyF2,
		TriggerRune: 'l',
	}
}

func CreateScrollBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Text1:       MessageScroll,
		Text2:       "[↑↓]",
		Text1Style:  BottomBarItemDefaultStyle,
		Text2Style:  BottomBarActionBarKeyBold,
		TriggerKey:  -1,
		TriggerRune: -1,
	}
}

func CreateArrowsNavigateItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Text1:       MessageNavigate,
		Text2:       "[←→↑↓]",
		Text1Style:  BottomBarItemDefaultStyle,
		Text2Style:  BottomBarActionBarKeyBold,
		TriggerKey:  -1,
		TriggerRune: -1,
	}
}

func CreateMoveArrowsItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Text1:       MessageMoveIssue,
		Text2:       "[←→]",
		Text1Style:  BottomBarItemDefaultStyle,
		Text2Style:  BottomBarActionBarKeyBold,
		TriggerKey:  -1,
		TriggerRune: -1,
	}
}

func CreateSelectItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:          int(ActionSelect),
		Text1:       MessageSelect,
		Text2:       "[enter]",
		Text1Style:  BottomBarItemDefaultStyle,
		Text2Style:  BottomBarActionBarKeyBold,
		TriggerKey:  tcell.KeyEnter,
		TriggerRune: -1,
	}
}

func CreateUnSelectItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:          int(ActionUnselect),
		Text1:       MessageUnselect,
		Text2:       "[enter]",
		Text1Style:  BottomBarItemDefaultStyle,
		Text2Style:  BottomBarActionBarKeyBold,
		TriggerKey:  tcell.KeyEnter,
		TriggerRune: -1,
	}
}

func NewYesBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:          int(ActionYes),
		Text1:       MessageYes,
		Text2:       "[y]",
		Text1Style:  BottomBarItemDefaultStyle,
		Text2Style:  BottomBarActionBarKeyBold,
		TriggerRune: 'y',
	}
}

func NewOpenBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:          int(ActionOpen),
		Text1:       MessageOpen,
		Text2:       "[o]",
		Text1Style:  BottomBarItemDefaultStyle,
		Text2Style:  BottomBarActionBarKeyBold,
		TriggerRune: 'o',
	}
}

func NewSaveBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:         int(ActionYes),
		Text1:      MessageSave,
		Text2:      "[F1]",
		Text1Style: BottomBarItemDefaultStyle,
		Text2Style: BottomBarActionBarKeyBold,
		TriggerKey: tcell.KeyF1,
	}
}
