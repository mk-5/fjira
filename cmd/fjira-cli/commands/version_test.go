package commands

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetVersionCmd(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should create&execute VersionCmd"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			cmd := GetVersionCmd("abc")

			// then
			assert.NotNil(t, cmd)

			// and when
			err := cmd.Execute()

			// then
			assert.Nil(t, err)
		})
	}
}
