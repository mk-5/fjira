package jira

import (
	"net/http"
	"reflect"
	"testing"
)

func Test_httpJiraApi_FindProjectStatuses(t *testing.T) {
	type args struct {
		projectId string
	}
	tests := []struct {
		name    string
		args    args
		want    []JiraIssueStatus
		wantErr bool
	}{
		{"should get project statuses",
			args{projectId: "123"},
			[]JiraIssueStatus{{"10011", "To Do", ""}, {"10012", "In Progress", "ABC"}, {"10013", "Done", ""}, {"10016", "Verification", "XXX"}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				body := `
[
    {
        "id": "10013",
        "name": "Task",
        "subtask": false,
        "statuses": [
            {
                "description": "",
                "name": "To Do",
                "untranslatedName": "To Do",
                "id": "10011",
                "statusCategory": {
                    "id": 2,
                    "key": "new",
                    "colorName": "blue-gray",
                    "name": "To Do"
                },
                "scope": {
                    "type": "PROJECT",
                    "project": {
                        "id": "10003"
                    }
                }
            },
            {
                "description": "ABC",
                "name": "In Progress",
                "untranslatedName": "In Progress",
                "id": "10012",
                "statusCategory": {
                    "id": 4,
                    "key": "indeterminate",
                    "colorName": "yellow",
                    "name": "In Progress"
                },
                "scope": {
                    "type": "PROJECT",
                    "project": {
                        "id": "10003"
                    }
                }
            },
            {
                "description": "",
                "name": "Done",
                "untranslatedName": "Done",
                "id": "10013",
                "statusCategory": {
                    "id": 3,
                    "key": "done",
                    "colorName": "green",
                    "name": "Done"
                },
                "scope": {
                    "type": "PROJECT",
                    "project": {
                        "id": "10003"
                    }
                }
            }
        ]
    },
    {
        "id": "10014",
        "name": "Epic",
        "subtask": false,
        "statuses": [
            {
                "description": "",
                "name": "To Do",
                "untranslatedName": "To Do",
                "id": "10011",
                "statusCategory": {
                    "id": 2,
                    "key": "new",
                    "colorName": "blue-gray",
                    "name": "To Do"
                },
                "scope": {
                    "type": "PROJECT",
                    "project": {
                        "id": "10003"
                    }
                }
            },
            {
                "description": "ABC",
                "name": "In Progress",
                "untranslatedName": "In Progress",
                "id": "10012",
                "statusCategory": {
                    "id": 4,
                    "key": "indeterminate",
                    "colorName": "yellow",
                    "name": "In Progress"
                },
                "scope": {
                    "type": "PROJECT",
                    "project": {
                        "id": "10003"
                    }
                }
            },
            {
                "description": "",
                "name": "Done",
                "untranslatedName": "Done",
                "id": "10013",
                "statusCategory": {

                    "id": 3,
                    "key": "done",
                    "colorName": "green",
                    "name": "Done"
                },
                "scope": {
                    "type": "PROJECT",
                    "project": {
                        "id": "10003"
                    }
                }
            }
        ]
    },
    {
        "id": "10015",
        "name": "Subtask",
        "subtask": true,
        "statuses": [
            {
                "description": "",
                "name": "To Do",
                "untranslatedName": "To Do",
                "id": "10011",
                "statusCategory": {
                    "id": 2,
                    "key": "new",
                    "colorName": "blue-gray",
                    "name": "To Do"
                },
                "scope": {
                    "type": "PROJECT",
                    "project": {
                        "id": "10003"
                    }
                }
            },
            {
                "description": "XXX",
                "name": "Verification",
                "untranslatedName": "Verification",
                "id": "10016",
                "statusCategory": {
                    "id": 4,
                    "key": "indeterminate",
                    "colorName": "yellow",
                    "name": "In Progress"
                },
                "scope": {
                    "type": "PROJECT",
                    "project": {
                        "id": "10003"
                    }
                }
            }
        ]
    }
]
`
				w.Write([]byte(body)) //nolint:errcheck
			})
			got, err := api.FindProjectStatuses(tt.args.projectId)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindProjectStatuses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindProjectStatuses() got = %v, want %v", got, tt.want)
			}
		})
	}
}
