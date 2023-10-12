package filters

import (
	"bytes"
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	os2 "github.com/mk-5/fjira/internal/os"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
	"time"
)

func prepareTestScreen(t *testing.T) tcell.SimulationScreen {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	app.InitTestApp(screen)
	tempDir := t.TempDir()
	_ = os2.SetUserHomeDir(tempDir)
	return screen
}

func TestNewFiltersSearchView(t *testing.T) {
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
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				body := `[
{
    "self": "https://test/rest/api/2/filter/10006",
    "id": "10006",
    "name": "Filter for FJIR board",
    "owner": {
        "self": "https://test",
        "accountId": "test",
        "avatarUrls": {
            "48x48": "https://test/48",
            "24x24": "https://test/24",
            "16x16": "https://test/16",
            "32x32": "https://test/32"
        },
        "displayName": "Test",
        "active": true
    },
    "jql": "project = FJIR ORDER BY Rank ASC",
    "viewUrl": "https://test/issues/?filter=10006",
    "searchUrl": "https://test/search?jql=project+%3D+FJIR+ORDER+BY+Rank+ASC",
    "favourite": false,
    "favouritedCount": 0,
    "sharePermissions": [
        {
            "id": 10007,
            "type": "project",
            "project": {
                "self": "https://test/rest/api/2/project/10006",
                "id": "10006",
                "key": "FJIR",
                "assigneeType": "PROJECT_LEAD",
                "name": "FJIRA",
                "roles": {},
                "avatarUrls": {
                    "48x48": "https://test",
                    "24x24": "https://test",
                    "16x16": "https://test",
                    "32x32": "https://test"
                },
                "projectTypeKey": "software",
                "simplified": true,
                "style": "next-gen",
                "properties": {},
                "entityId": "",
                "uuid": ""
            }
        }
    ],
    "editPermissions": [],
    "isWritable": true,
    "sharedUsers": {
        "size": 33,
        "items": [],
        "max-results": 1000,
        "start-index": 0,
        "end-index": 0
    },
    "subscriptions": {
        "size": 0,
        "items": [],
        "max-results": 1000,
        "start-index": 0,
        "end-index": 0
    }
}]
`
				w.Write([]byte(body)) //nolint:errcheck
			})
			view := NewFiltersView(api).(*filtersSearchView)

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
			assert.Contains(t, result, "Filter for FJIR board")

			// and then
			view.Destroy()
		})
	}
}

func Test_fjiraFiltersView_handleCancelActionWithQuit(t *testing.T) {
	type args struct {
		appQuit bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"should handle CANCEL and quit", args{appQuit: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			screen := prepareTestScreen(t)
			defer screen.Fini()

			view := NewFiltersView(jira.NewJiraApiMock(nil)).(*filtersSearchView)

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

func Test_fjiraFiltersView_startFiltersFuzzyFind(t *testing.T) {
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
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				body := `[
{
    "self": "https://test/rest/api/2/filter/10006",
    "id": "10006",
    "name": "Filter for FJIR board",
    "owner": {
        "self": "https://test",
        "accountId": "test",
        "avatarUrls": {
            "48x48": "https://test/48",
            "24x24": "https://test/24",
            "16x16": "https://test/16",
            "32x32": "https://test/32"
        },
        "displayName": "Test",
        "active": true
    },
    "jql": "project = FJIR ORDER BY Rank ASC",
    "viewUrl": "https://test/issues/?filter=10006",
    "searchUrl": "https://test/search?jql=project+%3D+FJIR+ORDER+BY+Rank+ASC",
    "favourite": false,
    "favouritedCount": 0,
    "sharePermissions": [
        {
            "id": 10007,
            "type": "project",
            "project": {
                "self": "https://test/rest/api/2/project/10006",
                "id": "10006",
                "key": "FJIR",
                "assigneeType": "PROJECT_LEAD",
                "name": "FJIRA",
                "roles": {},
                "avatarUrls": {
                    "48x48": "https://test",
                    "24x24": "https://test",
                    "16x16": "https://test",
                    "32x32": "https://test"
                },
                "projectTypeKey": "software",
                "simplified": true,
                "style": "next-gen",
                "properties": {},
                "entityId": "",
                "uuid": ""
            }
        }
    ],
    "editPermissions": [],
    "isWritable": true,
    "sharedUsers": {
        "size": 33,
        "items": [],
        "max-results": 1000,
        "start-index": 0,
        "end-index": 0
    },
    "subscriptions": {
        "size": 0,
        "items": [],
        "max-results": 1000,
        "start-index": 0,
        "end-index": 0
    }
}]
`
				w.Write([]byte(body)) //nolint:errcheck
			})
			view := NewFiltersView(api).(*filtersSearchView)

			// when
			go view.startFiltersFuzzyFind()
			for view.fuzzyFind == nil {
				<-time.After(10 * time.Millisecond)
			}
			<-time.After(100 * time.Millisecond)
			view.HandleKeyEvent(tcell.NewEventKey(tcell.KeyEnter, 'e', tcell.ModNone))

			// then
			assert.NotNil(t, view.fuzzyFind)
		})
	}
}
