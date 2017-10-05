package api

import "github.com/aerogo/aero"

// Action defines a single action on a given datatype.
type Action struct {
	Table string
	Route string
	Run   func(obj interface{}, ctx *aero.Context) error
}
