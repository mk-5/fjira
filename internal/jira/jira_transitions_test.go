package jira

import (
	"net/http"
	"reflect"
	"testing"
)

func Test_httpJiraApi_DoTransition(t *testing.T) {
	type args struct {
		issueId    string
		transition *IssueTransition
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"should do transition without error",
			args{transition: &IssueTransition{
				Id:   "test",
				Name: "test",
				To: struct {
					StatusUrl string `json:"self"`
					StatusId  string `json:"id"`
					Name      string `json:"name"`
				}(struct {
					StatusUrl string
					StatusId  string
					Name      string
				}{"test2", "test2", "test2"}),
			}, issueId: "ABC-123"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Write([]byte(``)) //nolint:errcheck
			})
			if err := api.DoTransition(tt.args.issueId, tt.args.transition); (err != nil) != tt.wantErr {
				t.Errorf("DoTransition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_httpJiraApi_FindTransitions(t *testing.T) {
	type args struct {
		issueId string
	}
	tests := []struct {
		name    string
		args    args
		want    []IssueTransition
		wantErr bool
	}{
		{"should find transitions without error",
			args{issueId: "ABC-123"},
			[]IssueTransition{{Id: "11", Name: "To Do"}, {Id: "21", Name: "In Progress"}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				body := `
{
    "transitions": [
        {
            "id": "11",
            "name": "To Do"
        },
        {
            "id": "21",
            "name": "In Progress"
        }
    ]
}
`
				w.Write([]byte(body)) //nolint:errcheck
			})
			got, err := api.FindTransitions(tt.args.issueId)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindTransitions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindTransitions() got = %v, want %v", got, tt.want)
			}
		})
	}
}
