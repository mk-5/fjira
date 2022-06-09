package fjira

import (
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_shouldReturnErrorWhenNoEnvironments(t *testing.T) {
	// given
	assert := assert2.New(t)
	os.Setenv(JiraTokenEnv, "")
	os.Setenv(JiraUsernameEnv, "")
	os.Setenv(JiraRestUrlEnv, "")

	// when
	_, error := readFromEnvironments()

	// then
	assert.Error(error, "Should return error when no fjira environments")
}

func Test_shouldReturnNoErrorWhenEnvironments(t *testing.T) {
	// given
	assert := assert2.New(t)
	os.Setenv(JiraTokenEnv, "test")
	os.Setenv(JiraUsernameEnv, "test")
	os.Setenv(JiraRestUrlEnv, "http://test.test")

	// when
	_, error := readFromEnvironments()

	// then
	assert.NoError(error, "Should return no error when fjira environments")
}

func Test_readFromUserSettings(t *testing.T) {
	// TODO - it's not multi-platform
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	os.Mkdir(tempDir+"/.fjira", os.ModePerm) //nolint:errcheck
	defer os.Remove(tempDir + "/.fjira")

	type args struct {
		workspace               string
		storedWorkspaceFilename string
		storedWorkspaceJson     string
	}
	tests := []struct {
		name string
		args args
		want *fjiraSettings
	}{
		{"should read user settings from current workspace",
			args{workspace: "xyz", storedWorkspaceFilename: "_xyz.json", storedWorkspaceJson: "{\"jiraRestUrl\":\"https://test.atlassian.net\",\"jiraToken\":\"123\",\"jiraUsername\":\"test@test.pl\"}"},
			&fjiraSettings{JiraToken: "123", JiraUsername: "test@test.pl", JiraRestUrl: "https://test.atlassian.net"},
		},
		{"should read user settings from another workspace",
			args{workspace: "abc", storedWorkspaceFilename: "_abc.json", storedWorkspaceJson: "{\"jiraRestUrl\":\"https://test\",\"jiraToken\":\"111\",\"jiraUsername\":\"test_user\"}"},
			&fjiraSettings{JiraToken: "111", JiraUsername: "test_user", JiraRestUrl: "https://test"},
		},
		{"should read user settings from default workspace",
			args{workspace: "", storedWorkspaceFilename: "_default.json", storedWorkspaceJson: "{\"jiraRestUrl\":\"https://test\",\"jiraToken\":\"111\",\"jiraUsername\":\"test_user\"}"},
			&fjiraSettings{JiraToken: "111", JiraUsername: "test_user", JiraRestUrl: "https://test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, _ := os.Create(tempDir + "/.fjira/" + tt.args.storedWorkspaceFilename) //nolint:errcheck
			file.WriteString(tt.args.storedWorkspaceJson)                                //nolint:errcheck

			got, _ := readFromUserSettings(tt.args.workspace)
			assert2.Equalf(t, tt.want, got, "readFromUserSettings(%v)", tt.args.workspace)
		})
	}
}
