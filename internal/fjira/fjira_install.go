package fjira

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/fatih/color"
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

func readFromUserInputAndStore(workspace string) (*fjiraSettings, error) {
	workspaceName := workspace
	if workspace == EmptyWorkspace {
		workspaceName = DefaultWorkspaceName
	}
	fmt.Print("\033[?1049h\033[H") // alternate screen buffer
	defer func() {
		fmt.Print("\033[?1049l")
	}()
	fmt.Print(MessageCreatingNewWorkspace)
	fmt.Println(color.CyanString(workspaceName))
	fmt.Println("")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(color.HiYellowString(MessageQuestionMark))
	fmt.Print(MessageEnterUsername)
	username, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	fmt.Print(color.HiYellowString(MessageQuestionMark))
	fmt.Print(MessageEnterJiraUrl)
	fmt.Print(color.BlueString(MessageEnterJiraUrlExample))
	url, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	fmt.Print(color.HiYellowString(MessageQuestionMark))
	fmt.Print(MessageEnterJiraApiToken)
	token, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	settings := &fjiraSettings{
		JiraToken:    strings.TrimSuffix(token, "\n"),
		JiraUsername: strings.TrimSuffix(username, "\n"),
		JiraRestUrl:  strings.TrimSuffix(url, "\n"),
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
