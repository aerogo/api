package api

import (
	"reflect"

	"github.com/aerogo/aero"
)

// An Editable can authorize changes, be changed and be saved in the database.
type Editable interface {
	Authorizable
	Savable
}

// A CustomEditable has its own implementation on how to edit certain object fields.
type CustomEditable interface {
	Edit(ctx aero.Context, key string, value reflect.Value, newValue reflect.Value) (consumed bool, err error)
}

// An AfterEditable is called after the editing process happens and before the object is saved.
type AfterEditable interface {
	AfterEdit(ctx aero.Context) error
}

// A VirtualEditable has virtual properties that do not really exist but can be set.
type VirtualEditable interface {
	VirtualEdit(ctx aero.Context, key string, newValue reflect.Value) (bool, error)
}
