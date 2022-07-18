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
	CurrentWorkspaceFilePattern    = "%s/.fjira/_current.json"
	AvailableWorkspacesPattern     = "%s/.fjira/[^_]*.json"
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
	linkPath := fmt.Sprintf(CurrentWorkspaceFilePattern, userHomeDir)
	workspaceFilePath, err := os.Readlink(linkPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	if workspaceFilePath == "" {
		return DefaultWorkspaceName, nil
	}
	workspace := u.normalizeWorkspaceFilename(workspaceFilePath)
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
		if normalized == "" {
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
	workspaceFilepath := u.getWorkspaceFilepath(workspace, false)
	if _, err := os.Stat(workspaceFilepath); errors.Is(err, os.ErrNotExist) {
		return WorkspaceNotFoundErr
	}
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err.Error())
	}
	currentWorkspacePath := fmt.Sprintf(CurrentWorkspaceFilePattern, userHomeDir)
	if _, err := os.Lstat(currentWorkspacePath); err == nil {
		_ = os.Remove(currentWorkspacePath)
	}
	err = os.Symlink(workspaceFilepath, currentWorkspacePath)
	if err != nil {
		return err
	}
	return err
}

func (*userHomeWorkspaces) getWorkspaceFilepath(workspace string, current bool) string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err.Error())
	}
	return fmt.Sprintf("%s/.fjira/%s.json", userHomeDir, workspace)
}

func (*userHomeWorkspaces) normalizeWorkspaceFilename(workspace string) string {
	workspace = filepath.Base(workspace)
	workspace = strings.TrimSpace(workspace)
	workspace = strings.Replace(workspace, CurrentWorkspaceFilenamePrefix, "", 1)
	workspace = strings.Replace(workspace, WorkspaceFileExtension, "", 1)
	workspace = strings.Join(strings.Fields(workspace), "")
	return workspace
}
