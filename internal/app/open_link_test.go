package app

import "testing"

func TestOpenLink(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should run OpenLink func without error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			OpenLink("/dev/null")
		})
	}
}
