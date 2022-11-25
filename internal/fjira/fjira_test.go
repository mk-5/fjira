package fjira

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateNewFjira(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should create&run fjira without error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			f := CreateNewFjira(&fjiraSettings{JiraRestUrl: "test", JiraToken: "test", JiraUsername: "test"})

			// then
			assert.NotNil(t, f)

			// and then
			go f.Run(&CliArgs{})
			<-time.NewTimer(100 * time.Millisecond).C

			// and then
			f.Close()
		})
	}
}
