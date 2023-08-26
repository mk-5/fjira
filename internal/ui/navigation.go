package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
)

// The whole navigation is placed here in `ui` package
// because of cyclic-imports - that are not allowed.

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
	ActionNew
	ActionDelete
)

var (
	BottomBarItemDefaultStyle  = app.DefaultStyle.Background(tcell.NewRGBColor(95, 135, 175)).Foreground(tcell.ColorWhite)
	BottomBarActionBarItemBold = app.DefaultStyle.Bold(true).Foreground(tcell.ColorDarkKhaki)
	BottomBarActionBarKeyBold  = BottomBarItemDefaultStyle.Foreground(tcell.NewRGBColor(21, 21, 21))
	TopBarItemDefault          = app.DefaultStyle.Background(tcell.NewRGBColor(95, 135, 95)).Foreground(tcell.ColorWhite)
	TopBarItemBold             = TopBarItemDefault.Foreground(app.AppBackground) // DarkOrange looks good here as well
)

type NavItemConfig struct {
	Text1  string
	Text2  string
	Action app.ActionBarAction
	Key    tcell.Key
	Rune   rune
}

func CreateBottomActionBar(text1 string, text2 string) *app.ActionBar {
	actionBar := app.NewActionBar(app.Bottom, app.Left)
	actionBar.AddItemWithStyles(
		text1,
		text2,
		BottomBarItemDefaultStyle,
		BottomBarActionBarItemBold,
	)
	return actionBar
}

func CreateTopActionBar(text1 string, text2 string) *app.ActionBar {
	actionBar := app.NewActionBar(app.Top, app.Right)
	actionBar.AddItemWithStyles(
		text1,
		text2,
		TopBarItemDefault, TopBarItemBold,
	)
	return actionBar
}

func CreateBottomActionBarWithItems(items []NavItemConfig) *app.ActionBar {
	actionBar := app.NewActionBar(app.Bottom, app.Left)
	for _, i := range items {
		actionBar.AddItem(&app.ActionBarItem{
			Id:          int(i.Action),
			Text1:       i.Text1,
			Text2:       i.Text2,
			Text1Style:  BottomBarItemDefaultStyle,
			Text2Style:  BottomBarActionBarKeyBold,
			TriggerKey:  i.Key,
			TriggerRune: i.Rune,
		})
	}
	return actionBar
}

func CreateTopActionBarWithItems(items []NavItemConfig) *app.ActionBar {
	actionBar := app.NewActionBar(app.Top, app.Left)
	for _, i := range items {
		actionBar.AddItem(&app.ActionBarItem{
			Id:          int(i.Action),
			Text1:       i.Text1,
			Text2:       i.Text2,
			Text1Style:  TopBarItemDefault,
			Text2Style:  TopBarItemBold,
			TriggerKey:  i.Key,
			TriggerRune: i.Rune,
		})
	}
	return actionBar
}

func CreateIssueTopBar(issue *jira.Issue) *app.ActionBar {
	items := []NavItemConfig{
		NavItemConfig{Text1: MessageIssueLabel, Text2: app.ActionBarLabel(issue.Key)},
		NavItemConfig{Text1: MessageLabelReporter, Text2: issue.Fields.Reporter.DisplayName},
		NavItemConfig{Text1: MessageLabelAssignee, Text2: issue.Fields.Assignee.DisplayName},
		NavItemConfig{Text1: MessageTypeStatus, Text2: issue.Fields.Type.Name},
		NavItemConfig{Text1: MessageLabelStatus, Text2: issue.Fields.Status.Name},
	}
	return CreateTopActionBarWithItems(items)
}

func CreateBottomLeftBar() *app.ActionBar {
	return app.NewActionBar(app.Bottom, app.Left)
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

func NewDeleteItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:         int(ActionDelete),
		Text1:      MessageDelete,
		Text2:      "[F2]",
		Text1Style: BottomBarItemDefaultStyle,
		Text2Style: BottomBarActionBarKeyBold,
		TriggerKey: tcell.KeyF2,
	}
}
