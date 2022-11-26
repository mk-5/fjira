package fjira

import (
	"github.com/stretchr/testify/assert"
	"os"
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
		{"should write settings without error", args{workspace: "test2"}},
		{"should write settings without error", args{workspace: "test3"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			tempDir := t.TempDir()
			_ = os.Setenv("HOME", tempDir) // TODO - will not work on windows
			s := &userHomeSettingsStorage{}
			settings := &fjiraSettings{JiraRestUrl: "http://test", JiraUsername: "test_user", JiraToken: "test_token"}
			filepath, _ := s.settingsFilePath(tt.args.workspace)
			assert.NoFileExists(t, filepath)

			// when
			err := s.write(tt.args.workspace, settings)

			// then
			assert.Nil(t, err)
			assert.FileExists(t, filepath)
		})
	}
}
