package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"time"
)
// HistoryEntry is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type HistoryEntry struct {
	ID       int       `json:"id" db:"id"`
	TrackId  string    `json:"track_id" db:"track_id"`
	PlayedAt time.Time `json:"played_at" db:"played_at"`
}

// String is not required by pop and may be deleted
func (h HistoryEntry) String() string {
	jh, _ := json.Marshal(h)
	return string(jh)
}

// HistoryEntries is not required by pop and may be deleted
type HistoryEntries []HistoryEntry

// String is not required by pop and may be deleted
func (h HistoryEntries) String() string {
	jh, _ := json.Marshal(h)
	return string(jh)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (h *HistoryEntry) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (h *HistoryEntry) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (h *HistoryEntry) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
