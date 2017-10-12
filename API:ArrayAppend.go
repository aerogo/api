package api

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/aerogo/aero"
	"github.com/aerogo/mirror"
)

// ArrayAppend ...
func (api *API) ArrayAppend(table string) (string, aero.Handle) {
	objType := api.db.Type(table)
	objTypeName := objType.Name()
	editableInterface := reflect.TypeOf((*Editable)(nil)).Elem()

	if !reflect.PtrTo(objType).Implements(editableInterface) {
		return "", nil
	}

	route := api.root + strings.ToLower(objTypeName) + "/:id/field/:field/append"
	handler := func(ctx *aero.Context) string {
		objID := ctx.Get("id")
		field := ctx.Get("field")

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

		// Parse body
		edits, err := ctx.RequestBodyJSONObject()

		if err != nil {
			return ctx.Error(http.StatusBadRequest, "Invalid data format (expected JSON)", err)
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

		// Determine the type of elements the slice is holding
		sliceType := arrayType.Elem()

		if sliceType.Kind() == reflect.Ptr {
			sliceType = sliceType.Elem()
		}

		// Create new item
		newItem := reflect.New(sliceType)
		SetObjectProperties(newItem.Interface(), edits, ctx)

		// Append item
		var newSlice reflect.Value

		if arrayType.Elem().Kind() == reflect.Ptr {
			newSlice = reflect.Append(arrayValue, newItem)
		} else {
			newSlice = reflect.Append(arrayValue, newItem.Elem())
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
