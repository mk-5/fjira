package fjira

import (
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_shouldReturnErrorWhenNoEnvironments(t *testing.T) {
	// given
	assert := assert2.New(t)
	os.Setenv(JiraTokenEnv, "")
	os.Setenv(JiraUsernameEnv, "")
	os.Setenv(JiraRestUrlEnv, "")

	// when
	_, error := readFromEnvironments()

	// then
	assert.Error(error, "Should return error when no fjira environments")
}

func Test_shouldReturnNoErrorWhenEnvironments(t *testing.T) {
	// given
	assert := assert2.New(t)
	os.Setenv(JiraTokenEnv, "test")
	os.Setenv(JiraUsernameEnv, "test")
	os.Setenv(JiraRestUrlEnv, "http://test.test")

	// when
	_, error := readFromEnvironments()

	// then
	assert.NoError(error, "Should return no error when fjira environments")
}
