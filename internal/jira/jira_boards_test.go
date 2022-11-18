package jira

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_httpApi_FindBoards(t *testing.T) {
	tests := []struct {
		name    string
		want    []*BoardItem
		wantErr bool
	}{
		{"should find boards without error", []*BoardItem{&BoardItem{Id: 1, Name: "GEN board", Self: "https://test.net/rest/agile/1.0/board/1", Type: "kanban"}}, false},
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
