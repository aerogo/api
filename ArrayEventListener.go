package api

import (
	"github.com/aerogo/aero"
)

// ArrayEventListener means the data type can authorize changes.
type ArrayEventListener interface {
	OnAppend(ctx aero.Context, field string, index int, obj interface{})
	OnRemove(ctx aero.Context, field string, index int, obj interface{})
}
