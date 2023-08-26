package boards

import (
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	assert2 "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGoIntoBoardView(t *testing.T) {
	RegisterGoTo()
	app.InitTestApp(nil)

	type args struct {
		gotoMethod    func()
		viewPredicate func() bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"should switch view into board view", args{
			gotoMethod: func() {
				app.GoTo("boards", &jira.Project{}, &jira.BoardItem{Id: 1}, func() {}, jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
					_, _ = w.Write([]byte("{}"))
				}))
			},
			viewPredicate: func() bool {
				_, ok := app.GetApp().CurrentView().(*boardView)
				return ok
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			tt.args.gotoMethod()

			// then
			ok := tt.args.viewPredicate()
			assert2.New(t).True(ok, "Current view is invalid.")
		})
	}
}
