package api

import "github.com/aerogo/aero"

// An Editable can authorize changes, be changed and be saved in the database.
type Editable interface {
	Authorizable
	Savable
	Update(ctx *aero.Context, data interface{}) error
}
