package fjira

import (
	os2 "github.com/mk-5/fjira/internal/os"
	"github.com/stretchr/testify/assert"
	"os"
	"runtime"
	"testing"
)

func Test_userHomeWorkspaces_normalizeWorkspaceFilename(t *testing.T) {
	// userHomeWorkspaces is deprecated and it's not working on windows
	if runtime.GOOS == "windows" {
		return
	}
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
	// userHomeWorkspaces is deprecated and it's not working on windows
	if runtime.GOOS == "windows" {
		return
	}
	tempDir := t.TempDir()
	_ = os2.SetUserHomeDir(tempDir)
	os.Mkdir(tempDir+"/.fjira", os.ModePerm)  //nolint:errcheck
	os.Create(tempDir + "/.fjira/test1.json") //nolint:errcheck
	os.Create(tempDir + "/.fjira/test2.json") //nolint:errcheck
	os.Create(tempDir + "/.fjira/test3.json") //nolint:errcheck
	t.Cleanup(func() {
		_ = os.RemoveAll(tempDir)
	})

	tests := []struct {
		name string
		want []string
	}{
		{"should read all fjiraSettings", []string{"test1", "test2", "test3"}},
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
	// userHomeWorkspaces is deprecated and it's not working on windows
	if runtime.GOOS == "windows" {
		return
	}
	tempDir := t.TempDir()
	_ = os2.SetUserHomeDir(tempDir)
	_ = os.Mkdir(tempDir+"/.fjira", os.ModePerm)                                //nolint:errcheck
	_, _ = os.Create(tempDir + "/.fjira/xyz.json")                              //nolint:errcheck
	_ = os.Symlink(tempDir+"/.fjira/xyz.json", tempDir+"/.fjira/_current.json") //nolint:errcheck
	t.Cleanup(func() {
		_ = os.RemoveAll(tempDir)
	})

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
