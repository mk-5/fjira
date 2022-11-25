package fjira

import (
	"fmt"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
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

func TestFjira_Run(t *testing.T) {
	type fields struct {
		app       *app.App
		api       jira.Api
		formatter FjiraFormatter
		jiraUrl   string
	}
	type args struct {
		args *CliArgs
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fjira{
				app:       tt.fields.app,
				api:       tt.fields.api,
				formatter: tt.fields.formatter,
				jiraUrl:   tt.fields.jiraUrl,
			}
			f.Run(tt.args.args)
		})
	}
}

func TestFjira_SetApi(t *testing.T) {
	type fields struct {
		app       *app.App
		api       jira.Api
		formatter FjiraFormatter
		jiraUrl   string
	}
	type args struct {
		api jira.Api
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fjira{
				app:       tt.fields.app,
				api:       tt.fields.api,
				formatter: tt.fields.formatter,
				jiraUrl:   tt.fields.jiraUrl,
			}
			f.SetApi(tt.args.api)
		})
	}
}

func TestFjira_bootstrap(t *testing.T) {
	type fields struct {
		app       *app.App
		api       jira.Api
		formatter FjiraFormatter
		jiraUrl   string
	}
	type args struct {
		args *CliArgs
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fjira{
				app:       tt.fields.app,
				api:       tt.fields.api,
				formatter: tt.fields.formatter,
				jiraUrl:   tt.fields.jiraUrl,
			}
			f.bootstrap(tt.args.args)
		})
	}
}

func TestGetApi(t *testing.T) {
	tests := []struct {
		name    string
		want    jira.Api
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetApi()
			if !tt.wantErr(t, err, fmt.Sprintf("GetApi()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetApi()")
		})
	}
}

func TestGetFormatter(t *testing.T) {
	tests := []struct {
		name    string
		want    FjiraFormatter
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFormatter()
			if !tt.wantErr(t, err, fmt.Sprintf("GetFormatter()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetFormatter()")
		})
	}
}

func TestGetJiraUrl(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetJiraUrl()
			if !tt.wantErr(t, err, fmt.Sprintf("GetJiraUrl()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetJiraUrl()")
		})
	}
}

func TestInstall(t *testing.T) {
	type args struct {
		args CliArgs
	}
	tests := []struct {
		name    string
		args    args
		want    *fjiraSettings
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Install(tt.args.args)
			if !tt.wantErr(t, err, fmt.Sprintf("Install(%v)", tt.args.args)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Install(%v)", tt.args.args)
		})
	}
}

func TestSetApi(t *testing.T) {
	type args struct {
		api jira.Api
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, SetApi(tt.args.api), fmt.Sprintf("SetApi(%v)", tt.args.api))
		})
	}
}
