package jira

import (
	"net/http"
	"testing"
)

func Test_httpJiraApi_Search(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		args    args
		want    []JiraIssue
		want1   int32
		wantErr bool
	}{
		{"should do search without error",
			args{query: "test"},
			[]JiraIssue{
				{"ISSUE-1", "", "", JiraIssueFields{Description: "Desc1", Status: struct {
					Name string `json:"name"`
				}(struct{ Name string }{"Status1"})}, ""},
				{"ISSUE-2", "", "", JiraIssueFields{Description: "Desc2", Status: struct {
					Name string `json:"name"`
				}(struct{ Name string }{"Status2"})}, ""},
				{"ISSUE-3", "", "", JiraIssueFields{Description: "Desc3", Status: struct {
					Name string `json:"name"`
				}(struct{ Name string }{"Status3"})}, ""},
			},
			3,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				body := `
{
    "expand": "schema,names",
    "startAt": 0,
    "maxResults": 100,
    "total": 3,
    "issues": [
        {
            "key": "ISSUE-1",
            "fields": {
                "summary": "Issue summary 1",
				"description": "Desc1",
                "status": {
                    "name": "Status1"
                }
            }
        },
        {
            "key": "ISSUE-2",
            "fields": {
                "summary": "Issue summary 2",
				"description": "Desc2",
                "status": {
                    "name": "Status2"
                }
            }
        },
        {
            "key": "ISSUE-3",
            "fields": {
                "summary": "Issue summary 3",
				"description": "Desc3",
                "status": {
                    "name": "Status3"
                }
            }
        }
    ]
}
`
				w.Write([]byte(body)) //nolint:errcheck
			})
			got, got1, err := api.Search(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := range got {
				if got[i].Key != tt.want[i].Key {
					t.Errorf("Search() got = %v, want %v", got[i].Key, tt.want[i].Key)
				}
			}
			if got1 != tt.want1 {
				t.Errorf("Search() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_httpJiraApi_SearchJql(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		args    args
		want    []JiraIssue
		wantErr bool
	}{
		{"should do search-jql without error",
			args{query: "test"},
			[]JiraIssue{
				{"ISSUE-1", "", "", JiraIssueFields{Description: "Desc1", Status: struct {
					Name string `json:"name"`
				}(struct{ Name string }{"Status1"})}, ""},
				{"ISSUE-2", "", "", JiraIssueFields{Description: "Desc2", Status: struct {
					Name string `json:"name"`
				}(struct{ Name string }{"Status2"})}, ""},
				{"ISSUE-3", "", "", JiraIssueFields{Description: "Desc3", Status: struct {
					Name string `json:"name"`
				}(struct{ Name string }{"Status3"})}, ""},
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
    "expand": "schema,names",
    "startAt": 0,
    "maxResults": 100,
    "total": 3,
    "issues": [
        {
            "key": "ISSUE-1",
            "fields": {
                "summary": "Issue summary 1",
				"description": "Desc1",
                "status": {
                    "name": "Status1"
                }
            }
        },
        {
            "key": "ISSUE-2",
            "fields": {
                "summary": "Issue summary 2",
				"description": "Desc2",
                "status": {
                    "name": "Status2"
                }
            }
        },
        {
            "key": "ISSUE-3",
            "fields": {
                "summary": "Issue summary 3",
				"description": "Desc3",
                "status": {
                    "name": "Status3"
                }
            }
        }
    ]
}
`
				w.Write([]byte(body)) //nolint:errcheck
			})
			got, err := api.SearchJql(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := range got {
				if got[i].Key != tt.want[i].Key {
					t.Errorf("Search() got = %v, want %v", got[i].Key, tt.want[i].Key)
				}
			}
		})
	}
}
