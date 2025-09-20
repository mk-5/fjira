package workspaces

import (
	"errors"
	os2 "github.com/mk-5/fjira/internal/os"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_userHomeSettingsStorage_write(t *testing.T) {
	type args struct {
		workspace string
	}
	tests := []struct {
		name string
		args args
	}{
		{"should Write settings without error", args{workspace: "test2"}},
		{"should Write settings without error", args{workspace: "test3"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			tempDir := t.TempDir()
			_ = os2.SetUserHomeDir(tempDir)
			s := &userHomeSettingsStorage{}
			settings := &WorkspaceSettings{JiraRestUrl: "http://test", JiraUsername: "test_user", JiraToken: "test_token"}
			filepath, _ := s.settingsFilePath()
			assert.NoFileExists(t, filepath)

			// when
			err := s.Write(tt.args.workspace, settings)

			// then
			assert.Nil(t, err)
			assert.FileExists(t, filepath)
		})
	}
}

func Test_userHomeSettingsStorage_read(t *testing.T) {
	type args struct {
		workspace string
	}
	tests := []struct {
		name string
		args args
	}{
		{"should return ErrWorkspaceNotFound if workspace doesn't exit", args{workspace: "test2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			tempDir := t.TempDir()
			_ = os2.SetUserHomeDir(tempDir)
			s := &userHomeSettingsStorage{}
			filepath, _ := s.settingsFilePath()
			assert.NoFileExists(t, filepath)

			// when
			_, err := s.Read(tt.args.workspace)

			// then
			assert.True(t, errors.Is(err, ErrWorkspaceNotFound))
		})
	}
}
