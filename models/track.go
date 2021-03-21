package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
)
// Track is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type Track struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	TrackNumber int       `json:"track_number" db:"track_number"`
	DiscNumber  int       `json:"disc_number" db:"disc_number"`
	Explicit    bool      `json:"explicit" db:"explicit"`
}

// String is not required by pop and may be deleted
func (t Track) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Tracks is not required by pop and may be deleted
type Tracks []Track

// String is not required by pop and may be deleted
func (t Tracks) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Track) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Track) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Track) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
