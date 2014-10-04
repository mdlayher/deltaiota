package models

import (
	"errors"
)

var (
	// ErrInvalid is returned when a model fails a call to Validate.
	ErrInvalid = errors.New("models: invalid data")
)

// Validator provides the Validate method, which ensures that fields on a struct
// contain valid values.  An error is returned if any values are not valid.
type Validator interface {
	Validate() error
}
