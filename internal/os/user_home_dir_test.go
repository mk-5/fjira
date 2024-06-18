package os

import (
	"os"
	"testing"
)

func TestMustGetUserHomeDir_xdxPath(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should return XDX path if available"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv("XDX_CONFIG_HOME", "abc")
			if got := MustGetUserHomeDir(); got != "abc" {
				t.Errorf("MustGetUserHomeDir() = %v, want %v", got, "abc")
			}
			_ = os.Setenv("XDX_CONFIG_HOME", "")
		})
	}
}
