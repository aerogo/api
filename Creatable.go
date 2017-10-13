package api

import "github.com/aerogo/aero"

// Creatable defines an object that can be created with some initial data.
type Creatable interface {
	Create(*aero.Context) error
}
