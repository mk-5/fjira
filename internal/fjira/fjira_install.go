package fjira

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"io"
	"os"
	"regexp"
	"strings"
)

const (
	JiraTokenEnv    = "FJIRA_TOKEN"
	JiraUsernameEnv = "FJIRA_USERNAME"
	JiraRestUrlEnv  = "FJIRA_REST_URL"
)

var (
	EnvironmentsMissingErr    = errors.New("cannot find " + JiraTokenEnv + " or " + JiraUsernameEnv + " or " + JiraRestUrlEnv + " environments. Please add them in order to use Jira REST API")
	WorkspaceFormatInvalidErr = errors.New("workspace name needs to match pattern [a-z0-9]{2,50}")
	workspaceRegExp           = regexp.MustCompile("^[a-z0-9]{2,50}$")
)

func validateWorkspaceName(workspace string) error {
	if workspace != EmptyWorkspace && !workspaceRegExp.MatchString(workspace) {
		return WorkspaceFormatInvalidErr
	}
	return nil
}

func readFromEnvironments() (*fjiraSettings, error) {
	var token = os.Getenv(JiraTokenEnv)
	var restUrl = os.Getenv(JiraRestUrlEnv)
	var username = os.Getenv(JiraUsernameEnv)
	if token == "" || restUrl == "" || username == "" {
		return nil, EnvironmentsMissingErr
	}
	return &fjiraSettings{
		JiraToken:    token,
		JiraRestUrl:  restUrl,
		JiraUsername: username,
	}, nil
}

func readFromUserSettings(workspace string) (*fjiraSettings, error) {
	var err error
	if workspace == EmptyWorkspace {
		workspaces := &userHomeWorkspaces{}
		workspace, err = workspaces.readCurrentWorkspace()
	}
	if err != nil {
		return nil, err
	}
	var settingsStorage = &userHomeSettingsStorage{}
	settings, err := settingsStorage.read(workspace)
	if err != nil {
		return nil, err
	}
	return settings, err
}

func readFromUserInputAndStore(input io.Reader, workspace string, existingSettings *fjiraSettings) (*fjiraSettings, error) {
	workspaceName := workspace
	if workspace == EmptyWorkspace {
		workspaceName = DefaultWorkspaceName
	}
	fmt.Print("\033[?1049h\033[H") // alternate screen buffer
	defer func() {
		fmt.Print("\033[?1049l")
	}()
	if existingSettings == nil {
		fmt.Print(MessageCreateNewWorkspace)
	} else {
		fmt.Print(MessageEditWorkspace)
	}
	fmt.Println(color.CyanString(workspaceName))
	fmt.Println("")
	reader := bufio.NewReader(input)
	fmt.Print(color.HiYellowString(MessageQuestionMark))
	fmt.Print(MessageEnterUsername)
	if existingSettings != nil {
		fmt.Print(color.BlueString("[%s] ", existingSettings.JiraUsername))
	}
	username, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	fmt.Print(color.HiYellowString(MessageQuestionMark))
	fmt.Print(MessageEnterJiraUrl)
	if existingSettings != nil {
		fmt.Print(color.BlueString("[%s] ", existingSettings.JiraRestUrl))
	} else {
		fmt.Print(color.BlueString(MessageEnterJiraUrlExample))
	}
	url, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	fmt.Print(color.HiYellowString(MessageQuestionMark))
	fmt.Print(MessageEnterJiraApiToken)
	if existingSettings != nil {
		fmt.Print(color.BlueString("[%s] ", existingSettings.JiraToken))
	}
	token, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	settings := &fjiraSettings{
		JiraToken:    strings.TrimSpace(token),
		JiraUsername: strings.TrimSpace(username),
		JiraRestUrl:  strings.TrimSpace(url),
	}
	if existingSettings != nil && settings.JiraUsername == "" {
		settings.JiraUsername = existingSettings.JiraUsername
	}
	if existingSettings != nil && settings.JiraToken == "" {
		settings.JiraToken = existingSettings.JiraToken
	}
	if existingSettings != nil && settings.JiraRestUrl == "" {
		settings.JiraRestUrl = existingSettings.JiraRestUrl
	}
	var settingsStorage = &userHomeSettingsStorage{}
	err = settingsStorage.write(workspace, settings)
	if err != nil {
		return nil, err
	}
	workspaces := &userHomeWorkspaces{}
	_ = workspaces.setCurrentWorkspace(workspace)
	return settings, err
}

func readFromWorkspaceEdit(input io.Reader, workspace string) (*fjiraSettings, error) {
	var settingsStorage = &userHomeSettingsStorage{}
	settings, err := settingsStorage.read(workspace)
	if err != nil {
		return nil, err
	}
	editedSettings, err := readFromUserInputAndStore(input, workspace, settings)
	if err != nil {
		return nil, err
	}
	err = settingsStorage.write(workspace, editedSettings)
	if err != nil {
		return nil, err
	}
	return editedSettings, nil
}
