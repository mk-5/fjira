package jira

import (
	"net/http"
	"testing"
)

func Test_httpJiraApi_DoAssignee(t *testing.T) {
	type args struct {
		issueId   string
		accountId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"should do assignment without error",
			args{accountId: "123456", issueId: "ABC-123"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Write([]byte(``)) //nolint:errcheck
			})
			if err := api.DoAssignee(tt.args.issueId, tt.args.accountId); (err != nil) != tt.wantErr {
				t.Errorf("DoAssignee() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
