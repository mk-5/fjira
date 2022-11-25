package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestActionBar_AddItem(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"should create action bar item"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ab := NewActionBar(0, 1)
			ab.Update()

			// when
			ab.AddTextItem("id", "test")

			// then
			assert.Equal(t, 1, len(ab.items))

			// and when
			ab.RemoveItemAtIndex(0)

			// then
			assert.Equal(t, 0, len(ab.items))

			// when
			ab.AddTextItem("id1", "test")
			ab.AddTextItem("id2", "test")
			ab.AddTextItem("id3", "test")
			ab.AddTextItem("id3", "test")

			// and when
			ab.TrimItemsTo(2)

			// then
			assert.Equal(t, 2, len(ab.items))

			// and when
			ab.RemoveItem(ab.items[0].Id)

			// then
			assert.Equal(t, 1, len(ab.items))
		})
	}
}
