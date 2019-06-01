package api

import "github.com/aerogo/aero"

// Filter describes an object with private data that needs to be filtered in the public API.
type Filter interface {
	ShouldFilter(aero.Context) bool
	Filter()
}
