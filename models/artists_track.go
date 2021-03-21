package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
)
// ArtistsTrack is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type ArtistsTrack struct {
	ID       int    `json:"id" db:"id"`
	ArtistID string `json:"artist_id" db:"artist_id"`
	TrackID  string `json:"track_id" db:"track_id"`
}

// String is not required by pop and may be deleted
func (a ArtistsTrack) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// ArtistsTracks is not required by pop and may be deleted
type ArtistsTracks []ArtistsTrack

// String is not required by pop and may be deleted
func (a ArtistsTracks) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *ArtistsTrack) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *ArtistsTrack) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *ArtistsTrack) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
