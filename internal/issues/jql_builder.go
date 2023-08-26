package issues

import (
	"fmt"
	"github.com/mk-5/fjira/internal/jira"
	"github.com/mk-5/fjira/internal/ui"
	"strings"
)

func BuildSearchIssuesJql(project *jira.Project, query string, status *jira.IssueStatus, user *jira.User, label string) string {
	jql := ""
	if project != nil && project.Id != ui.MessageAll {
		jql = jql + fmt.Sprintf("project=%s", project.Id)
	}
	orderBy := "ORDER BY status"
	query = strings.TrimSpace(query)
	if query != "" {
		jql = jql + fmt.Sprintf(" AND summary~\"%s*\"", query)
	}
	if status != nil && status.Name != ui.MessageAll {
		jql = jql + fmt.Sprintf(" AND status=%s", status.Id)
	}
	if user != nil && user.DisplayName != ui.MessageAll {
		jql = jql + fmt.Sprintf(" AND assignee=%s", user.AccountId)
	}
	// TODO - would be safer to check the index of inserted all message, instead of checking it like this / same for all All checks
	if label != "" && label != ui.MessageAll {
		jql = jql + fmt.Sprintf(" AND labels=%s", label)
	}
	if query != "" && issueRegExp.MatchString(query) {
		jql = jql + fmt.Sprintf(" OR issuekey=\"%s\"", query)
	}
	return fmt.Sprintf("%s %s", strings.TrimLeft(jql, " AND"), orderBy)
}
