package boards

import (
	"net/http"
	"testing"

	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	assert2 "github.com/stretchr/testify/assert"
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

func TestGoIntoBoardView_WithSprints(t *testing.T) {
	RegisterGoTo()
	app.InitTestApp(nil)

	api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/rest/agile/1.0/board/1/configuration":
			// Scrum board configuration with filter id
			w.WriteHeader(200)
			_, _ = w.Write([]byte(`{
  "id": 1,
  "type": "scrum",
  "filter": { "id": "10000" },
  "columnConfig": { "columns": [] }
}`))
		case "/rest/api/2/filter/10000":
			// Filter with some JQL
			w.WriteHeader(200)
			_, _ = w.Write([]byte(`{
  "id": "10000",
  "name": "Test Filter",
  "jql": "project = GEN",
  "favourite": true
}`))
		case "/rest/agile/1.0/board/1/sprint":
			// Ensure correct query param is passed for sprints
			if got := r.URL.Query().Get("state"); got != "active,future" {
				t.Fatalf("expected state query param 'active,future', got '%s'", got)
			}
			w.WriteHeader(200)
			_, _ = w.Write([]byte(`{
  "maxResults": 50,
  "startAt": 0,
  "total": 2,
  "isLast": true,
  "values": [
    {
      "id": 11,
      "self": "https://test.net/rest/agile/1.0/sprint/11",
      "state": "future",
      "name": "Sprint Future",
      "startDate": null,
      "endDate": null,
      "completeDate": null,
      "originBoardId": 1,
      "goal": ""
    },
    {
      "id": 10,
      "self": "https://test.net/rest/agile/1.0/sprint/10",
      "state": "active",
      "name": "Sprint Active",
      "startDate": null,
      "endDate": null,
      "completeDate": null,
      "originBoardId": 1,
      "goal": ""
    }
  ]
}`))
		default:
			w.WriteHeader(200)
			_, _ = w.Write([]byte(`{}`))
		}
	})

	// when
	app.GoTo("boards", &jira.Project{}, &jira.BoardItem{Id: 1}, func() {}, api)

	// then
	view, ok := app.GetApp().CurrentView().(*boardView)
	if !ok {
		t.Fatalf("expected boardView to be current view")
	}
	assert := assert2.New(t)
	assert.NotNil(view.activeSprint, "active sprint should be set")
	assert.Equal("active", view.activeSprint.State, "active sprint should be selected")
	assert.Equal("Sprint Active", view.activeSprint.Name, "active sprint name mismatch")
	assert.True(len(view.sprints) >= 2, "sprints should be set on the view")
}
