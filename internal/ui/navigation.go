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
	ActionCopyIssue
	ActionCancel
	ActionOpen
	ActionYes
	ActionAddLabel
	ActionSelect
	ActionUnselect
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
		bottomBarItemDefaultStyle(),
		bottomBarActionBarItemBold(),
	)
	return actionBar
}

func CreateTopActionBar(text1 string, text2 string) *app.ActionBar {
	actionBar := app.NewActionBar(app.Top, app.Right)
	actionBar.AddItemWithStyles(
		text1,
		text2,
		topBarItemDefault(), topBarItemBold(),
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
			Text1Style:  bottomBarItemDefaultStyle(),
			Text2Style:  bottomBarActionBarKeyBold(),
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
			Text1Style:  topBarItemDefault(),
			Text2Style:  topBarItemBold(),
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
		Text1Style: bottomBarItemDefaultStyle(),
		Text2Style: bottomBarActionBarKeyBold(),
		TriggerKey: tcell.KeyEscape,
	}
}

func CreateScrollBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Text1:       MessageScroll,
		Text2:       "[↑↓]",
		Text1Style:  bottomBarItemDefaultStyle(),
		Text2Style:  bottomBarActionBarKeyBold(),
		TriggerKey:  -1,
		TriggerRune: -1,
	}
}

func CreateArrowsNavigateItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Text1:       MessageNavigate,
		Text2:       "[←→↑↓]",
		Text1Style:  bottomBarItemDefaultStyle(),
		Text2Style:  bottomBarActionBarKeyBold(),
		TriggerKey:  -1,
		TriggerRune: -1,
	}
}

func CreateMoveArrowsItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Text1:       MessageMoveIssue,
		Text2:       "[←→]",
		Text1Style:  bottomBarItemDefaultStyle(),
		Text2Style:  bottomBarActionBarKeyBold(),
		TriggerKey:  -1,
		TriggerRune: -1,
	}
}

func CreateSelectItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:          int(ActionSelect),
		Text1:       MessageSelect,
		Text2:       "[enter]",
		Text1Style:  bottomBarItemDefaultStyle(),
		Text2Style:  bottomBarActionBarKeyBold(),
		TriggerKey:  tcell.KeyEnter,
		TriggerRune: -1,
	}
}

func CreateUnSelectItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:          int(ActionUnselect),
		Text1:       MessageUnselect,
		Text2:       "[enter]",
		Text1Style:  bottomBarItemDefaultStyle(),
		Text2Style:  bottomBarActionBarKeyBold(),
		TriggerKey:  tcell.KeyEnter,
		TriggerRune: -1,
	}
}

func NewYesBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:          int(ActionYes),
		Text1:       MessageYes,
		Text2:       "[y]",
		Text1Style:  bottomBarItemDefaultStyle(),
		Text2Style:  bottomBarActionBarKeyBold(),
		TriggerRune: 'y',
	}
}

func NewOpenBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:          int(ActionOpen),
		Text1:       MessageOpen,
		Text2:       "[o]",
		Text1Style:  bottomBarItemDefaultStyle(),
		Text2Style:  bottomBarActionBarKeyBold(),
		TriggerRune: 'o',
	}
}

func NewSaveBarItem() *app.ActionBarItem {
	return &app.ActionBarItem{
		Id:         int(ActionYes),
		Text1:      MessageSave,
		Text2:      "[F1]",
		Text1Style: bottomBarItemDefaultStyle(),
		Text2Style: bottomBarActionBarKeyBold(),
		TriggerKey: tcell.KeyF1,
	}
}

func bottomBarItemDefaultStyle() tcell.Style {
	return app.DefaultStyle().Background(app.Color("navigation.bottom.background")).Foreground(app.Color("navigation.bottom.foreground1"))
}

func bottomBarActionBarItemBold() tcell.Style {
	return app.DefaultStyle().Bold(true).Foreground(app.Color("navigation.bottom.foreground2"))
}

func bottomBarActionBarKeyBold() tcell.Style {
	return bottomBarItemDefaultStyle().Foreground(app.Color("navigation.bottom.foreground2"))
}

func topBarItemDefault() tcell.Style {
	return app.DefaultStyle().Background(app.Color("navigation.top.background")).Foreground(app.Color("navigation.top.foreground1"))
}

func topBarItemBold() tcell.Style {
	return topBarItemDefault().Foreground(app.Color("navigation.top.foreground2")) // DarkOrange looks good here as well
}
