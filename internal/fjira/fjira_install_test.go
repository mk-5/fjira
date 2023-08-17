package fjira

import (
	"bytes"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	os2 "github.com/mk-5/fjira/internal/os"
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func Test_shouldReturnErrorWhenNoEnvironments(t *testing.T) {
	// given
	assert := assert2.New(t)
	_ = os.Setenv(JiraTokenEnv, "")
	_ = os.Setenv(JiraUsernameEnv, "")
	_ = os.Setenv(JiraRestUrlEnv, "")

	// when
	_, err := readFromEnvironments()

	// then
	assert.Error(err, "Should return error when no fjira environments")
}

func Test_shouldReturnNoErrorWhenEnvironments(t *testing.T) {
	// given
	assert := assert2.New(t)
	_ = os.Setenv(JiraTokenEnv, "test")
	_ = os.Setenv(JiraUsernameEnv, "test")
	_ = os.Setenv(JiraRestUrlEnv, "http://test.test")

	// when
	_, err := readFromEnvironments()

	// then
	assert.NoError(err, "Should return no error when fjira environments")
}

func Test_shouldInstallWithoutErrorWhenInstallEnvironmentProps(t *testing.T) {
	// given
	assert := assert2.New(t)
	_ = os.Setenv(JiraTokenEnv, "test")
	_ = os.Setenv(JiraUsernameEnv, "test")
	_ = os.Setenv(JiraRestUrlEnv, "http://test.test")

	// when
	_, err := Install(CliArgs{Workspace: "abc"})

	// then
	assert.NoError(err, "Should Install app without error when EnvironmentProps")
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
			_ = screen.Init() //nolint:errcheck
			defer screen.Fini()
			app.CreateNewAppWithScreen(screen)
			fjira := CreateNewFjira(&fjiraWorkspaceSettings{})

			// when
			fjira.Close()

			// then
			assert2.True(t, true) // no error during execution
		})
	}
}

func Test_readFromUserSettings(t *testing.T) {
	tempDir := t.TempDir()
	_ = os2.SetUserHomeDir(tempDir)
	_ = os.Mkdir(tempDir+"/.fjira", os.ModePerm) //nolint:errcheck
	t.Cleanup(func() {
		_ = os.Remove(tempDir)
	})

	type args struct {
		workspace           string
		storedWorkspaceYaml string
	}
	tests := []struct {
		name string
		args args
		want *fjiraWorkspaceSettings
	}{
		{"should read user settings from current workspace",
			args{workspace: "xyz", storedWorkspaceYaml: "\ncurrent: xyz\nworkspaces:\n    xyz:\n        jiraRestUrl: https://test.atlassian.net\n        jiraToken: 123\n        jiraUsername: test@test.pl"},
			&fjiraWorkspaceSettings{JiraToken: "123", JiraUsername: "test@test.pl", JiraRestUrl: "https://test.atlassian.net", Workspace: "xyz"},
		},
		{"should read user settings from another workspace",
			args{workspace: "abc", storedWorkspaceYaml: "\ncurrent: default\nworkspaces:\n    abc:\n        jiraRestUrl: https://test\n        jiraToken: 111\n        jiraUsername: test_user"},
			&fjiraWorkspaceSettings{JiraToken: "111", JiraUsername: "test_user", JiraRestUrl: "https://test", Workspace: "abc"},
		},
		{"should read user settings from default workspace",
			args{workspace: "", storedWorkspaceYaml: "\ncurrent: default\nworkspaces:\n    default:\n        jiraRestUrl: https://test2\n        jiraToken: 333\n        jiraUsername: test_user2"},
			&fjiraWorkspaceSettings{JiraToken: "333", JiraUsername: "test_user2", JiraRestUrl: "https://test2", Workspace: "default"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Create(tempDir + "/.fjira/fjira.yaml") //nolint:errcheck
			defer file.Close()
			if err != nil {
				panic(err)
			}
			_, err = file.WriteString(tt.args.storedWorkspaceYaml) //nolint:errcheck
			if err != nil {
				panic(err)
			}
			assert2.Nil(t, err)

			got, err := readFromUserSettings(tt.args.workspace)
			assert2.Nil(t, err)
			assert2.Equalf(t, tt.want, got, "readFromUserSettings(%v)", tt.args.workspace)
		})
	}
}

func Test_readFromUserInputAndWorkspaceEdit(t *testing.T) {
	tempDir := t.TempDir()
	_ = os2.SetUserHomeDir(tempDir)
	_ = os.Mkdir(tempDir+"/.fjira", os.ModePerm) //nolint:errcheck
	t.Cleanup(func() {
		_ = os.RemoveAll(tempDir)
	})

	type args struct {
		workspace        string
		existingSettings *fjiraWorkspaceSettings
	}
	tests := []struct {
		name string
		args args
	}{
		{"should read workspace data from user input",
			args{workspace: "xyz", existingSettings: &fjiraWorkspaceSettings{JiraToken: "123", JiraUsername: "test@test.pl", JiraRestUrl: "https://test.atlassian.net"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			stdin := bytes.NewBufferString("TestUser\nTestUrl\nTestToken\n1\n")

			// when
			settings, err := readFromUserInputAndStore(stdin, tt.args.workspace, tt.args.existingSettings)
			if err != nil {
				assert2.Fail(t, err.Error())
			}
			<-time.NewTimer(100 * time.Millisecond).C

			// then
			assert2.Equal(t, "TestToken", settings.JiraToken)
			assert2.Equal(t, "TestUrl", settings.JiraRestUrl)
			assert2.Equal(t, "TestUser", settings.JiraUsername)
			assert2.Equal(t, jira.ApiToken, settings.JiraTokenType)
			assert2.FileExists(t, fmt.Sprintf(tempDir+"/.fjira/fjira.yaml"))

			// and when
			stdin = bytes.NewBufferString("TestUser2\nTestUrl2\n\n")
			settings, err = readFromWorkspaceEdit(stdin, tt.args.workspace)
			if err != nil {
				assert2.Fail(t, err.Error())
			}
			<-time.NewTimer(100 * time.Millisecond).C

			// then
			assert2.Equal(t, "TestToken", settings.JiraToken)
			assert2.Equal(t, "TestUrl2", settings.JiraRestUrl)
			assert2.Equal(t, "TestUser2", settings.JiraUsername)
			assert2.Equal(t, jira.ApiToken, settings.JiraTokenType)
			assert2.FileExists(t, fmt.Sprintf(tempDir+"/.fjira/fjira.yaml"))
		})
	}
}

func Test_fjira_ValidateWorkspaceName(t *testing.T) {
	type args struct {
		workspace string
	}
	tests := []struct {
		name string
		args args
		wont error
	}{
		{"should validate workspace name", args{workspace: "xyz"}, nil},
		{"should validate workspace name", args{workspace: ""}, nil},
		{"should validate workspace name", args{workspace: ";asd;231"}, WorkspaceFormatInvalidErr},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			resultErr := validateWorkspaceName(tt.args.workspace)

			// then
			assert2.Equal(t, tt.wont, resultErr)
		})
	}
}
