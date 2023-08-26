package issues

import (
	"fmt"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	"github.com/mk-5/fjira/internal/ui"
	"strconv"
	"strings"
)

func FormatJiraIssue(issue *jira.Issue) string {
	return fmt.Sprintf("%s %s [%s] - %s",
		issue.Key,
		issue.Fields.Summary,
		issue.Fields.Status.Name,
		FormatAssignee(issue))
}

func FormatJiraIssueTable(issue *jira.Issue, summaryColWidth int, statusColWidth int) string {
	assignee := issue.Fields.Assignee.DisplayName
	if assignee == "" {
		assignee = ui.MessageUnassigned
	}
	summaryColWidth = app.MinInt(summaryColWidth, ui.MaxSummaryColWidth)
	summaryCut := app.MinInt(summaryColWidth, len(issue.Fields.Summary))
	statusColWidth = app.MinInt(statusColWidth, ui.MaxStatusColWidth)
	statusCut := app.MinInt(statusColWidth, len(issue.Fields.Status.Name))
	return fmt.Sprintf("%10s %"+strconv.Itoa(summaryColWidth+ui.TableColumnPadding)+"s %"+strconv.Itoa(statusColWidth+4+ui.TableColumnPadding)+"s %s",
		issue.Key,
		issue.Fields.Summary[:summaryCut],
		fmt.Sprintf("[%s]", strings.ToUpper(issue.Fields.Status.Name[:statusCut])),
		fmt.Sprintf("- %s", assignee))
}

func FormatJiraIssues(issues []jira.Issue) []string {
	formatted := make([]string, 0, len(issues))
	summaryColWidth := findIssueColumnSize(&issues, func(i jira.Issue) string {
		return i.Fields.Summary
	})
	statusColWidth := findIssueColumnSize(&issues, func(i jira.Issue) string {
		return i.Fields.Status.Name
	})
	for _, issue := range issues {
		formatted = append(formatted, FormatJiraIssueTable(&issue, summaryColWidth, statusColWidth))
	}
	return formatted
}

func FormatAssignee(issue *jira.Issue) string {
	assignee := issue.Fields.Assignee.DisplayName
	if assignee == "" {
		assignee = ui.MessageUnassigned
	}
	return assignee
}

func findIssueColumnSize(items *[]jira.Issue, colSupplier func(issue jira.Issue) string) int {
	max := 0
	for _, item := range *items {
		current := colSupplier(item)
		if max < len(current) {
			max = len(current)
		}
	}
	return max
}
