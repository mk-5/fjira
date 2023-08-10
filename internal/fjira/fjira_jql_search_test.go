package fjira

import (
	"bytes"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	os2 "github.com/mk-5/fjira/internal/os"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func prepareTestScreen(t *testing.T) tcell.SimulationScreen {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	app.CreateNewAppWithScreen(screen)
	CreateNewFjira(&fjiraWorkspaceSettings{Workspace: "default"})
	tempDir := t.TempDir()
	_ = os2.SetUserHomeDir(tempDir)
	return screen
}

func TestNewJqlSearchView(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should render the view"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			screen := prepareTestScreen(t)
			defer screen.Fini()
			view := NewJqlSearchView()
			err := view.jqlStorage.addNew("test jql query")

			// when
			view.Init()
			<-time.After(time.Millisecond * 200)
			view.HandleKeyEvent(tcell.NewEventKey(tcell.KeyACK, 't', tcell.ModNone))
			// keep going app for a while
			i := 0
			for {
				view.Update()
				view.Draw(screen)
				i++
				if i > 1000000 {
					break
				}
			}
			var buffer bytes.Buffer
			contents, x, y := screen.GetContents()
			screen.Show()
			for i := 0; i < x*y; i++ {
				if len(contents[i].Bytes) != 0 {
					buffer.Write(contents[i].Bytes)
				}
			}
			result := strings.TrimSpace(buffer.String())

			// then
			assert.Nil(t, err)
			assert.Contains(t, result, "test jql query", "should contain jql from storage")

			// and then
			view.Destroy()
		})
	}
}

func Test_fjiraJqlSearchView_confirmJqlDelete(t *testing.T) {
	type args struct {
		jql string
	}
	tests := []struct {
		name string
		args args
	}{
		{"should confirm jql delete", args{jql: "test jql"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			screen := prepareTestScreen(t)
			defer screen.Fini()
			view := NewJqlSearchView()
			_ = view.jqlStorage.addNew(tt.args.jql)
			jqls, _ := view.jqlStorage.readAll()
			assert.Contains(t, jqls, tt.args.jql) // ensure that jql is added

			// when
			done := make(chan struct{})
			started := make(chan struct{})
			go func() {
				started <- struct{}{}
				view.confirmJqlDelete(tt.args.jql)
				done <- struct{}{}
			}()
			<-started
			<-time.After(200 * time.Millisecond)
			if confirmation, ok := (app.GetApp().LastDrawable()).(app.KeyListener); ok {
				confirmation.HandleKeyEvent(tcell.NewEventKey(tcell.KeyACK, 'y', tcell.ModNone))
			}
			view.Update()
			<-done

			// then
			jqls2, _ := view.jqlStorage.readAll()
			assert.NotContains(t, jqls2, tt.args.jql)
		})
	}
}

func Test_fjiraJqlSearchView_handleCancelActionWithQuit(t *testing.T) {
	type args struct {
		fuzzyFind *app.FuzzyFind
		appQuit   bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"should handle CANCEL and quit", args{fuzzyFind: nil, appQuit: false}},
		{"should handle CANCEL stay in the app", args{fuzzyFind: app.NewFuzzyFind("test", []string{}), appQuit: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			screen := prepareTestScreen(t)
			defer screen.Fini()

			view := NewJqlSearchView()
			view.fuzzyFind = tt.args.fuzzyFind

			// when
			go view.handleBottomBarActions()
			<-time.After(100 * time.Millisecond)
			go view.HandleKeyEvent(tcell.NewEventKey(tcell.KeyEscape, 'e', tcell.ModNone))

			// then
			<-time.After(100 * time.Millisecond)
			assert.Equal(t, tt.args.appQuit, app.GetApp().IsQuit())
		})
	}
}

func Test_fjiraJqlSearchView_handleNewJqlAction(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should handle new JQL action"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			screen := prepareTestScreen(t)
			defer screen.Fini()
			view := NewJqlSearchView()

			// when
			go view.handleBottomBarActions()
			<-time.After(100 * time.Millisecond)
			go view.HandleKeyEvent(tcell.NewEventKey(tcell.KeyF1, 'e', tcell.ModNone))

			// then
			<-time.After(100 * time.Millisecond)
			_, ok := app.GetApp().CurrentView().(*fjiraTextWriterView)
			assert.True(t, ok)
		})
	}
}

func Test_fjiraJqlSearchView_handleDeleteJqlAction(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should handle delete JQL action"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			screen := prepareTestScreen(t)
			defer screen.Fini()
			view := NewJqlSearchView()
			view.fuzzyFind = app.NewFuzzyFind("test", []string{})

			// when
			go view.handleBottomBarActions()
			<-time.After(100 * time.Millisecond)
			go view.HandleKeyEvent(tcell.NewEventKey(tcell.KeyF2, 'e', tcell.ModNone))

			// then
			<-time.After(100 * time.Millisecond)
			assert.Nil(t, view.fuzzyFind)
		})
	}
}

func Test_fjiraJqlSearchView_startJqlFuzzyFind(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should start jql fuzzy find"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			screen := prepareTestScreen(t)
			defer screen.Fini()
			view := NewJqlSearchView()

			// when
			go view.startJqlFuzzyFind()
			<-time.After(100 * time.Millisecond)
			view.HandleKeyEvent(tcell.NewEventKey(tcell.KeyEnter, 'e', tcell.ModNone))

			// then
			assert.NotNil(t, view.fuzzyFind)
		})
	}
}
