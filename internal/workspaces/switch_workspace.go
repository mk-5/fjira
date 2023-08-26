package workspaces

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/ui"
	"time"
)

type switchWorkspaceView struct {
	fuzzyFind     *app.FuzzyFind
	fjiraSettings SettingsStorage
}

func NewSwitchWorkspaceView() app.View {
	return &switchWorkspaceView{
		fjiraSettings: NewUserHomeSettingsStorage(),
	}
}

func (s *switchWorkspaceView) Init() {
	records, err := s.fjiraSettings.ReadAllWorkspaces()
	if err != nil {
		panic(err.Error())
	}
	s.fuzzyFind = app.NewFuzzyFind(ui.MessageSelectWorkspace, records)
	s.fuzzyFind.MarginBottom = 0
	app.GetApp().SetDirty()
	go s.waitForFuzzyFindComplete()
}

func (s *switchWorkspaceView) Destroy() {
	// do nothing
}

func (s *switchWorkspaceView) Update() {
	if s.fuzzyFind != nil {
		s.fuzzyFind.Update()
	}
}

func (s *switchWorkspaceView) Draw(screen tcell.Screen) {
	if s.fuzzyFind != nil {
		s.fuzzyFind.Draw(screen)
	}
}

func (s *switchWorkspaceView) Resize(screenX, screenY int) {
	if s.fuzzyFind != nil {
		s.fuzzyFind.Resize(screenX, screenY)
	}
}

func (s *switchWorkspaceView) HandleKeyEvent(keyEvent *tcell.EventKey) {
	if s.fuzzyFind != nil {
		s.fuzzyFind.HandleKeyEvent(keyEvent)
	}
}

func (s *switchWorkspaceView) waitForFuzzyFindComplete() {
	if workspace := <-s.fuzzyFind.Complete; true {
		if workspace.Index < 0 {
			app.GetApp().Quit()
			return
		}
		err := s.fjiraSettings.SetCurrentWorkspace(workspace.Match)
		if err != nil {
			app.Error(err.Error())
			app.GetApp().Quit()
			return
		}
		app.Success(fmt.Sprintf(ui.MessageSelectWorkspaceSuccess, workspace.Match))
		time.Sleep(2 * time.Second)
		app.GetApp().Quit()
	}
}
