package jira

import (
	"net/http"
	"reflect"
	"testing"
)

func Test_httpJiraApi_FindProjects(t *testing.T) {
	tests := []struct {
		name    string
		want    []Project
		wantErr bool
	}{
		{"should find projects without error",
			[]Project{{"1", "FJIRA", "FJIR"}, {"2", "General", "GEN"}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				body := `
{
  "values": [
    {
        "expand": "description,lead,issueTypes,url,projectKeys,permissions,insight",
        "id": "1",
        "key": "FJIR",
        "name": "FJIRA",
        "projectTypeKey": "software",
        "simplified": true,
        "style": "next-gen",
        "isPrivate": false,
        "properties": {},
        "entityId": "250cd492-c831-44d9-ae5c-17bd93922fa6",
        "uuid": "250cd492-c831-44d9-ae5c-17bd93922fa6"
    },
    {
        "expand": "description,lead,issueTypes,url,projectKeys,permissions,insight",
        "id": "2",
        "key": "GEN",
        "name": "General",
        "projectTypeKey": "software",
        "simplified": false,
        "style": "classic",
        "isPrivate": false,
        "properties": {}
    }
]
}
`
				w.Write([]byte(body)) //nolint:errcheck
			})
			got, err := api.FindProjects()
			if (err != nil) != tt.wantErr {
				t.Errorf("FindProjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindProjects() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpJiraApi_FindProject(t *testing.T) {
	tests := []struct {
		name    string
		want    *Project
		wantErr bool
	}{
		{"should find project without error",
			&Project{"1", "FJIRA", "FJIR"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				body := `
    {
        "expand": "description,lead,issueTypes,url,projectKeys,permissions,insight",
        "id": "1",
        "key": "FJIR",
        "name": "FJIRA",
        "projectTypeKey": "software",
        "simplified": true,
        "style": "next-gen",
        "isPrivate": false,
        "properties": {},
        "entityId": "250cd492-c831-44d9-ae5c-17bd93922fa6",
        "uuid": "250cd492-c831-44d9-ae5c-17bd93922fa6"
    }
`
				w.Write([]byte(body)) //nolint:errcheck
			})
			got, err := api.FindProject("FJIR")
			if (err != nil) != tt.wantErr {
				t.Errorf("FindProjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindProjects() got = %v, want %v", got, tt.want)
			}
		})
	}
}
