package app

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	os2 "github.com/mk-5/fjira/internal/os"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestColor(t *testing.T) {
	MustLoadColorScheme()

	type args struct {
		c string
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		{"should get color from default colors", args{c: "navigation.top.background"}, 6260575},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Color(tt.args.c).Hex(), "Color(%v)", tt.args.c)
		})
	}
}

func TestMustLoadColorScheme(t *testing.T) {
	tempDir := t.TempDir()
	_ = os2.SetUserHomeDir(tempDir)
	_ = os.MkdirAll(fmt.Sprintf("%s/.fjira", tempDir), os.ModePerm)
	t.Cleanup(func() {
		_ = os.RemoveAll(tempDir)
	})

	tests := []struct {
		name string
	}{
		{"should load color scheme from user directory"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			f := `
navigation:
  top:
    background: "#0000FF"
    foreground1: "#EEEEEE"

`
			p := fmt.Sprintf("%s/.fjira/colors.yml", tempDir)
			_ = os.WriteFile(p, []byte(f), 0644)
			_, err := os.Stat(p)
			assert.False(t, os.IsNotExist(err))

			// when
			schemeMap = map[string]interface{}{}
			colorsMap = map[string]tcell.Color{}
			MustLoadColorScheme()

			// then
			assert.Equalf(t, int32(255), Color("navigation.top.background").Hex(), "MustLoadColorScheme()")
			assert.Equalf(t, int32(15658734), Color("navigation.top.foreground1").Hex(), "MustLoadColorScheme()")
		})
	}
}
