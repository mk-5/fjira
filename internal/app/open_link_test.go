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

func TestOpenLink_forSystem(t *testing.T) {
	type args struct {
		system string
	}
	tests := []struct {
		name string
		args args
	}{
		{"should run OpenLink func for system", args{system: "windows"}},
		{"should run OpenLink func for system", args{system: "linux"}},
		{"should run OpenLink func for system", args{system: "darwin"}},
		{"should run OpenLink func for system", args{system: "unknown"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = openLinkForSystem(tt.args.system, "/dev/null")
		})
	}
}
