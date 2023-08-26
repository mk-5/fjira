package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterGoto(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should register, and run valid goTo function"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			done := make(chan interface{})
			RegisterGoto("test-screen", func(args ...interface{}) {
				done <- args[0]
			})
			RegisterGoto("test-screen2", func(args ...interface{}) {
			})

			// when
			go GoTo("test-screen", "test-arg")

			// then
			arg1 := <-done
			assert.Equal(t, "test-arg", arg1)
			assert.Equal(t, "test-screen", CurrentScreenName())

			// and when: move forward
			GoTo("test-screen2")

			// then
			assert.Equal(t, "test-screen2", CurrentScreenName())
			assert.Equal(t, "test-screen", PreviousScreenName())

			// and when: go back
			go GoBack()

			// then
			<-done
			assert.Equal(t, "test-screen", CurrentScreenName())
			assert.Equal(t, "test-screen2", PreviousScreenName())
		})
	}
}
