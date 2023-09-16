package fjira

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/mk-5/fjira/internal/jira"
	"github.com/mk-5/fjira/internal/ui"
	"github.com/mk-5/fjira/internal/workspaces"
	goinput "github.com/tcnksm/go-input"
	"io"
	"os"
	"regexp"
	"strings"
)

const (
	JiraTokenEnv    = "FJIRA_TOKEN"
	JiraUsernameEnv = "FJIRA_USERNAME"
	JiraTokenType   = "FJIRA_JIRA_TOKEN_TYPE"
	JiraRestUrlEnv  = "FJIRA_REST_URL"
)

var (
	EnvironmentsMissingErr    = errors.New("cannot find " + JiraTokenEnv + " or " + JiraUsernameEnv + " or " + JiraRestUrlEnv + " environments. Please add them in order to use Jira REST API")
	WorkspaceFormatInvalidErr = errors.New("workspace name needs to match pattern [a-z0-9]{2,50}")
	workspaceRegExp           = regexp.MustCompile("^[a-z0-9]{2,50}$")
)

func Install(workspace string) (*workspaces.WorkspaceSettings, error) {
	// it will be removed after a few version
	u := workspaces.NewDeprecatedUserHomeWorkspaces()
	_ = u.MigrateFromGlobWorkspacesToYaml()

	err := validateWorkspaceName(workspace)
	if err != nil {
		return nil, err
	}
	s, err := readFromEnvironments()
	if err == nil {
		return s, nil // envs found
	}
	if err != EnvironmentsMissingErr {
		return nil, err
	}
	s2, err := readFromUserSettings(workspace)
	if err == workspaces.WorkspaceNotFoundErr || errors.Unwrap(err) == workspaces.WorkspaceNotFoundErr {
		return readFromUserInputAndStore(os.Stdin, workspace, nil)
	}
	if err != nil {
		return nil, err
	}
	return s2, nil
}

func EditWorkspaceAndReadSettings(input io.Reader, workspace string) (*workspaces.WorkspaceSettings, error) {
	var settingsStorage = workspaces.NewUserHomeSettingsStorage()
	settings, err := settingsStorage.Read(workspace)
	if err != nil {
		return nil, err
	}
	editedSettings, err := readFromUserInputAndStore(input, workspace, settings)
	if err != nil {
		return nil, err
	}
	err = settingsStorage.Write(workspace, editedSettings)
	if err != nil {
		return nil, err
	}
	return editedSettings, nil
}

func readFromEnvironments() (*workspaces.WorkspaceSettings, error) {
	var token = os.Getenv(JiraTokenEnv)
	var restUrl = os.Getenv(JiraRestUrlEnv)
	var username = os.Getenv(JiraUsernameEnv)
	var tokenTypeStr = os.Getenv(JiraTokenType)
	if token == "" || restUrl == "" || username == "" {
		return nil, EnvironmentsMissingErr
	}
	if tokenTypeStr == "" {
		tokenTypeStr = string(jira.ApiToken)
	}
	tokenType := jira.JiraTokenType(tokenTypeStr)
	return &workspaces.WorkspaceSettings{
		JiraToken:     token,
		JiraRestUrl:   restUrl,
		JiraUsername:  username,
		JiraTokenType: tokenType,
	}, nil
}

func readFromUserSettings(workspace string) (*workspaces.WorkspaceSettings, error) {
	var err error
	settingsStorage := workspaces.NewUserHomeSettingsStorage()
	if workspace == workspaces.EmptyWorkspace {
		workspace, err = settingsStorage.ReadCurrentWorkspace()
	}
	if err != nil {
		return nil, err
	}
	settings, err := settingsStorage.Read(workspace)
	if err != nil {
		return nil, err
	}
	return settings, err
}

func readFromUserInputAndStore(input io.Reader, workspace string, existingSettings *workspaces.WorkspaceSettings) (*workspaces.WorkspaceSettings, error) {
	workspaceName := workspace
	if workspace == workspaces.EmptyWorkspace {
		workspaceName = workspaces.DefaultWorkspaceName
	}
	fmt.Print("\033[?1049h\033[H") // alternate screen buffer
	defer func() {
		fmt.Print("\033[?1049l")
	}()
	if existingSettings == nil {
		fmt.Print(ui.MessageCreateNewWorkspace)
	} else {
		fmt.Print(ui.MessageEditWorkspace)
	}
	fmt.Println(color.CyanString(workspaceName))
	fmt.Println("")

	reader := bufio.NewReader(input)
	fmt.Print(color.HiYellowString(ui.MessageQuestionMark))
	fmt.Print(ui.MessageEnterUsername)
	if existingSettings != nil {
		fmt.Print(color.BlueString("[%s] ", existingSettings.JiraUsername))
	}
	username, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	fmt.Print(color.HiYellowString(ui.MessageQuestionMark))
	fmt.Print(ui.MessageEnterJiraUrl)
	if existingSettings != nil {
		fmt.Print(color.BlueString("[%s] ", existingSettings.JiraRestUrl))
	} else {
		fmt.Print(color.BlueString(ui.MessageEnterJiraUrlExample))
	}
	url, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	fmt.Print(color.HiYellowString(ui.MessageQuestionMark))
	fmt.Print(ui.MessageEnterJiraApiToken)
	if existingSettings != nil {
		fmt.Print(color.BlueString("[%s] ", existingSettings.JiraToken))
	}
	token, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	fmt.Print(color.HiYellowString(ui.MessageQuestionMark))
	u := &goinput.UI{
		Writer: os.Stdout,
		Reader: input,
	}
	tokenType := jira.ApiToken
	if existingSettings != nil && existingSettings.JiraTokenType != "" {
		tokenType = existingSettings.JiraTokenType
	}
	tokenTypeOptions := []string{string(jira.ApiToken), string(jira.PersonalToken)}
	tokenTypeStr, err := u.Select(ui.MessageEnterJiraTokenType, tokenTypeOptions, &goinput.Options{
		Default:  string(tokenType),
		Required: true,
		Loop:     true,
	})
	if err != nil {
		return nil, err
	}
	settings := &workspaces.WorkspaceSettings{
		JiraToken:     strings.TrimSpace(token),
		JiraUsername:  strings.TrimSpace(username),
		JiraRestUrl:   strings.TrimSpace(url),
		JiraTokenType: jira.JiraTokenType(tokenTypeStr),
	}
	if existingSettings != nil && existingSettings.JiraUsername == "" {
		settings.JiraUsername = existingSettings.JiraUsername
	}
	if existingSettings != nil && settings.JiraToken == "" {
		settings.JiraToken = existingSettings.JiraToken
	}
	if existingSettings != nil && settings.JiraRestUrl == "" {
		settings.JiraRestUrl = existingSettings.JiraRestUrl
	}
	var settingsStorage = workspaces.NewUserHomeSettingsStorage()
	err = settingsStorage.Write(workspace, settings)
	if err != nil {
		return nil, err
	}
	_ = settingsStorage.SetCurrentWorkspace(workspace)
	return settings, err
}

func validateWorkspaceName(workspace string) error {
	if workspace != workspaces.EmptyWorkspace && !workspaceRegExp.MatchString(workspace) {
		return WorkspaceFormatInvalidErr
	}
	return nil
}
