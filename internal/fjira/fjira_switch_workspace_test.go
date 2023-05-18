package fjira

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewSwitchWorkspaceView(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should create new workspace view"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := newSwitchWorkspaceView()
			view.Destroy()
			assert.NotNil(t, view, "newSwitchWorkspaceView()")
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

func Test_fjiraSwitchWorkspaceView_should_handle_empty_fuzzy_find_result(t *testing.T) {
	// given
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()
	app.CreateNewAppWithScreen(screen)
	s := newSwitchWorkspaceView()

	// when
	s.Init()
	<-time.After(100 * time.Millisecond)
	s.HandleKeyEvent(tcell.NewEventKey(tcell.KeyEnter, -1, tcell.ModNone))
	<-time.After(300 * time.Millisecond)

	// then
	assert.True(t, app.GetApp().IsQuit())
}

func Test_fjiraSwitchWorkspaceView_should_handle_fuzzy_find_result(t *testing.T) {
	// given
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()
	app.CreateNewAppWithScreen(screen)
	s := newSwitchWorkspaceView()

	// when
	s.Init()
	<-time.After(100 * time.Millisecond)
	s.fuzzyFind.Complete <- app.FuzzyFindResult{Index: 0, Match: "test"}
	<-time.After(300 * time.Millisecond)

	// then
	assert.True(t, app.GetApp().IsQuit())
}
