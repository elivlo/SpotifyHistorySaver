package models

// ArtistsTrack is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type ArtistsTrack struct {
	ID       int    `json:"id" db:"id"`
	ArtistID string `json:"artist_id" db:"artist_id"`
	TrackID  string `json:"track_id" db:"track_id"`
}

// ArtistsTracks is not required by pop and may be deleted
type ArtistsTracks []ArtistsTrack
