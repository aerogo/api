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
func (api *API) ArrayRemove(collection string) (string, aero.Handle) {
	objType := api.Type(collection)
	objTypeName := objType.Name()
	editableInterface := reflect.TypeOf((*Editable)(nil)).Elem()

	if !reflect.PtrTo(objType).Implements(editableInterface) {
		return "", nil
	}

	route := api.root + strings.ToLower(objTypeName) + "/:id/field/:field/remove/:index"
	handler := func(ctx *aero.Context) string {
		objID := ctx.Get("id")
		field := ctx.Get("field")
		indexString := ctx.Get("index")

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

		// Find index
		value, index, err := mirror.GetSliceElement(arrayValue.Interface(), indexString)

		if err != nil {
			return ctx.Error(http.StatusBadRequest, "Could not find array element using index "+indexString, err)
		}

		// Create a new slice where the removed item does not exist anymore
		oldLen := arrayValue.Len()
		newLen := oldLen - 1

		if newLen < 0 {
			return ctx.Error(http.StatusBadRequest, "Can't remove array element", errors.New("Array is empty"))
		}

		newSlice := reflect.MakeSlice(arrayType, newLen, newLen)

		offset := 0
		for i := 0; i < newLen; i++ {
			if i == index {
				offset++
			}

			newSlice.Index(i).Set(arrayValue.Index(i + offset))
		}

		arrayValue.Set(newSlice)

		// Call OnRemove
		listener, isRemoveListener := obj.(ArrayEventListener)

		if isRemoveListener {
			listener.OnRemove(ctx, field, index, value.Interface())
		}

		// Save
		editable.Save()

		return "ok"
	}

	return route, handler
}
