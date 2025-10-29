package jira

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_httpApi_FindBoards(t *testing.T) {
	tests := []struct {
		name    string
		want    []BoardItem
		wantErr bool
	}{
		{"should find boards without error", []BoardItem{BoardItem{Id: 1, Name: "GEN board", Self: "https://test.net/rest/agile/1.0/board/1", Type: "kanban"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				body := `
{
    "maxResults": 50,
    "startAt": 0,
    "total": 1,
    "isLast": true,
    "values": [
        {
            "id": 1,
            "self": "https://test.net/rest/agile/1.0/board/1",
            "name": "GEN board",
            "type": "kanban",
            "location": {
                "projectId": 10000,
                "displayName": "General (GEN)",
                "projectName": "General",
                "projectKey": "GEN",
                "projectTypeKey": "software",
                "avatarURI": "https://test.net/rest/api/2/universal_avatar/view/type/project/avatar/10416?size=small",
                "name": "General (GEN)"
            }
        }
    ]
}
`
				w.Write([]byte(body)) //nolint:errcheck
			})
			got, err := api.FindBoards("1")
			if (err != nil) != tt.wantErr {
				t.Errorf("FindBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.EqualValues(t, tt.want[0], got[0], "FindBoards()")
		})
	}
}

func Test_httpApi_GetBoardConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		want    *BoardConfiguration
		wantErr bool
	}{
		{"should get board configuration", &BoardConfiguration{Name: "GEN board"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				body := `
{
    "id": 1,
    "name": "GEN board",
    "type": "kanban",
    "self": "https://test.net/rest/agile/1.0/board/1/configuration",
    "location": {
        "type": "project",
        "key": "GEN",
        "id": "10000",
        "self": "https://test.net/rest/api/2/project/10000",
        "name": "General"
    },
    "filter": {
        "id": "10000",
        "self": "https://test.net/rest/api/2/filter/10000"
    },
    "subQuery": {
        "query": "fixVersion in unreleasedVersions() OR fixVersion is EMPTY"
    },
    "columnConfig": {
        "columns": [
            {
                "name": "Backlog",
                "statuses": []
            },
            {
                "name": "Backlog",
                "statuses": [
                    {
                        "id": "10000",
                        "self": "https://test.net/rest/api/2/status/10000"
                    }
                ]
            },
            {
                "name": "Selected for Development",
                "statuses": [
                    {
                        "id": "10001",
                        "self": "https://test.net/rest/api/2/status/10001"
                    }
                ]
            },
            {
                "name": "In Progress",
                "statuses": [
                    {
                        "id": "3",
                        "self": "https://test.net/rest/api/2/status/3"
                    }
                ]
            },
            {
                "name": "Done",
                "statuses": [
                    {
                        "id": "10002",
                        "self": "https://test.net/rest/api/2/status/10002"
                    }
                ]
            }
        ],
        "constraintType": "issueCount"
    },
    "ranking": {
        "rankCustomFieldId": 10019
    }
}
`
				w.Write([]byte(body)) //nolint:errcheck
			})
			got, err := api.GetBoardConfiguration(1)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBoardConfiguration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want.Name, got.Name, "GetBoardConfiguration()")
			assert.Equal(t, 5, len(got.ColumnConfig.Columns), "GetBoardConfiguration()")
		})
	}
}

func Test_httpApi_GetBoardSprints(t *testing.T) {
	api := NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
		// validate query parameters
		assert.Contains(t, r.URL.RawQuery, "state=active%2Cfuture")

		w.WriteHeader(200)
		body := `
{
  "maxResults": 50,
  "startAt": 0,
  "total": 2,
  "isLast": true,
  "values": [
    {
      "id": 10,
      "self": "https://test.net/rest/agile/1.0/sprint/10",
      "state": "active",
      "name": "Sprint 10",
      "startDate": null,
      "endDate": null,
      "completeDate": null,
      "originBoardId": 1,
      "goal": ""
    },
    {
      "id": 11,
      "self": "https://test.net/rest/agile/1.0/sprint/11",
      "state": "future",
      "name": "Sprint 11",
      "startDate": null,
      "endDate": null,
      "completeDate": null,
      "originBoardId": 1,
      "goal": ""
    }
  ]
}
`
		w.Write([]byte(body)) //nolint:errcheck
	})

	got, err := api.GetBoardSprints(1)
	assert.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, 10, got[0].Id)
	assert.Equal(t, "future", got[1].State)
}

func Test_httpApi_GetBoardSprintIssues(t *testing.T) {
	page := int32(2)
	pageSize := int32(25)
	api := NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
		// validate pagination and fields params
		q := r.URL.Query()
		assert.Equal(t, "25", q.Get("maxResults"))
		assert.Equal(t, "50", q.Get("startAt"))
		assert.Equal(t, "id,key,summary,issuetype,project,reporter,status,assignee", q.Get("fields"))

		w.WriteHeader(200)
		body := `
{
  "total": 123,
  "maxResults": 25,
  "isLast": false,
  "issues": [
    { "id": "10001", "key": "GEN-1" },
    { "id": "10002", "key": "GEN-2" }
  ]
}
`
		w.Write([]byte(body)) //nolint:errcheck
	})

	issues, total, max, err := api.GetBoardSprintIssues(1, 10, page, pageSize)
	assert.NoError(t, err)
	assert.Equal(t, int32(123), total)
	assert.Equal(t, int32(25), max)
	assert.Len(t, issues, 2)
	assert.Equal(t, "GEN-1", issues[0].Key)
	assert.Equal(t, "10002", issues[1].Id)
}
