package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
)
// Artist is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type Artist struct {
	ID       string    `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
}

// String is not required by pop and may be deleted
func (a Artist) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Artists is not required by pop and may be deleted
type Artists []Artist

// String is not required by pop and may be deleted
func (a Artists) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *Artist) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *Artist) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *Artist) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
