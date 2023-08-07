package fjira

import (
	"bytes"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
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
			fjira := CreateNewFjira(&fjiraSettings{})

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
			&fjiraSettings{JiraToken: "123", JiraUsername: "test@test.pl", JiraRestUrl: "https://test.atlassian.net", Workspace: "xyz"},
		},
		{"should read user settings from another workspace",
			args{workspace: "abc", storedWorkspaceFilename: "abc.json", storedWorkspaceJson: "{\"jiraRestUrl\":\"https://test\",\"jiraToken\":\"111\",\"jiraUsername\":\"test_user\"}"},
			&fjiraSettings{JiraToken: "111", JiraUsername: "test_user", JiraRestUrl: "https://test", Workspace: "abc"},
		},
		{"should read user settings from default workspace",
			args{workspace: "", storedWorkspaceFilename: "default.json", storedWorkspaceJson: "{\"jiraRestUrl\":\"https://test\",\"jiraToken\":\"111\",\"jiraUsername\":\"test_user\"}"},
			&fjiraSettings{JiraToken: "111", JiraUsername: "test_user", JiraRestUrl: "https://test", Workspace: "default"},
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

func Test_readFromUserInputAndWorkspaceEdit(t *testing.T) {
	tempDir := t.TempDir()
	_ = os2.SetUserHomeDir(tempDir)
	_ = os.Mkdir(tempDir+"/.fjira", os.ModePerm) //nolint:errcheck
	defer func(name string) {
		_ = os.Remove(name)
	}(tempDir + "/.fjira")

	type args struct {
		workspace        string
		existingSettings *fjiraSettings
	}
	tests := []struct {
		name string
		args args
	}{
		{"should read workspace data from user input",
			args{workspace: "xyz", existingSettings: &fjiraSettings{JiraToken: "123", JiraUsername: "test@test.pl", JiraRestUrl: "https://test.atlassian.net"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			stdin := bytes.NewBufferString("TestUser\nTestUrl\nTestToken\n")

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
			assert2.FileExists(t, fmt.Sprintf(tempDir+"/.fjira/%s.json", tt.args.workspace))

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
			assert2.FileExists(t, fmt.Sprintf(tempDir+"/.fjira/%s.json", tt.args.workspace))
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
