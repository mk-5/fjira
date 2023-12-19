package users

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mk-5/fjira/internal/app"
	"github.com/mk-5/fjira/internal/jira"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func TestNewFuzzyFind(t *testing.T) {
	screen := tcell.NewSimulationScreen("utf-8")
	_ = screen.Init() //nolint:errcheck
	defer screen.Fini()
	app.InitTestApp(screen)

	tests := []struct {
		name string
	}{
		{"should use api find up to typeaheadThreshold, then fuzzy-find"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			apiCall := false
			sb2User := strings.Builder{}
			sb2User.WriteString(`[{"id": "U1", "displayName": "Bob"}, {"id": "U2", "displayName": "John"}]`)
			sb1000Users := strings.Builder{}
			sb1000Users.WriteString("[")
			for i := 0; i < 1000; i++ {
				sb1000Users.WriteString(`{"id": "U1", "displayName": "Bob"}`)
				if i != 999 {
					sb1000Users.WriteString(",")
				}
			}
			sb1000Users.WriteString("]")
			var apiResult *strings.Builder
			api := jira.NewJiraApiMock(func(w http.ResponseWriter, r *http.Request) {
				apiCall = true
				w.WriteHeader(200)
				_, _ = w.Write([]byte(apiResult.String()))
			})
			fuzzyFind, us := NewFuzzyFind("ABC", api)
			fuzzyFind.SetDebounceDisabled(true)

			// when
			apiCall = false
			apiResult = &sb1000Users
			fuzzyFind.SetQuery("")
			fuzzyFind.Update()

			// then
			assert.True(t, apiCall)
			assert.Equal(t, 1001, len(*us))

			// when
			apiCall = false
			apiResult = &sb1000Users
			fuzzyFind.SetQuery("b")
			fuzzyFind.Update()

			// then
			assert.True(t, apiCall)
			assert.Equal(t, 1001, len(*us))

			// when
			apiCall = false
			apiResult = &sb2User
			fuzzyFind.SetQuery("bo")
			fuzzyFind.Update()

			// then
			assert.True(t, apiCall)
			assert.Equal(t, 3, len(*us))

			// when
			apiCall = false
			apiResult = &sb2User
			fuzzyFind.SetQuery("bo")
			fuzzyFind.Update()

			// then
			assert.False(t, apiCall, "api shouldn't be called because previous call returned 2 records")
			assert.Equal(t, 3, len(*us))
		})
	}
}
