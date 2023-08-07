package os

import (
	"os"
	"runtime"
)

func SetUserHomeDir(dir string) error {
	// vars taken from os.UserHomeDir
	env := "HOME"
	switch runtime.GOOS {
	case "windows":
		env = "USERPROFILE"
	case "plan9":
		env = "home"
	}
	return os.Setenv(env, dir)
}
