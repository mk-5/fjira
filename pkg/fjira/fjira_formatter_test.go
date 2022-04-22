package fjira

import (
	"github.com/mk5/fjira/internal/jira"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_shouldFormatProjects(t *testing.T) {
	assert := assert.New(t)

	// given
	formatter := defaultFormatter{}
	issues := []jira.JiraIssue{
		jira.JiraIssue{
			Key: "TEST-123",
			Fields: jira.JiraIssueFields{
				Summary: "Test issue",
				Status: struct {
					Name string `json:"name"`
				}(struct{ Name string }{Name: "DONE"}),
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
	assert.Len(result, 1)
	assert.Equal("  TEST-123   Test issue     [DONE] - Bob", result[0])
}
