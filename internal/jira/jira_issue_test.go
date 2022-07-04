package jira

import (
	"net/http"
	"reflect"
	"testing"
)

func Test_httpJiraApi_GetIssueDetailed(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    *JiraIssue
		wantErr bool
	}{
		{"should get detailed jira issue",
			args{id: "10011"},
			&JiraIssue{
				Key: "JWC-3", Id: "10011",
				Fields: JiraIssueFields{
					Summary:     "Tutorial - create tutorial",
					Description: "Lorem ipsum",
					Project:     JiraProject{Id: "10003", Name: "JIRA WORK CHART", Key: "JWC"},
					Reporter: struct {
						AccountId   string `json:"accountId"`
						DisplayName string `json:"displayName"`
					}(struct {
						AccountId   string
						DisplayName string
					}{"607f55ba074a0b006a6cb482", "Mateusz Kulawik"}),
					Assignee: struct {
						AccountId   string `json:"accountId"`
						DisplayName string `json:"displayName"`
					}(struct {
						AccountId   string
						DisplayName string
					}{"", ""}),
					Type:   JiraIssueType{Name: "Task"},
					Labels: []string{"TestLabel"},
					Status: struct {
						Name string `json:"name"`
					}(struct{ Name string }{"Done"}),
					Comment: struct {
						Comments   []JiraComment `json:"comments"`
						MaxResults int32         `json:"maxResults"`
						Total      int32         `json:"total"`
						StartAt    int32         `json:"startAt"`
					}(struct {
						Comments   []JiraComment
						MaxResults int32
						Total      int32
						StartAt    int32
					}{
						Comments: []JiraComment{
							{Body: "Comment 123-ABC", Created: "2022-06-09T22:53:42.057+0200", Author: JiraUser{DisplayName: "Mateusz Kulawik"}},
						},
						MaxResults: 1, Total: 1, StartAt: 0},
					),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				body := `
{
    "expand": "renderedFields,names,schema,operations,editmeta,changelog,versionedRepresentations,customfield_10010.requestTypePractice",
    "id": "10011",
    "self": "https://test/rest/api/2/issue/10011",
    "key": "JWC-3",
    "fields": {
        "issuetype": {
            "id": "10013",
            "description": "A small, distinct piece of work.",
            "name": "Task",
            "subtask": false
        },
        "timespent": 14400,
        "project": {
            "id": "10003",
            "key": "JWC",
            "name": "JIRA WORK CHART",
            "projectTypeKey": "software"
        },
        "fixVersions": [],
        "aggregatetimespent": 14400,
        "resolutiondate": "2022-02-22T00:27:11.861+0100",
        "workratio": -1,
        "issuerestriction": {
            "issuerestrictions": {},
            "shouldDisplay": true
        },
        "lastViewed": "2022-02-22T00:27:17.356+0100",
        "created": "2021-10-02T22:34:22.521+0200",
        "aggregatetimeoriginalestimate": null,
        "timeestimate": 0,
        "versions": [],
        "issuelinks": [],
        "assignee": null,
        "updated": "2022-02-22T00:27:19.792+0100",
        "status": {
            "description": "",
            "name": "Done",
            "id": "10013"
        },
		"labels": ["TestLabel"],
        "description": "Lorem ipsum",
        "summary": "Tutorial - create tutorial",
        "creator": {
            "accountId": "607f55ba074a0b006a6cb482",
            "emailAddress": "test@test.dev",
            "displayName": "Mateusz Kulawik",
            "active": true,
            "timeZone": "Europe/Warsaw",
            "accountType": "atlassian"
        },
        "subtasks": [],
        "reporter": {
			"accountId": "607f55ba074a0b006a6cb482",
            "emailAddress": "test@test.dev",
            "displayName": "Mateusz Kulawik",
            "active": true,
            "timeZone": "Europe/Warsaw",
            "accountType": "atlassian"
        },
 		"comment": {
            "comments": [
                {
                    "author": {
                        "displayName": "Mateusz Kulawik"
                    },
                    "body": "Comment 123-ABC",
                    "created": "2022-06-09T22:53:42.057+0200",
                    "updated": "2022-06-09T22:53:42.057+0200"
                }
            ],
            "maxResults": 1,
            "total": 1,
            "startAt": 0
        }
    }
}
`
				w.Write([]byte(body)) //nolint:errcheck
			})
			got, err := api.GetIssueDetailed(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIssueDetailed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIssueDetailed() got = %v, want %v", got, tt.want)
			}
		})
	}
}
