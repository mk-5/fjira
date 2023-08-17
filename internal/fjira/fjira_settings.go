package fjira

import (
	"errors"
	"fmt"
	"github.com/mk-5/fjira/internal/jira"
	"gopkg.in/yaml.v3"
	"os"
)

type fjiraSettings struct {
	Current    string `json:"current" yaml:"current"`
	Workspaces map[string]fjiraWorkspaceSettings
}

type fjiraWorkspaceSettings struct {
	JiraRestUrl   string             `json:"jiraRestUrl" yaml:"jiraRestUrl"`
	JiraToken     string             `json:"jiraToken" yaml:"jiraToken"`
	JiraUsername  string             `json:"jiraUsername" yaml:"jiraUsername"`
	JiraTokenType jira.JiraTokenType `json:"jiraTokenType" yaml:"jiraTokenType"`
	Workspace     string             `json:"-" yaml:"-"`
}

var (
	WorkspaceNotFoundErr = errors.New("workspace doesn't exist")
)

const (
	EmptyWorkspace       = ""
	DefaultWorkspaceName = "default"
	SettingsFilename     = "fjira.yaml"
)

type userHomeSettingsStorage struct{}

type settingsStorage interface { //nolint
	write(workspace string, settings *fjiraWorkspaceSettings) error
	read(workspace string) (*fjiraWorkspaceSettings, error)
	readAllWorkspaces() ([]string, error)
	readCurrentWorkspace() (string, error)
	setCurrentWorkspace(workspace string) error
}

func (s *userHomeSettingsStorage) read(workspace string) (*fjiraWorkspaceSettings, error) {
	settings, err := s.createOrGetSettings()
	if err != nil {
		return nil, err
	}
	if w, ok := settings.Workspaces[workspace]; ok {
		w.Workspace = workspace
		return &w, nil
	}
	return nil, WorkspaceNotFoundErr
}

func (s *userHomeSettingsStorage) write(workspace string, workspaceSettings *fjiraWorkspaceSettings) error {
	settings, err := s.createOrGetSettings()
	if err != nil {
		return err
	}
	settings.Workspaces[workspace] = *workspaceSettings
	err = s.writeSettings(settings)
	return err
}

func (s *userHomeSettingsStorage) writeSettings(settings *fjiraSettings) error {
	settingsFilePath, err := s.settingsFilePath()
	if err != nil {
		return err
	}
	settingsYml, err := yaml.Marshal(settings)
	if err != nil {
		return err
	}
	err = os.WriteFile(settingsFilePath, settingsYml, 0644)
	return err
}

func (s *userHomeSettingsStorage) createOrGetSettings() (*fjiraSettings, error) {
	settingsFilePath, err := s.settingsFilePath()
	if err != nil {
		return nil, err
	}
	var settings fjiraSettings
	settingsBytes, err := os.ReadFile(settingsFilePath)
	if errors.Is(err, os.ErrNotExist) {
		settings = fjiraSettings{
			Current:    "",
			Workspaces: map[string]fjiraWorkspaceSettings{},
		}
	} else if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(settingsBytes, &settings)
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func (s *userHomeSettingsStorage) readCurrentWorkspace() (string, error) {
	settings, err := s.createOrGetSettings()
	if err != nil {
		return "", err
	}
	return settings.Current, nil
}

func (s *userHomeSettingsStorage) setCurrentWorkspace(workspace string) error {
	settings, err := s.createOrGetSettings()
	if err != nil {
		return err
	}
	if _, ok := settings.Workspaces[workspace]; ok {
		settings.Current = workspace
		return s.writeSettings(settings)
	}
	return WorkspaceNotFoundErr
}

func (s *userHomeSettingsStorage) readAllWorkspaces() ([]string, error) {
	settings, err := s.createOrGetSettings()
	if err != nil {
		return nil, err
	}
	w := make([]string, 0, len(settings.Workspaces))
	for k := range settings.Workspaces {
		w = append(w, k)
	}
	return w, nil
}

func (s *userHomeSettingsStorage) settingsFilePath() (string, error) {
	configDir, err := s.configDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", configDir, SettingsFilename), err
}

func (s *userHomeSettingsStorage) configDir() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := fmt.Sprintf("%s/.fjira", userHomeDir)
	if _, err := os.Stat(configDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(configDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	return configDir, nil
}
