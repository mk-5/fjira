package fjira

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestCreateNewFjira(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should create&run fjira without error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			f := CreateNewFjira(&fjiraSettings{JiraRestUrl: "test", JiraToken: "test", JiraUsername: "test"})

			// then
			assert.NotNil(t, f)

			// and then
			settings, err := Install(CliArgs{})
			go f.Run(&CliArgs{})
			<-time.NewTimer(100 * time.Millisecond).C

			// and then
			f.Close()
			assert.NotNil(t, settings)
			assert.Nil(t, err)
		})
	}
}

func TestInstall(t *testing.T) {
	// TODO - not working on windows
	tempDir := t.TempDir()
	_ = os.Setenv("HOME", tempDir)
	_ = os.Mkdir(tempDir+"/.fjira", os.ModePerm) //nolint:errcheck
	defer func(name string) {
		_ = os.Remove(name)
	}(tempDir + "/.fjira")

	tests := []struct {
		name string
	}{
		{"should run install without error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			u := &userHomeWorkspaces{}
			s := &userHomeSettingsStorage{}
			settings := &fjiraSettings{JiraRestUrl: "http://test", JiraUsername: "test_user", JiraToken: "test_token"}
			filepath, _ := s.settingsFilePath("xyz")
			assert.NoFileExists(t, filepath)

			// when
			_ = s.write("xyz", settings)
			_ = u.setCurrentWorkspace("xyz")
			settings, err := Install(CliArgs{})

			// when
			assert.NotNil(t, settings)
			assert.Nil(t, err)
		})
	}
}
