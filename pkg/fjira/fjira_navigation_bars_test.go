package fjira

import (
	"github.com/mk5/fjira/internal/jira"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateNavigationBars(t *testing.T) {
	tests := []struct {
		name     string
		supplier func() interface{}
	}{
		{"should create comment bar item", func() interface{} {
			return CreateCommentBarItem()
		}},
		{"should create project bottom bar", func() interface{} {
			return CreateNewEmptyProjectBottomBar()
		}},
		{"should create issue bottom bar", func() interface{} {
			return CreateNewIssueBottomBar(&jira.JiraIssue{})
		}},
		{"should create issue top bar", func() interface{} {
			return CreateNewIssueTopBar(&jira.JiraIssue{})
		}},
		{"should create search issues bottom bar", func() interface{} {
			return CreateNewSearchIssuesBottomBar(&jira.JiraProject{})
		}},
		{"should create search issues top bar", func() interface{} {
			return CreateNewSearchIssuesTopBar()
		}},
		{"should create assignee change bar item", func() interface{} {
			return NewAssigneeChangeBarItem()
		}},
		{"should create by-assignee change bar item", func() interface{} {
			return NewByAssigneeBarItem()
		}},
		{"should create by-status change bar item", func() interface{} {
			return NewByStatusBarItem()
		}},
		{"should create cancel bar item", func() interface{} {
			return NewCancelBarItem()
		}},
		{"should create new assignee bar item", func() interface{} {
			return NewNewAssigneeBarItem(&jira.JiraUser{})
		}},
		{"should create new status bar item", func() interface{} {
			return NewNewStatusBarItem("test")
		}},
		{"should create new save bar item", func() interface{} {
			return NewSaveBarItem()
		}},
		{"should create new status change bar item", func() interface{} {
			return NewStatusChangeBarItem()
		}},
		{"should create new YES change bar item", func() interface{} {
			return NewYesBarItem()
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNilf(t, tt.supplier(), "CreateCommentBarItem()")
		})
	}
}
