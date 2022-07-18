package fjira

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mk5/fjira/internal/app"
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

func Test_shouldInstallWithoutErrorWhenInstallEnvironmentProps(t *testing.T) {
	// given
	assert := assert2.New(t)
	os.Setenv(JiraTokenEnv, "test")
	os.Setenv(JiraUsernameEnv, "test")
	os.Setenv(JiraRestUrlEnv, "http://test.test")

	// when
	_, error := Install("abc")

	// then
	assert.NoError(error, "Should Install app without error when EnvironmentProps")
}

func Test_fjira_Close(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should close without error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			screen := tcell.NewSimulationScreen("utf-8")
			screen.Init() //nolint:errcheck
			defer screen.Fini()
			app.CreateNewAppWithScreen(screen)
			fjira := CreateNewFjira(&fjiraSettings{})

			// when
			fjira.Close()

			// then
			assert2.True(t, true) // no error during execution
		})
	}
}

func Test_readFromUserSettings(t *testing.T) {
	// TODO - not working on windows
	tempDir := t.TempDir()
	_ = os.Setenv("HOME", tempDir)
	_ = os.Mkdir(tempDir+"/.fjira", os.ModePerm) //nolint:errcheck
	defer func(name string) {
		_ = os.Remove(name)
	}(tempDir + "/.fjira")

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
			args{workspace: "xyz", storedWorkspaceFilename: "xyz.json", storedWorkspaceJson: "{\"jiraRestUrl\":\"https://test.atlassian.net\",\"jiraToken\":\"123\",\"jiraUsername\":\"test@test.pl\"}"},
			&fjiraSettings{JiraToken: "123", JiraUsername: "test@test.pl", JiraRestUrl: "https://test.atlassian.net"},
		},
		{"should read user settings from another workspace",
			args{workspace: "abc", storedWorkspaceFilename: "abc.json", storedWorkspaceJson: "{\"jiraRestUrl\":\"https://test\",\"jiraToken\":\"111\",\"jiraUsername\":\"test_user\"}"},
			&fjiraSettings{JiraToken: "111", JiraUsername: "test_user", JiraRestUrl: "https://test"},
		},
		{"should read user settings from default workspace",
			args{workspace: "", storedWorkspaceFilename: "default.json", storedWorkspaceJson: "{\"jiraRestUrl\":\"https://test\",\"jiraToken\":\"111\",\"jiraUsername\":\"test_user\"}"},
			&fjiraSettings{JiraToken: "111", JiraUsername: "test_user", JiraRestUrl: "https://test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, _ := os.Create(tempDir + "/.fjira/" + tt.args.storedWorkspaceFilename) //nolint:errcheck
			_, _ = file.WriteString(tt.args.storedWorkspaceJson)                         //nolint:errcheck
			_ = os.Symlink(tempDir+"/.fjira/_current.json", file.Name())                 //nolint:errcheck

			got, _ := readFromUserSettings(tt.args.workspace)
			assert2.Equalf(t, tt.want, got, "readFromUserSettings(%v)", tt.args.workspace)
		})
	}
}
