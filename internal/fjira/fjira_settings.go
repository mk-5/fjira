package fjira

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type fjiraSettings struct {
	JiraRestUrl  string `json:"jiraRestUrl"`
	JiraToken    string `json:"jiraToken"`
	JiraUsername string `json:"jiraUsername"`
}

var (
	WorkspaceNotFoundErr = errors.New("workspace not initialized")
)

const (
	EmptyWorkspace       = ""
	DefaultWorkspaceName = "default"
)

type userHomeSettingsStorage struct{}

type settingsStorage interface { //nolint
	write(workspace string, settings *fjiraSettings) error
}

func (s *userHomeSettingsStorage) read(workspace string) (*fjiraSettings, error) {
	settingsFilePath, err := s.settingsFilePath(workspace)
	if _, err := os.Stat(settingsFilePath); errors.Is(err, os.ErrNotExist) {
		return nil, WorkspaceNotFoundErr
	}
	if err != nil {
		return nil, err
	}
	fileBytes, err := os.ReadFile(settingsFilePath)
	if err != nil {
		return nil, err
	}
	var settings fjiraSettings
	err = json.Unmarshal(fileBytes, &settings)
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func (s *userHomeSettingsStorage) write(workspace string, settings *fjiraSettings) error {
	settingsFilePath, err := s.settingsFilePath(workspace)
	if err != nil {
		return err
	}
	settingsJson, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	err = os.WriteFile(settingsFilePath, settingsJson, 0644)
	return err
}

func (s *userHomeSettingsStorage) settingsFilePath(workspace string) (string, error) {
	if workspace == EmptyWorkspace {
		workspace = DefaultWorkspaceName
	}
	settingsFilename := fmt.Sprintf("%s.json", workspace)
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
	return fmt.Sprintf("%s/%s", configDir, settingsFilename), err
}
