package fjira

import (
	"github.com/gdamore/tcell"
	"github.com/mk5/fjira/internal/app"
	"github.com/mk5/fjira/internal/jira"
)

const (
	MessageLabelIssue    = "Issue: "
	MessageLabelStatus   = "Status: "
	MessageTypeStatus    = "Type: "
	MessageLabelAssignee = "Assignee: "
	MessageLabelReporter = "Reporter: "
)

var (
	IssueBarActionBarItemBold = tcell.StyleDefault.Bold(true).Foreground(tcell.ColorDarkKhaki)
)

func CreateNewIssueTopBar(issue *jira.JiraIssue) *app.ActionBar {
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
