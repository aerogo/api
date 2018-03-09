package api

import (
	"github.com/aerogo/aero"
)

// Deletable defines an object type that can be deleted from the database.
type Deletable interface {
	Authorizable
	DeleteInContext(*aero.Context) error
	Delete() error
}
