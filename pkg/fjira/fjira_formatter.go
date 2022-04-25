package fjira

import (
	"fmt"
	"github.com/mk5/fjira/internal/app"
	"github.com/mk5/fjira/internal/jira"
	"strconv"
	"strings"
)

type FjiraFormatter interface {
	formatJiraProject(project *jira.JiraProject) string
	formatJiraProjects(projects []jira.JiraProject) []string
	formatJiraIssue(issue *jira.JiraIssue) string
	formatJiraIssues(issues []jira.JiraIssue) []string
	formatJiraUser(user *jira.JiraUser) string
	formatJiraUsers(user []jira.JiraUser) []string
	formatJiraStatuses(statuses []jira.JiraIssueTransition) []string
}

type defaultFormatter struct{}

const (
	TableColumnPadding = 2
	MaxSummaryColWidth = 45
	MaxStatusColWidth  = 12
)

func (f *defaultFormatter) formatJiraIssue(issue *jira.JiraIssue) string {
	return fmt.Sprintf("%s %s [%s] - %s",
		issue.Key,
		issue.Fields.Summary,
		issue.Fields.Status.Name,
		f.formatAssignee(issue))
}

func (*defaultFormatter) formatJiraIssueTable(issue *jira.JiraIssue, summaryColWidth int, statusColWidth int) string {
	assignee := issue.Fields.Assignee.DisplayName
	if assignee == "" {
		assignee = Unassigned
	}
	summaryColWidth = app.MinInt(summaryColWidth, MaxSummaryColWidth)
	summaryCut := app.MinInt(summaryColWidth, len(issue.Fields.Summary))
	statusColWidth = app.MinInt(statusColWidth, MaxStatusColWidth)
	statusCut := app.MinInt(statusColWidth, len(issue.Fields.Status.Name))
	return fmt.Sprintf("%10s %"+strconv.Itoa(summaryColWidth+TableColumnPadding)+"s %"+strconv.Itoa(statusColWidth+4+TableColumnPadding)+"s %s",
		issue.Key,
		issue.Fields.Summary[:summaryCut],
		fmt.Sprintf("[%s]", strings.ToUpper(issue.Fields.Status.Name[:statusCut])),
		fmt.Sprintf("- %s", assignee))
}

func (f *defaultFormatter) formatJiraIssues(issues []jira.JiraIssue) []string {
	formatted := make([]string, 0, len(issues))
	summaryColWidth := f.findIssueColumnSize(&issues, func(i jira.JiraIssue) string {
		return i.Fields.Summary
	})
	statusColWidth := f.findIssueColumnSize(&issues, func(i jira.JiraIssue) string {
		return i.Fields.Status.Name
	})
	for _, issue := range issues {
		formatted = append(formatted, f.formatJiraIssueTable(&issue, summaryColWidth, statusColWidth))
	}
	return formatted
}

func (*defaultFormatter) formatJiraUser(user *jira.JiraUser) string {
	return fmt.Sprintf("%s <%s>", user.DisplayName, user.EmailAddress)
}

func (f *defaultFormatter) formatJiraUsers(users []jira.JiraUser) []string {
	formatted := make([]string, 0, len(users))
	for _, user := range users {
		formatted = append(formatted, f.formatJiraUser(&user))
	}
	return formatted
}

func (*defaultFormatter) formatJiraProject(project *jira.JiraProject) string {
	return fmt.Sprintf("[%s] %s", project.Key, project.Name)
}

func (f *defaultFormatter) formatJiraProjects(projects []jira.JiraProject) []string {
	formatted := make([]string, 0, len(projects))
	for _, project := range projects {
		formatted = append(formatted, f.formatJiraProject(&project))
	}
	return formatted
}

func (f *defaultFormatter) formatJiraStatuses(statuses []jira.JiraIssueTransition) []string {
	formatted := make([]string, 0, len(statuses))
	for _, status := range statuses {
		formatted = append(formatted, fmt.Sprintf("%s", status.Name))
	}
	return formatted
}

func (f *defaultFormatter) findIssueColumnSize(items *[]jira.JiraIssue, colSupplier func(issue jira.JiraIssue) string) int {
	max := 0
	for _, item := range *items {
		current := colSupplier(item)
		if max < len(current) {
			max = len(current)
		}
	}
	return max
}

func (f *defaultFormatter) formatAssignee(issue *jira.JiraIssue) string {
	assignee := issue.Fields.Assignee.DisplayName
	if assignee == "" {
		assignee = Unassigned
	}
	return assignee
}
