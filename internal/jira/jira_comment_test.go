package jira

import (
	"net/http"
	"testing"
)

func Test_httpJiraApi_DoComment(t *testing.T) {
	type args struct {
		issueId     string
		commentBody string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"should do comment without error",
			args{commentBody: "Lorem ipsum", issueId: "ABC-123"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Write([]byte(``)) //nolint:errcheck
			})
			if err := api.DoComment(tt.args.issueId, tt.args.commentBody); (err != nil) != tt.wantErr {
				t.Errorf("DoComment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
