package models

import (
	"sort"
	"time"
)

// HistoryEntry is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type HistoryEntry struct {
	ID       int       `json:"id" db:"id"`
	TrackId  string    `json:"track_id" db:"track_id"`
	PlayedAt time.Time `json:"played_at" db:"played_at"`
}

// HistoryEntries is not required by pop and may be deleted
type HistoryEntries []HistoryEntry

// SortByDate will sort an array of HistoryEntries by date ascending
func (h HistoryEntries) SortByDate() {
	sort.Slice(h, func(i, j int) bool {
		return h[i].PlayedAt.Before(h[j].PlayedAt)
	})
}
