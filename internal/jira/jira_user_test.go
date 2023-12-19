package jira

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func Test_httpJiraApi_FindUsers(t *testing.T) {
	type args struct {
		project string
		query   string
	}
	tests := []struct {
		name    string
		args    args
		want    []User
		wantErr bool
	}{
		{"should find users without error",
			args{project: "FJIR"},
			[]User{
				{AccountId: "456", EmailAddress: "test@test.pl", DisplayName: "Mateusz Kulawik", Active: true, TimeZone: "Europe/Warsaw", Locale: "en_GB", AvatarUrls: nil},
				{AccountId: "123", EmailAddress: "", DisplayName: "mateusz.test", Active: true, TimeZone: "Europe/Warsaw", Locale: "en_US", AvatarUrls: nil},
			},
			false,
		},
		{"should find users with query without error",
			args{project: "FJIR", query: "test"},
			[]User{
				{AccountId: "456", EmailAddress: "test@test.pl", DisplayName: "Mateusz Kulawik", Active: true, TimeZone: "Europe/Warsaw", Locale: "en_GB", AvatarUrls: nil},
				{AccountId: "123", EmailAddress: "", DisplayName: "mateusz.test", Active: true, TimeZone: "Europe/Warsaw", Locale: "en_US", AvatarUrls: nil},
			},
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
        "accountId": "456",
        "accountType": "atlassian",
        "emailAddress": "test@test.pl",
        "displayName": "Mateusz Kulawik",
        "active": true,
        "timeZone": "Europe/Warsaw",
        "locale": "en_GB"
    },
    {
        "accountId": "123",
        "accountType": "atlassian",
        "emailAddress": "",
        "displayName": "mateusz.test",
        "active": true,
        "timeZone": "Europe/Warsaw",
        "locale": "en_US"
    }
]
`
				w.Write([]byte(body)) //nolint:errcheck
			})
			var got []User
			var err error
			if tt.args.query == "" {
				got, err = api.FindUsers(tt.args.project)
			} else {
				got, err = api.FindUsersWithQuery(tt.args.project, tt.args.query)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUsers() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpJiraApi_FindUsers_returnError(t *testing.T) {
	type args struct {
		project string
		query   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"should return error when search failed",
			args{project: "FJIR"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			api := NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(500)
			})

			// when
			_, err := api.FindUsersWithQuery(tt.args.project, tt.args.query)

			// then
			assert.Error(t, err)
		})
	}
}
