package api

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/aerogo/aero"
	"github.com/aerogo/mirror"
)

// EditProperty ...
func (api *API) EditProperty(table string) (string, aero.Handle) {
	objType := api.db.Type(table)
	objTypeName := objType.Name()
	editableInterface := reflect.TypeOf((*Editable)(nil)).Elem()

	if !reflect.PtrTo(objType).Implements(editableInterface) {
		return "", nil
	}

	route := api.root + strings.ToLower(objTypeName) + "/:id/:property"
	handler := func(ctx *aero.Context) string {
		objID := ctx.Get("id")
		property := ctx.Get("property")

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

		// Get the object that we're editing
		_, _, value, err := mirror.GetField(obj, property)

		if err != nil {
			return ctx.Error(http.StatusBadRequest, "This property does not exist in type "+objTypeName, err)
		}

		// Set properties
		fieldToEdit := value.Interface().(Editable)
		err = SetObjectProperties(fieldToEdit, edits, ctx)

		if err != nil {
			return ctx.Error(http.StatusInternalServerError, value.Type().Name()+" could not be edited", err)
		}

		// Save
		err = editable.Save()

		if err != nil {
			return ctx.Error(http.StatusInternalServerError, objTypeName+" could not be saved", err)
		}

		return "ok"
	}

	return route, handler
}
