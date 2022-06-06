package fjira

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mk5/fjira/internal/app"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSwitchWorkspaceView(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should create new workspace view"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, NewSwitchWorkspaceView(), "NewSwitchWorkspaceView()")
		})
	}
}

func Test_fjiraSwitchWorkspaceView_Draw(t *testing.T) {
	type fields struct {
		fuzzyFind  *app.FuzzyFind
		workspaces *userHomeWorkspaces
	}
	type args struct {
		screen tcell.Screen
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"should run Draw without errors", fields{fuzzyFind: app.NewFuzzyFind("test", []string{})}, args{screen: tcell.NewSimulationScreen("utf-8")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &fjiraSwitchWorkspaceView{
				fuzzyFind:  tt.fields.fuzzyFind,
				workspaces: tt.fields.workspaces,
			}
			s.Draw(tt.args.screen)
		})
	}
}

func Test_fjiraSwitchWorkspaceView_HandleKeyEvent(t *testing.T) {
	type fields struct {
		fuzzyFind  *app.FuzzyFind
		workspaces *userHomeWorkspaces
	}
	type args struct {
		keyEvent *tcell.EventKey
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"should run HandleKeyEvent without errors", fields{fuzzyFind: app.NewFuzzyFind("test", []string{})}, args{keyEvent: &tcell.EventKey{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &fjiraSwitchWorkspaceView{
				fuzzyFind:  tt.fields.fuzzyFind,
				workspaces: tt.fields.workspaces,
			}
			s.HandleKeyEvent(tt.args.keyEvent)
		})
	}
}

func Test_fjiraSwitchWorkspaceView_Init(t *testing.T) {
	type fields struct {
		fuzzyFind  *app.FuzzyFind
		workspaces *userHomeWorkspaces
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"should run Init without errors", fields{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &fjiraSwitchWorkspaceView{
				fuzzyFind:  tt.fields.fuzzyFind,
				workspaces: tt.fields.workspaces,
			}
			s.Init()

			assert.NotNil(t, s.fuzzyFind)
		})
	}
}

func Test_fjiraSwitchWorkspaceView_Resize(t *testing.T) {
	type fields struct {
		fuzzyFind  *app.FuzzyFind
		workspaces *userHomeWorkspaces
	}
	type args struct {
		screenX int
		screenY int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"should run resize without errors", fields{fuzzyFind: app.NewFuzzyFind("test", []string{})}, args{screenX: 10, screenY: 10}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &fjiraSwitchWorkspaceView{
				fuzzyFind:  tt.fields.fuzzyFind,
				workspaces: tt.fields.workspaces,
			}
			s.Resize(tt.args.screenX, tt.args.screenY)
		})
	}
}

func Test_fjiraSwitchWorkspaceView_Update(t *testing.T) {
	type fields struct {
		fuzzyFind  *app.FuzzyFind
		workspaces *userHomeWorkspaces
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"should run update without errors", fields{fuzzyFind: app.NewFuzzyFind("test", []string{})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &fjiraSwitchWorkspaceView{
				fuzzyFind:  tt.fields.fuzzyFind,
				workspaces: tt.fields.workspaces,
			}
			s.Update()
		})
	}
}
