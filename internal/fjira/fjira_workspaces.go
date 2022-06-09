package fjira

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	CurrentWorkspaceFilenamePrefix = "_"
	CurrentWorkspaceFilePattern    = "%s/.fjira/_*.json"
	AvailableWorkspacesPattern     = "%s/.fjira/*.json"
	WorkspaceFileExtension         = ".json"
)

type userHomeWorkspaces struct{}

type workspaces interface { //nolint
	readCurrentWorkspace() (string, error)
	readAllWorkspaces() ([]string, error)
	setCurrentWorkspace(workspace string) error
}

func (u *userHomeWorkspaces) readCurrentWorkspace() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	pattern := fmt.Sprintf(CurrentWorkspaceFilePattern, userHomeDir)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return DefaultWorkspaceName, nil
	}
	workspace := u.normalizeWorkspaceFilename(matches[0])
	return workspace, nil
}

func (u *userHomeWorkspaces) readAllWorkspaces() ([]string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	pattern := fmt.Sprintf(AvailableWorkspacesPattern, userHomeDir)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	workspaces := make([]string, 0, len(matches))
	for _, filename := range matches {
		normalized := u.normalizeWorkspaceFilename(filename)
		if strings.TrimSpace(normalized) == "" {
			continue
		}
		workspaces = append(workspaces, normalized)
	}
	return workspaces, nil
}

func (u *userHomeWorkspaces) setCurrentWorkspace(workspace string) error {
	if workspace == EmptyWorkspace {
		workspace = DefaultWorkspaceName
	}
	previousWorkspace, err := u.readCurrentWorkspace()
	if err != nil {
		return err
	}
	if previousWorkspace == workspace {
		return nil
	}
	previousWorkspaceFilepath := u.getWorkspaceFilepath(previousWorkspace, true)
	workspaceFilepath := u.getWorkspaceFilepath(workspace, false)
	if _, err := os.Stat(workspaceFilepath); errors.Is(err, os.ErrNotExist) {
		return WorkspaceNotFoundErr
	}
	err = os.Rename(workspaceFilepath, u.getWorkspaceFilepath(workspace, true))
	if err != nil {
		return err
	}
	err = os.Rename(previousWorkspaceFilepath, u.getWorkspaceFilepath(previousWorkspace, false))
	return err
}

func (u *userHomeWorkspaces) getWorkspaceFilepath(workspace string, current bool) string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err.Error())
	}
	if current {
		return fmt.Sprintf("%s/.fjira/_%s.json", userHomeDir, workspace)
	}
	return fmt.Sprintf("%s/.fjira/%s.json", userHomeDir, workspace)
}

func (u *userHomeWorkspaces) normalizeWorkspaceFilename(workspace string) string {
	workspace = filepath.Base(workspace)
	workspace = strings.Replace(workspace, CurrentWorkspaceFilenamePrefix, "", 1)
	workspace = strings.Replace(workspace, WorkspaceFileExtension, "", 1)
	workspace = strings.Join(strings.Fields(workspace), "")
	return workspace
}
