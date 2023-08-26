package users

import (
	"fmt"
	"github.com/mk-5/fjira/internal/jira"
)

func FormatJiraUser(user *jira.User) string {
	return fmt.Sprintf("%s <%s>", user.DisplayName, user.EmailAddress)
}

func FormatJiraUsers(users []jira.User) []string {
	formatted := make([]string, 0, len(users))
	for _, user := range users {
		formatted = append(formatted, FormatJiraUser(&user))
	}
	return formatted
}
