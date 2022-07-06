package fjira

import (
	"github.com/mk5/fjira/internal/jira"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_buildSearchIssuesJql(t *testing.T) {
	type args struct {
		project *jira.JiraProject
		query   string
		status  *jira.JiraIssueStatus
		user    *jira.JiraUser
		label   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"should create valid jql", args{project: &jira.JiraProject{Id: "123"}}, "project=123 ORDER BY status"},
		{"should create valid jql", args{project: &jira.JiraProject{Id: "123"}, query: "abc"}, "project=123 AND summary~\"abc*\" ORDER BY status"},
		{"should create valid jql", args{project: &jira.JiraProject{Id: MessageAll, Key: MessageAll}, query: "abc"}, "summary~\"abc*\" ORDER BY status"},
		{"should create valid jql", args{
			project: &jira.JiraProject{Id: "123"}, query: "abc", status: &jira.JiraIssueStatus{Id: "st1"}},
			"project=123 AND summary~\"abc*\" AND status=st1 ORDER BY status",
		},
		{"should create valid jql", args{
			project: &jira.JiraProject{Id: "123"}, query: "abc", status: &jira.JiraIssueStatus{Id: "st1"}, user: &jira.JiraUser{AccountId: "us1"}},
			"project=123 AND summary~\"abc*\" AND status=st1 AND assignee=us1 ORDER BY status",
		},
		{"should create valid jql", args{project: &jira.JiraProject{Id: "123"}, label: "test"}, "project=123 AND labels=test ORDER BY status"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, buildSearchIssuesJql(tt.args.project, tt.args.query, tt.args.status, tt.args.user, tt.args.label), "buildSearchIssuesJql(%v, %v, %v, %v)", tt.args.project, tt.args.query, tt.args.status, tt.args.user)
		})
	}
}
