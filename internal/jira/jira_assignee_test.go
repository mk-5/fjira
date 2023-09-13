package jira

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_httpJiraApi_DoAssigneeCloud(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should do CLOUD assignment without error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			var url string
			api := NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Write([]byte(``)) //nolint:errcheck
				url = r.URL.Path
			})

			// when
			err := api.DoAssignee("ISS123", &User{AccountId: "acc123"})

			// then
			assert.Nil(t, err)
			assert.Equal(t, "/rest/api/2/issue/ISS123/assignee", url)
		})
	}
}

func Test_httpJiraApi_DoAssigneeOnPremise(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should do SERVER assignment without error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			var url string
			api := NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Write([]byte(``)) //nolint:errcheck
				url = r.URL.Path
			})

			// when
			err := api.DoAssignee("ISS123", &User{Name: "username"})

			// then
			assert.Nil(t, err)
			assert.Equal(t, "/rest/api/2/issue/ISS123", url)
		})
	}
}
