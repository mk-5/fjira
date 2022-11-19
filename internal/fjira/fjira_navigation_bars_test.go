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
			return CreateProjectBottomBar()
		}},
		{"should create issue bottom bar", func() interface{} {
			return CreateBottomLeftBar()
		}},
		{"should create issue top bar", func() interface{} {
			return CreateIssueTopBar(&jira.Issue{})
		}},
		{"should create search issues bottom bar", func() interface{} {
			return CreateSearchIssuesBottomBar()
		}},
		{"should create search issues top bar", func() interface{} {
			return CreateSearchIssuesTopBar(&jira.Project{})
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
		{"should create new save bar item", func() interface{} {
			return NewSaveBarItem()
		}},
		{"should create new status change bar item", func() interface{} {
			return NewStatusChangeBarItem()
		}},
		{"should create new YES change bar item", func() interface{} {
			return NewYesBarItem()
		}},
		{"should create new OPEN bar item", func() interface{} {
			return NewOpenBarItem()
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNilf(t, tt.supplier(), "CreateCommentBarItem()")
		})
	}
}
