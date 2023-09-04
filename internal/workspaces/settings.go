package workspaces

import (
	"errors"
	"fmt"
	"github.com/mk-5/fjira/internal/jira"
	"gopkg.in/yaml.v3"
	"os"
)

type Settings struct {
	Current    string `json:"current" yaml:"current"`
	Workspaces map[string]WorkspaceSettings
}

type WorkspaceSettings struct {
	JiraRestUrl   string             `json:"jiraRestUrl" yaml:"jiraRestUrl"`
	JiraToken     string             `json:"jiraToken" yaml:"jiraToken"`
	JiraUsername  string             `json:"jiraUsername" yaml:"jiraUsername"`
	JiraTokenType jira.JiraTokenType `json:"jiraTokenType" yaml:"jiraTokenType"`
	Workspace     string             `json:"-" yaml:"-"`
}

type SettingsStorage interface { //nolint
	Write(workspace string, settings *WorkspaceSettings) error
	Read(workspace string) (*WorkspaceSettings, error)
	ReadAllWorkspaces() ([]string, error)
	ReadCurrentWorkspace() (string, error)
	SetCurrentWorkspace(workspace string) error
	ConfigDir() (string, error)
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

func NewUserHomeSettingsStorage() SettingsStorage {
	return &userHomeSettingsStorage{}
}

func (s *userHomeSettingsStorage) Read(workspace string) (*WorkspaceSettings, error) {
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

func (s *userHomeSettingsStorage) Write(workspace string, workspaceSettings *WorkspaceSettings) error {
	settings, err := s.createOrGetSettings()
	if err != nil {
		return err
	}
	settings.Workspaces[workspace] = *workspaceSettings
	err = s.writeSettings(settings)
	return err
}

func (s *userHomeSettingsStorage) writeSettings(settings *Settings) error {
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

func (s *userHomeSettingsStorage) createOrGetSettings() (*Settings, error) {
	settingsFilePath, err := s.settingsFilePath()
	if err != nil {
		return nil, err
	}
	var settings Settings
	settingsBytes, err := os.ReadFile(settingsFilePath)
	if errors.Is(err, os.ErrNotExist) {
		settings = Settings{
			Current:    DefaultWorkspaceName,
			Workspaces: map[string]WorkspaceSettings{},
		}
	} else if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(settingsBytes, &settings)
	if err != nil {
		return nil, err
	}
	// temporary for migration, "" was a default workspace before. Should be removed after some time
	for k := range settings.Workspaces {
		if k == "" {
			if settings.Current == "" {
				settings.Current = DefaultWorkspaceName
			}
			settings.Workspaces[DefaultWorkspaceName] = settings.Workspaces[k]
			delete(settings.Workspaces, k)
			_ = s.writeSettings(&settings)
			break
		}
	}
	return &settings, nil
}

func (s *userHomeSettingsStorage) ReadCurrentWorkspace() (string, error) {
	settings, err := s.createOrGetSettings()
	if err != nil {
		return "", err
	}
	return settings.Current, nil
}

func (s *userHomeSettingsStorage) SetCurrentWorkspace(workspace string) error {
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

func (s *userHomeSettingsStorage) ReadAllWorkspaces() ([]string, error) {
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

func (s *userHomeSettingsStorage) ConfigDir() (string, error) {
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

func (s *userHomeSettingsStorage) settingsFilePath() (string, error) {
	configDir, err := s.ConfigDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", configDir, SettingsFilename), err
}
