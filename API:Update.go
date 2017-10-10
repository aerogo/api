package api

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/aerogo/aero"
)

// Update ...
func (api *API) Update(table string) (string, aero.Handle) {
	objType := api.db.Type(table)
	objTypeName := objType.Name()
	editableInterface := reflect.TypeOf((*Updatable)(nil)).Elem()

	if !reflect.PtrTo(objType).Implements(editableInterface) {
		return "", nil
	}

	route := api.root + strings.ToLower(objTypeName) + "/:id"
	handler := func(ctx *aero.Context) string {
		objID := ctx.Get("id")
		obj, err := api.db.Get(objTypeName, objID)

		if err != nil {
			return ctx.Error(http.StatusNotFound, "Not found", err)
		}

		// Authorize
		editable := obj.(Updatable)
		err = editable.Authorize(ctx)

		if err != nil {
			return ctx.Error(http.StatusForbidden, "Not authorized", err)
		}

		// Parse body
		body := ctx.RequestBody()

		var data interface{}
		err = json.Unmarshal(body, &data)

		if err != nil {
			return ctx.Error(http.StatusBadRequest, "Invalid data format (expected JSON)", err)
		}

		// Edit
		err = editable.Update(ctx, data)

		if err != nil {
			return ctx.Error(http.StatusBadRequest, objTypeName+" could not be updated", err)
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