package jira

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_httpApi_GetFilter(t *testing.T) {
	tests := []struct {
		name    string
		want    *Filter
		wantErr bool
	}{
		{"should get filter", &Filter{Name: "Filter for FJIR board"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				body := `
{
    "self": "https://test/rest/api/2/filter/10006",
    "id": "10006",
    "name": "Filter for FJIR board",
    "owner": {
        "self": "https://test",
        "accountId": "test",
        "avatarUrls": {
            "48x48": "https://test/48",
            "24x24": "https://test/24",
            "16x16": "https://test/16",
            "32x32": "https://test/32"
        },
        "displayName": "Test",
        "active": true
    },
    "jql": "project = FJIR ORDER BY Rank ASC",
    "viewUrl": "https://test/issues/?filter=10006",
    "searchUrl": "https://test/search?jql=project+%3D+FJIR+ORDER+BY+Rank+ASC",
    "favourite": false,
    "favouritedCount": 0,
    "sharePermissions": [
        {
            "id": 10007,
            "type": "project",
            "project": {
                "self": "https://test/rest/api/2/project/10006",
                "id": "10006",
                "key": "FJIR",
                "assigneeType": "PROJECT_LEAD",
                "name": "FJIRA",
                "roles": {},
                "avatarUrls": {
                    "48x48": "https://test",
                    "24x24": "https://test",
                    "16x16": "https://test",
                    "32x32": "https://test"
                },
                "projectTypeKey": "software",
                "simplified": true,
                "style": "next-gen",
                "properties": {},
                "entityId": "",
                "uuid": ""
            }
        }
    ],
    "editPermissions": [],
    "isWritable": true,
    "sharedUsers": {
        "size": 33,
        "items": [],
        "max-results": 1000,
        "start-index": 0,
        "end-index": 0
    },
    "subscriptions": {
        "size": 0,
        "items": [],
        "max-results": 1000,
        "start-index": 0,
        "end-index": 0
    }
}
`
				w.Write([]byte(body)) //nolint:errcheck
			})
			got, err := api.GetFilter("10006")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want.Name, got.Name, "GetFilter()")
			assert.Equal(t, "project = FJIR ORDER BY Rank ASC", got.JQL, "GetFilter()")
		})
	}
}

func Test_httpApi_GetMyFilters(t *testing.T) {
	tests := []struct {
		name    string
		want    *Filter
		wantErr bool
	}{
		{"should get my-filters", &Filter{Name: "Filter for FJIR board"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				body := `[
{
    "self": "https://test/rest/api/2/filter/10006",
    "id": "10006",
    "name": "Filter for FJIR board",
    "owner": {
        "self": "https://test",
        "accountId": "test",
        "avatarUrls": {
            "48x48": "https://test/48",
            "24x24": "https://test/24",
            "16x16": "https://test/16",
            "32x32": "https://test/32"
        },
        "displayName": "Test",
        "active": true
    },
    "jql": "project = FJIR ORDER BY Rank ASC",
    "viewUrl": "https://test/issues/?filter=10006",
    "searchUrl": "https://test/search?jql=project+%3D+FJIR+ORDER+BY+Rank+ASC",
    "favourite": false,
    "favouritedCount": 0,
    "sharePermissions": [
        {
            "id": 10007,
            "type": "project",
            "project": {
                "self": "https://test/rest/api/2/project/10006",
                "id": "10006",
                "key": "FJIR",
                "assigneeType": "PROJECT_LEAD",
                "name": "FJIRA",
                "roles": {},
                "avatarUrls": {
                    "48x48": "https://test",
                    "24x24": "https://test",
                    "16x16": "https://test",
                    "32x32": "https://test"
                },
                "projectTypeKey": "software",
                "simplified": true,
                "style": "next-gen",
                "properties": {},
                "entityId": "",
                "uuid": ""
            }
        }
    ],
    "editPermissions": [],
    "isWritable": true,
    "sharedUsers": {
        "size": 33,
        "items": [],
        "max-results": 1000,
        "start-index": 0,
        "end-index": 0
    },
    "subscriptions": {
        "size": 0,
        "items": [],
        "max-results": 1000,
        "start-index": 0,
        "end-index": 0
    }
}]
`
				w.Write([]byte(body)) //nolint:errcheck
			})
			got, err := api.GetMyFilters()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMyFilters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want.Name, got[0].Name, "GetMyFilters()")
			assert.Equal(t, "project = FJIR ORDER BY Rank ASC", got[0].JQL, "GetMyFilters()")
		})
	}
}
