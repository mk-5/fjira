package statuses

import "github.com/mk-5/fjira/internal/jira"

func FormatJiraStatuses(statuses []jira.IssueStatus) []string {
	formatted := make([]string, 0, len(statuses))
	for _, status := range statuses {
		formatted = append(formatted, status.Name)
	}
	return formatted
}

func FormatJiraTransitions(statuses []jira.IssueTransition) []string {
	formatted := make([]string, 0, len(statuses))
	for _, status := range statuses {
		formatted = append(formatted, status.Name)
	}
	return formatted
}
