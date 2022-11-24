package fjira

import (
	"github.com/mk-5/fjira/internal/jira"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_shouldFormatProjects(t *testing.T) {
	assertions := assert.New(t)
	// given
	formatter := defaultFormatter{}
	issues := []jira.Issue{
		{
			Key: "TEST-123",
			Fields: jira.IssueFields{
				Summary: "Test issue",
				Status:  jira.Status{Name: "DONE"},
				Assignee: struct {
					AccountId   string `json:"accountId"`
					DisplayName string `json:"displayName"`
				}(struct {
					AccountId   string
					DisplayName string
				}{AccountId: "123", DisplayName: "Bob"}),
			},
		},
	}

	// when
	result := formatter.formatJiraIssues(issues)

	// then
	assertions.Len(result, 1)
	assertions.Equal("  TEST-123   Test issue     [DONE] - Bob", result[0])
}
