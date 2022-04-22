package fjira

import (
	"errors"
	"os"
)

const (
	JiraTokenEnv    = "FJIRA_TOKEN"
	JiraUsernameEnv = "FJIRA_USERNAME"
	JiraRestUrlEnv  = "FJIRA_REST_URL"
)

var EnvironmentsMissingErr = errors.New("Cannot find " + JiraTokenEnv + " or " + JiraUsernameEnv + " or " + JiraRestUrlEnv + " environments. Please add them in order to use Jira REST API.")

func checkJiraEnvironments() error {
	var token = os.Getenv(JiraTokenEnv)
	var restUrl = os.Getenv(JiraRestUrlEnv)
	var username = os.Getenv(JiraUsernameEnv)
	if token == "" || restUrl == "" || username == "" {
		return EnvironmentsMissingErr
	}
	return nil
}
