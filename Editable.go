package api

import "github.com/aerogo/aero"

// An Editable can authorize changes, be changed and be saved in the database.
type Editable interface {
	Authorizable
	Savable
	Edit(ctx *aero.Context, data map[string]interface{}) error
}
