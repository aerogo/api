package api

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/aerogo/aero"
	"github.com/aerogo/mirror"
)

// EditField ...
func (api *API) EditField(collection string) (string, aero.Handler) {
	objType := api.Type(collection)
	objTypeName := objType.Name()
	editableInterface := reflect.TypeOf((*Editable)(nil)).Elem()

	if !reflect.PtrTo(objType).Implements(editableInterface) {
		return "", nil
	}

	route := api.root + strings.ToLower(objTypeName) + "/:id/field/:field"
	handler := func(ctx aero.Context) error {
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
		edits, err := ctx.Request().Body().JSONObject()

		if err != nil {
			return ctx.Error(http.StatusBadRequest, "Invalid data format (expected JSON)", err)
		}

		// Get the field that we're editing
		_, _, fieldValue, err := mirror.GetField(obj, field)

		if err != nil {
			return ctx.Error(http.StatusBadRequest, "Could not find"+field+" in type "+objTypeName, err)
		}

		// Set properties
		fieldToEdit := fieldValue.Addr().Interface()
		err = SetObjectProperties(fieldToEdit, edits, ctx)

		if err != nil {
			return ctx.Error(http.StatusInternalServerError, fieldValue.Type().Name()+" could not be edited", err)
		}

		// Save
		editable.Save()

		return ctx.String("ok")
	}

	return route, handler
}
