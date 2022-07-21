package jira

import (
	assert2 "github.com/stretchr/testify/assert"
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
    "token": "",
    "suggestions": [
        {
            "label": "Design",
            "html": "<b></b>Design"
        },
        {
            "label": "TestLabel",
            "html": "<b></b>TestLabel"
        },
        {
            "label": "Windows",
            "html": "<b></b>Windows"
        }
    ]
}
`
				_, _ = w.Write([]byte(body))
			})
			got, err := api.FindLabels(&JiraIssue{}, "")
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

func Test_httpJiraApi_AddLabel(t *testing.T) {
	tests := []struct {
		name     string
		issueKey string
		label    string
		wantErr  bool
	}{
		{"should add label without error",
			"PROJ-123", "Test", false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				body := ""
				w.Write([]byte(body)) //nolint:errcheck
			})
			err := api.AddLabel(tt.issueKey, tt.label)
			assert2.Nil(t, err)
		})
	}
}
