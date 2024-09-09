package os

import (
	"fmt"
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

func MustGetUserHomeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return dir
}

func MustGetFjiraHomeDir() string {
	dir := MustGetUserHomeDir()
	userHomeDir := fmt.Sprintf("%s/.fjira", dir)
	xdg := os.Getenv("XDG_CONFIG_HOME")
	xdgDir := fmt.Sprintf("%s/fjira", xdg)
	if xdg != "" && dirExist(xdgDir) {
		return xdgDir
	}
	if xdg != "" && !dirExist(userHomeDir) {
		return xdgDir
	}
	return userHomeDir
}

func dirExist(path string) bool {
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return true
	}
	return false
}
