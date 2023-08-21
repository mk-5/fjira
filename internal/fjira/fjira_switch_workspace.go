package fjira

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"time"
)

type fjiraSwitchWorkspaceView struct {
	fuzzyFind     *app.FuzzyFind
	fjiraSettings *userHomeSettingsStorage
}

func newSwitchWorkspaceView() *fjiraSwitchWorkspaceView {
	return &fjiraSwitchWorkspaceView{
		fjiraSettings: &userHomeSettingsStorage{},
	}
}

func (s *fjiraSwitchWorkspaceView) Init() {
	records, err := s.fjiraSettings.readAllWorkspaces()
	if err != nil {
		panic(err.Error())
	}
	s.fuzzyFind = app.NewFuzzyFind(MessageSelectWorkspace, records)
	s.fuzzyFind.MarginBottom = 0
	app.GetApp().SetDirty()
	go s.waitForFuzzyFindComplete()
}

func (s *fjiraSwitchWorkspaceView) Destroy() {
	// do nothing
}

func (s *fjiraSwitchWorkspaceView) Update() {
	if s.fuzzyFind != nil {
		s.fuzzyFind.Update()
	}
}

func (s *fjiraSwitchWorkspaceView) Draw(screen tcell.Screen) {
	if s.fuzzyFind != nil {
		s.fuzzyFind.Draw(screen)
	}
}

func (s *fjiraSwitchWorkspaceView) Resize(screenX, screenY int) {
	if s.fuzzyFind != nil {
		s.fuzzyFind.Resize(screenX, screenY)
	}
}

func (s *fjiraSwitchWorkspaceView) HandleKeyEvent(keyEvent *tcell.EventKey) {
	if s.fuzzyFind != nil {
		s.fuzzyFind.HandleKeyEvent(keyEvent)
	}
}

func (s *fjiraSwitchWorkspaceView) waitForFuzzyFindComplete() {
	if workspace := <-s.fuzzyFind.Complete; true {
		if workspace.Index < 0 {
			app.GetApp().Quit()
			return
		}
		err := s.fjiraSettings.setCurrentWorkspace(workspace.Match)
		if err != nil {
			app.Error(err.Error())
			app.GetApp().Quit()
			return
		}
		app.Success(fmt.Sprintf(MessageSelectWorkspaceSuccess, workspace.Match))
		time.Sleep(2 * time.Second)
		app.GetApp().Quit()
	}
}
