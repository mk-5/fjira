package fjira

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	CurrentWorkspaceFilenamePrefix = "_"
	CurrentWorkspaceFileName       = "current"
	CurrentWorkspaceFilePattern    = "%s/.fjira/" + CurrentWorkspaceFilenamePrefix + CurrentWorkspaceFileName + ".json" // @deprecated
	AvailableWorkspacesPattern     = "%s/.fjira/[^_]*.json"                                                             // @deprecated
	WorkspaceFileExtension         = ".json"                                                                            // @deprecated
)

// @deprecated - it shouldn't be in use. Everything is handled by userHomeSettingsStorage now
type userHomeWorkspaces struct{}

type workspaces interface { //nolint
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

func (*userHomeWorkspaces) normalizeWorkspaceFilename(workspace string) string {
	workspace = filepath.Base(workspace)
	workspace = strings.TrimSpace(workspace)
	workspace = strings.Replace(workspace, CurrentWorkspaceFilenamePrefix, "", 1)
	workspace = strings.Replace(workspace, WorkspaceFileExtension, "", 1)
	workspace = strings.Join(strings.Fields(workspace), "")
	return workspace
}

// In the first version all workspaces have been stored ~/.fjira/ directory,
// and the current workspace pointer was just _current.json file.
// There is a problem with symlinks for windows platform, so it was not super future-proof solution.
// That method is migrating from the old to the new .yml settings approach
func (u *userHomeWorkspaces) migrateFromGlobWorkspacesToYaml() error {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err.Error())
	}
	oldCurrentWorkspacePointerLink := fmt.Sprintf(CurrentWorkspaceFilePattern, userHomeDir)
	if _, err := os.Lstat(oldCurrentWorkspacePointerLink); errors.Is(err, os.ErrNotExist) {
		// nothing to do
		return nil
	}

	workspaces, err := u.readAllWorkspaces()
	if err != nil {
		return err
	}

	settingsStorage := &userHomeSettingsStorage{}

	for _, w := range workspaces {
		file := fmt.Sprintf("%s/.fjira/%s.json", userHomeDir, w)
		bytes, err := os.ReadFile(file)
		if err != nil {
			// skip if it cannot read the workspace
			continue
		}
		var wSettings fjiraWorkspaceSettings
		err = json.Unmarshal(bytes, &wSettings)
		if err != nil {
			// skip if it cannot read the workspace
			continue
		}
		err = settingsStorage.write(w, &wSettings)
		if err != nil {
			return err
		}
		os.Remove(file)
	}

	current, err := u.readCurrentWorkspace()
	if err != nil {
		current = DefaultWorkspaceName
	}
	err = settingsStorage.setCurrentWorkspace(current)
	if err != nil {
		return err
	}

	// remove old files
	os.Remove(fmt.Sprintf("%s/.fjira/_current.json", userHomeDir))
	return nil
}
