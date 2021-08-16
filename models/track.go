package models

// Track is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type Track struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	TrackNumber int    `json:"track_number" db:"track_number"`
	DiscNumber  int    `json:"disc_number" db:"disc_number"`
	Explicit    bool   `json:"explicit" db:"explicit"`
}

// Tracks is not required by pop and may be deleted
type Tracks []Track
