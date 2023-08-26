package ui

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateNavigationBars(t *testing.T) {
	tests := []struct {
		name     string
		supplier func() interface{}
	}{
		{"should create issue bottom bar", func() interface{} {
			return CreateBottomLeftBar()
		}},
		{"should create cancel bar item", func() interface{} {
			return NewCancelBarItem()
		}},
		{"should create new save bar item", func() interface{} {
			return NewSaveBarItem()
		}},
		{"should create new YES change bar item", func() interface{} {
			return NewYesBarItem()
		}},
		{"should create new OPEN bar item", func() interface{} {
			return NewOpenBarItem()
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNilf(t, tt.supplier(), "CreateCommentBarItem()")
		})
	}
}
