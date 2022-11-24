package fjira

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"time"
)

type fjiraSwitchWorkspaceView struct {
	fuzzyFind  *app.FuzzyFind
	workspaces *userHomeWorkspaces
}

func NewSwitchWorkspaceView() *fjiraSwitchWorkspaceView {
	return &fjiraSwitchWorkspaceView{
		workspaces: &userHomeWorkspaces{},
	}
}

func (s *fjiraSwitchWorkspaceView) Init() {
	records, err := s.workspaces.readAllWorkspaces()
	if err != nil {
		panic(err.Error())
	}
	s.fuzzyFind = app.NewFuzzyFind(MessageSelectWorkspace, records)
	s.fuzzyFind.MarginBottom = 0
	go func() {
		if workspace := <-s.fuzzyFind.Complete; true {
			if workspace.Index < 0 {
				app.GetApp().Quit()
				return
			}
			err := s.workspaces.setCurrentWorkspace(workspace.Match)
			if err != nil {
				app.Error(err.Error())
				app.GetApp().Quit()
				return
			}
			app.Success(fmt.Sprintf(MessageSelectWorkspaceSuccess, workspace.Match))
			time.Sleep(2 * time.Second)
			app.GetApp().Quit()
		}
	}()
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
