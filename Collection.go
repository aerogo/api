package api

import (
	"github.com/aerogo/aero"
)

// Collection interface for lists, sets, etc.
type Collection interface {
	Add(interface{}) error
	Remove(interface{}) bool
	Contains(interface{}) bool

	Get(id interface{}) (interface{}, error)
	Set(id interface{}, value interface{}) error

	// TransformBody returns an item that is passed to methods like Add, Remove, etc.
	TransformBody(body []byte) interface{}

	// Authorize returns an error if the given API request is not authorized.
	Authorize(*aero.Context) error

	// Save saves the object in the database.
	Save() error
}
