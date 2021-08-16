package models

// Artist is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type Artist struct {
	ID       string    `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
}

// Artists is not required by pop and may be deleted
type Artists []Artist
