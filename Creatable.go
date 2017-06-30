package api

import "github.com/aerogo/aero"

// Creatable defines an object type where new instances can be created by users.
type Creatable interface {
	Savable
	Authorizable
	Create(ctx *aero.Context) error
}
