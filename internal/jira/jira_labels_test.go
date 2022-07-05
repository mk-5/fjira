package jira

import (
	"net/http"
	"reflect"
	"testing"
)

func Test_httpJiraApi_FindLabels(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{"should get labels without error",
			[]string{"Design", "TestLabel", "Windows"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				body := `
{
    "maxResults": 1000,
    "startAt": 0,
    "total": 3,
    "isLast": true,
    "values": [
        "Design",
        "TestLabel",
        "Windows"
    ]
}
`
				w.Write([]byte(body)) //nolint:errcheck
			})
			got, err := api.FindLabels()
			if (err != nil) != tt.wantErr {
				t.Errorf("FindLabels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindLabels() got = %v, want %v", got, tt.want)
			}
		})
	}
}
