package os

import (
	"os"
	"testing"
)

func TestMustGetUserHomeDir_xdxPath(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should return XDG path if available"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv("XDG_CONFIG_HOME", "abc")
			_ = SetUserHomeDir("something_for_just_for_test")
			if got := MustGetFjiraHomeDir(); got != "abc/fjira" {
				t.Errorf("MustGetUserHomeDir() = %v, want %v", got, "abc")
			}
			_ = os.Setenv("XDG_CONFIG_HOME", "")
		})
	}
}

func TestMustGetUserHomeDir_userHomeWhenExist(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should return userHome when exist, and XDG doesn't"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv("XDG_CONFIG_HOME", "abc")
			_ = SetUserHomeDir("./something_for_just_for_test")
			_ = os.MkdirAll("./something_for_just_for_test/.fjira", 0750)
			if got := MustGetFjiraHomeDir(); got != "./something_for_just_for_test/.fjira" {
				t.Errorf("MustGetUserHomeDir() = %v, want %v", got, "abc")
			}
			_ = os.Setenv("XDG_CONFIG_HOME", "")
			_ = os.RemoveAll("./something_for_just_for_test")
		})
	}
}

func TestMustGetUserHomeDir_homePath(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should return HOME path if available"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = SetUserHomeDir("user_home")
			if got := MustGetFjiraHomeDir(); got != "user_home/.fjira" {
				t.Errorf("MustGetUserHomeDir() = %v, want %v", got, "abc")
			}
			_ = os.Setenv("HOME", "")
		})
	}
}
