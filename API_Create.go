package api

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/aerogo/aero"
)

// Create ...
func (api *API) Create(table string) (string, aero.Handle) {
	objType := api.db.Type(table)
	objTypeName := objType.Name()
	creatableInDBInterface := reflect.TypeOf((*Newable)(nil)).Elem()

	if !reflect.PtrTo(objType).Implements(creatableInDBInterface) {
		return "", nil
	}

	route := api.root + "new/" + strings.ToLower(objTypeName)
	handler := func(ctx *aero.Context) string {
		obj := reflect.New(objType).Interface()
		creatable := obj.(Newable)

		// Authorize
		err := creatable.Authorize(ctx, "create")

		if err != nil {
			return ctx.Error(http.StatusForbidden, "Not authorized", err)
		}

		// Create
		err = creatable.Create(ctx)

		if err != nil {
			return ctx.Error(http.StatusBadRequest, objTypeName+" could not be created", err)
		}

		// Save
		err = creatable.Save()

		if err != nil {
			return ctx.Error(http.StatusInternalServerError, objTypeName+" could not be saved", err)
		}

		return ctx.JSON(obj)
	}

	return route, handler
}
