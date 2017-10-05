package api

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/aerogo/aero"
)

// Action ...
func (api *API) Action(table string, action *Action) (string, aero.Handle) {
	objType := api.db.Type(table)
	objTypeName := objType.Name()
	actionableInterface := reflect.TypeOf((*Actionable)(nil)).Elem()

	if !reflect.PtrTo(objType).Implements(actionableInterface) {
		return "", nil
	}

	route := api.root + strings.ToLower(objTypeName) + "/:id" + action.Route
	handler := func(ctx *aero.Context) string {
		objID := ctx.Get("id")
		obj, err := api.db.Get(objTypeName, objID)

		if err != nil {
			return ctx.Error(http.StatusNotFound, "Not found", err)
		}

		// Authorize
		actionable := obj.(Actionable)
		err = actionable.Authorize(ctx)

		if err != nil {
			return ctx.Error(http.StatusForbidden, "Not authorized", err)
		}

		// Action
		err = action.Run(obj, ctx)

		if err != nil {
			return ctx.Error(http.StatusBadRequest, objTypeName+" could not be updated", err)
		}

		return "ok"
	}

	return route, handler
}
