package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_HistoryEntry(t *testing.T) {
	entries := HistoryEntries{
		{ID: 0, PlayedAt: time.Now().Add(time.Minute)},
		{ID: 1, PlayedAt: time.Now()},
	}

	entries.SortByDate()

	assert.Equal(t, 1, entries[0].ID)
	assert.Equal(t, 0, entries[1].ID)
}
