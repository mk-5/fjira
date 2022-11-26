package fjira

import (
	"fmt"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	"strconv"
	"strings"
)

// TODO - change to static?
type fjiraFormatter interface {
	formatJiraProject(project *jira.Project) string
	formatJiraProjects(projects []jira.Project) []string
	formatJiraIssue(issue *jira.Issue) string
	formatJiraIssues(issues []jira.Issue) []string
	formatJiraUser(user *jira.User) string
	formatJiraUsers(user []jira.User) []string
	formatJiraBoards(boards []*jira.BoardItem) []string
	formatJiraTransitions(transitions []jira.IssueTransition) []string
	formatJiraStatuses(statuses []jira.IssueStatus) []string
}

type defaultFormatter struct{}

const (
	TableColumnPadding = 2
	MaxSummaryColWidth = 45
	MaxStatusColWidth  = 12
)

func (f *defaultFormatter) formatJiraIssue(issue *jira.Issue) string {
	return fmt.Sprintf("%s %s [%s] - %s",
		issue.Key,
		issue.Fields.Summary,
		issue.Fields.Status.Name,
		f.formatAssignee(issue))
}

func (*defaultFormatter) formatJiraIssueTable(issue *jira.Issue, summaryColWidth int, statusColWidth int) string {
	assignee := issue.Fields.Assignee.DisplayName
	if assignee == "" {
		assignee = MessageUnassigned
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

func (f *defaultFormatter) formatJiraIssues(issues []jira.Issue) []string {
	formatted := make([]string, 0, len(issues))
	summaryColWidth := f.findIssueColumnSize(&issues, func(i jira.Issue) string {
		return i.Fields.Summary
	})
	statusColWidth := f.findIssueColumnSize(&issues, func(i jira.Issue) string {
		return i.Fields.Status.Name
	})
	for _, issue := range issues {
		formatted = append(formatted, f.formatJiraIssueTable(&issue, summaryColWidth, statusColWidth))
	}
	return formatted
}

func (*defaultFormatter) formatJiraUser(user *jira.User) string {
	return fmt.Sprintf("%s <%s>", user.DisplayName, user.EmailAddress)
}

func (f *defaultFormatter) formatJiraUsers(users []jira.User) []string {
	formatted := make([]string, 0, len(users))
	for _, user := range users {
		formatted = append(formatted, f.formatJiraUser(&user))
	}
	return formatted
}

func (*defaultFormatter) formatJiraProject(project *jira.Project) string {
	return fmt.Sprintf("[%s] %s", project.Key, project.Name)
}

func (f *defaultFormatter) formatJiraProjects(projects []jira.Project) []string {
	formatted := make([]string, 0, len(projects))
	for _, project := range projects {
		formatted = append(formatted, f.formatJiraProject(&project))
	}
	return formatted
}

func (f *defaultFormatter) formatJiraTransitions(statuses []jira.IssueTransition) []string {
	formatted := make([]string, 0, len(statuses))
	for _, status := range statuses {
		formatted = append(formatted, status.Name)
	}
	return formatted
}

func (f *defaultFormatter) formatJiraStatuses(statuses []jira.IssueStatus) []string {
	formatted := make([]string, 0, len(statuses))
	for _, status := range statuses {
		formatted = append(formatted, status.Name)
	}
	return formatted
}

func (f *defaultFormatter) formatJiraBoards(boards []*jira.BoardItem) []string {
	formatted := make([]string, 0, len(boards))
	for _, board := range boards {
		formatted = append(formatted, board.Name)
	}
	return formatted
}

func (f *defaultFormatter) findIssueColumnSize(items *[]jira.Issue, colSupplier func(issue jira.Issue) string) int {
	max := 0
	for _, item := range *items {
		current := colSupplier(item)
		if max < len(current) {
			max = len(current)
		}
	}
	return max
}

func (f *defaultFormatter) formatAssignee(issue *jira.Issue) string {
	assignee := issue.Fields.Assignee.DisplayName
	if assignee == "" {
		assignee = MessageUnassigned
	}
	return assignee
}
