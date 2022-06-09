package fjira

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_userHomeWorkspaces_getWorkspaceFilepath(t *testing.T) {
	// TODO - it's not multi-platform
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Remove(tempDir + "/.fjira")

	type args struct {
		workspace string
		current   bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"should convert filename into filepath, not current", args{workspace: "test", current: false}, tempDir + "/.fjira/test.json"},
		{"should convert filename into filepath, current", args{workspace: "test", current: true}, tempDir + "/.fjira/_test.json"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userHomeWorkspaces{}
			assert.Equalf(t, tt.want, u.getWorkspaceFilepath(tt.args.workspace, tt.args.current), "getWorkspaceFilepath(%v, %v)", tt.args.workspace, tt.args.current)
		})
	}
}

func Test_userHomeWorkspaces_normalizeWorkspaceFilename(t *testing.T) {
	type args struct {
		workspace string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"should normalize not current workspace path", args{workspace: "/tmp/.fjira/test.json"}, "test"},
		{"should normalize current workspace path", args{workspace: "/tmp/.fjira/_test.json"}, "test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userHomeWorkspaces{}
			assert.Equalf(t, tt.want, u.normalizeWorkspaceFilename(tt.args.workspace), "normalizeWorkspaceFilename(%v)", tt.args.workspace)
		})
	}
}

func Test_userHomeWorkspaces_readAllWorkspaces(t *testing.T) {
	// TODO - it's not multi-platform
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	os.Mkdir(tempDir+"/.fjira", os.ModePerm)  //nolint:errcheck
	os.Create(tempDir + "/.fjira/test1.json") //nolint:errcheck
	os.Create(tempDir + "/.fjira/test2.json") //nolint:errcheck
	os.Create(tempDir + "/.fjira/test3.json") //nolint:errcheck
	defer os.Remove(tempDir + "/.fjira")

	tests := []struct {
		name string
		want []string
	}{
		{"should read all workspaces", []string{"test1", "test2", "test3"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userHomeWorkspaces{}
			got, _ := u.readAllWorkspaces()
			assert.Equalf(t, tt.want, got, "readAllWorkspaces()")
		})
	}
}

func Test_userHomeWorkspaces_readCurrentWorkspace(t *testing.T) {
	// TODO - it's not multi-platform
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	os.Mkdir(tempDir+"/.fjira", os.ModePerm) //nolint:errcheck
	os.Create(tempDir + "/.fjira/_xyz.json") //nolint:errcheck
	defer os.Remove(tempDir + "/.fjira")

	tests := []struct {
		name    string
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{"should read current workspace", "xyz", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userHomeWorkspaces{}
			got, err := u.readCurrentWorkspace()
			if tt.wantErr != nil && !tt.wantErr(t, err, "readCurrentWorkspace()") {
				return
			}
			assert.Equalf(t, tt.want, got, "readCurrentWorkspace()")
		})
	}
}

func Test_userHomeWorkspaces_setCurrentWorkspace(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	os.Mkdir(tempDir+"/.fjira", os.ModePerm)     //nolint:errcheck
	os.Create(tempDir + "/.fjira/_default.json") //nolint:errcheck
	os.Create(tempDir + "/.fjira/yyy.json")      //nolint:errcheck
	defer os.Remove(tempDir + "/.fjira")

	type args struct {
		workspace string
	}
	tests := []struct {
		name string
		args args
	}{
		{"should set current workspace", args{workspace: "yyy"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userHomeWorkspaces{}
			err := u.setCurrentWorkspace(tt.args.workspace)
			_, err2 := os.Stat(tempDir + "/.fjira/_yyy.json")
			_, err3 := os.Stat(tempDir + "/.fjira/_default.json")

			assert.Nil(t, err)
			assert.Nil(t, err2)
			assert.ErrorIs(t, err3, os.ErrNotExist)
		})
	}
}
