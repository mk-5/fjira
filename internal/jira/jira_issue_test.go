package jira

import (
	"io"
	"testing"
)

type getDetailedIssueMock struct{}

func (api *getDetailedIssueMock) jiraRequest(method string, restPath string, queryParams interface{}, reqBody io.Reader) ([]byte, error) {
	return []byte(`
{
    "expand": "renderedFields,names,schema,operations,editmeta,changelog,versionedRepresentations,customfield_10010.requestTypePractice",
    "id": "10055",
    "key": "HD-17",
    "fields": {
        "issuetype": {
            "id": "10019",
            "description": "A small, distinct piece of work.",
            "name": "Task",
        },
        "timespent": null,
        "project": {
            "id": "10005",
            "key": "HD",
            "name": "Hacker Dungeon",
        },
        "created": "2022-03-19T11:11:56.993+0100",
        "assignee": {
            "accountId": "123",
            "emailAddress": "test@test.dev",
            "displayName": "Test",
            "active": true,
        },
        "status": {
            "name": "In Progress",
            "id": "10018",
        },
        "description": "Lorem ipsum",
        "summary": "Test",
        "subtasks": []
    }
}
`), nil
}

func TestHttpJiraApi_GetIssueDetailed(t *testing.T) {
	// given
	//client := http.Client{}
	//api := &httpJiraApi{client: http.Client{}, restUrl: nil}
}
