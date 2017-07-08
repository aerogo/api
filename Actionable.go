package api

import "github.com/aerogo/aero"

// Actionable means the data type can execute actions.
type Actionable interface {
	Authorizable
	Savable
	Action(*aero.Context, string) error
}
