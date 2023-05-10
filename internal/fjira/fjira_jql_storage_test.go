package fjira

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func Test_jqlStorage(t *testing.T) {
	type args struct {
		jql string
	}
	tests := []struct {
		name string
		args args
	}{
		{"should add&remove jql storage without error", args{"test jql"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			app.CreateNewAppWithScreen(tcell.NewSimulationScreen("utf-8"))
			CreateNewFjira(&fjiraSettings{Workspace: "default"})
			tempDir := t.TempDir()
			_ = os.Setenv("HOME", tempDir) // TODO - will not work on windows
			s := &jqlStorage{}
			jqlFile, _ := s.jqlsFile()

			// when
			err := s.addNew(tt.args.jql)
			bytesRead, err2 := ioutil.ReadFile(jqlFile.Name())
			fileContent := string(bytesRead)
			workspace, errWorkspace := GetCurrentWorkspace()

			// then
			assert.Nil(t, err)
			assert.Nil(t, err2)
			assert.Nil(t, errWorkspace)
			assert.Equal(t, "default", workspace)
			assert.True(t, strings.HasSuffix(jqlFile.Name(), "default.jqls"), "invalid file %s", jqlFile.Name())
			assert.Equal(t, "test jql", fileContent)

			// and when
			err3 := s.remove(tt.args.jql)
			bytesRead2, err4 := ioutil.ReadFile(jqlFile.Name())
			fileContent2 := string(bytesRead2)

			// then
			assert.Nil(t, err3)
			assert.Nil(t, err4)
			assert.NotContains(t, fileContent2, "test jql")
		})
	}
}
