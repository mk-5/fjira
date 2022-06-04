package fjira

import (
	"errors"
	"fmt"
	"os"
)

type fjiraSettings struct {
	JiraRestUrl    string `json:"jiraRestUrl"`
	JiraBasicToken string `json:"jiraBasicToken"`
}

var WorkspaceNotFoundErr = errors.New("Workspace not initialized.")

const (
	DefaultWorkspace = ""
)

type userHomeSettingsStorage struct{}

type settingsStorage interface { //nolint
	write(workspace string, settings *fjiraSettings) error
	read(workspace string) (*fjiraSettings, error)
}

func (s *userHomeSettingsStorage) read(workspace string) (*fjiraSettings, error) {
	return nil, nil
}

func (s *userHomeSettingsStorage) write(workspace string, settings *fjiraSettings) error {
	if workspace != DefaultWorkspace {
		workspace = fmt.Sprintf("_%s", workspace)
	}
	//settingsFilePath, err := s.settingsFilePath(workspace)
	//if err != nil {
	//	return err
	//}
	//settingsJson, err := json.Marshal(settings)

	//err = os.WriteFile(settingsFilePath, settingsJson, 0644)
	return nil
}

func (s *userHomeSettingsStorage) settingsFilePath(workspace string) (string, error) {
	settingsFilename := fmt.Sprintf("fjira%s.json", workspace)
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
