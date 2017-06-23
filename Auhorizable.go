package api

import "github.com/aerogo/aero"

// Authorizable means the data type can authorize changes.
type Authorizable interface {
	// Authorize returns an error if the given API request is not authorized.
	Authorize(*aero.Context) error
}
