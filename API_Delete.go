package api

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/aerogo/aero"
)

// Delete ...
func (api *API) Delete(collection string) (string, aero.Handle) {
	objType := api.Type(collection)
	objTypeName := objType.Name()
	deletableInterface := reflect.TypeOf((*Deletable)(nil)).Elem()

	if !reflect.PtrTo(objType).Implements(deletableInterface) {
		return "", nil
	}

	route := api.root + strings.ToLower(objTypeName) + "/:id/delete"
	handler := func(ctx *aero.Context) string {
		objID := ctx.Get("id")
		obj, err := api.db.Get(objTypeName, objID)

		if err != nil {
			return ctx.Error(http.StatusNotFound, "Not found", err)
		}

		// Authorize
		deletable := obj.(Deletable)
		err = deletable.Authorize(ctx, "delete")

		if err != nil {
			return ctx.Error(http.StatusForbidden, "Not authorized", err)
		}

		// Delete
		err = deletable.DeleteInContext(ctx)

		if err != nil {
			return ctx.Error(http.StatusInternalServerError, objTypeName+" could not be deleted", err)
		}

		return "ok"
	}

	return route, handler
}
