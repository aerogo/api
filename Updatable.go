package api

import "github.com/aerogo/aero"

// An Updatable can authorize changes, be changed and be saved in the database.
type Updatable interface {
	Authorizable
	Savable
	Update(ctx *aero.Context, data map[string]interface{}) error
}
