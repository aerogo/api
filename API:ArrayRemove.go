package api

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/aerogo/aero"
	"github.com/aerogo/mirror"
)

// ArrayRemove ...
func (api *API) ArrayRemove(table string) (string, aero.Handle) {
	objType := api.db.Type(table)
	objTypeName := objType.Name()
	editableInterface := reflect.TypeOf((*Editable)(nil)).Elem()

	if !reflect.PtrTo(objType).Implements(editableInterface) {
		return "", nil
	}

	route := api.root + strings.ToLower(objTypeName) + "/:id/field/:field/remove/:index"
	handler := func(ctx *aero.Context) string {
		objID := ctx.Get("id")
		field := ctx.Get("field")
		index, err := ctx.GetInt("index")

		if err != nil {
			return ctx.Error(http.StatusBadRequest, "Index needs to be a number", err)
		}

		// Fetch object
		obj, err := api.db.Get(objTypeName, objID)

		if err != nil {
			return ctx.Error(http.StatusNotFound, "Not found", err)
		}

		// Authorize
		editable := obj.(Editable)
		err = editable.Authorize(ctx, "edit")

		if err != nil {
			return ctx.Error(http.StatusForbidden, "Not authorized", err)
		}

		// Get the field that we're editing
		_, arrayType, arrayValue, err := mirror.GetField(obj, field)

		if err != nil {
			return ctx.Error(http.StatusBadRequest, "Could not find"+field+" in type "+objTypeName, err)
		}

		// Is the field really a slice?
		if arrayType.Kind() != reflect.Slice {
			return ctx.Error(http.StatusBadRequest, "Invalid field", errors.New("Field "+field+" is not a slice"))
		}

		// Create a new slice where the removed item does not exist anymore
		oldLen := arrayValue.Len()
		newLen := oldLen - 1
		newSlice := reflect.MakeSlice(arrayType, newLen, newLen)

		offset := 0
		for i := 0; i < newLen; i++ {
			if i == index {
				offset++
			}

			newSlice.Index(i).Set(arrayValue.Index(i + offset))
		}

		arrayValue.Set(newSlice)

		// Save
		err = editable.Save()

		if err != nil {
			return ctx.Error(http.StatusInternalServerError, objTypeName+" could not be saved", err)
		}

		return "ok"
	}

	return route, handler
}
