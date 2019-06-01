package api

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/aerogo/aero"
)

// ActionHandler ...
func (api *API) ActionHandler(action *Action) (string, aero.Handler) {
	objType := api.Type(action.Collection)
	objTypeName := objType.Name()
	actionableInterface := reflect.TypeOf((*Actionable)(nil)).Elem()

	if !reflect.PtrTo(objType).Implements(actionableInterface) {
		return "", nil
	}

	route := api.root + strings.ToLower(objTypeName) + "/:id" + action.Route
	handler := func(ctx aero.Context) error {
		objID := ctx.Get("id")
		obj, err := api.db.Get(objTypeName, objID)

		if err != nil {
			return ctx.Error(http.StatusNotFound, "Not found", err)
		}

		// Authorize
		actionable := obj.(Actionable)
		err = actionable.Authorize(ctx, action.Name)

		if err != nil {
			return ctx.Error(http.StatusForbidden, "Not authorized", err)
		}

		// Action
		err = action.Run(obj, ctx)

		if err != nil {
			return ctx.Error(http.StatusBadRequest, objTypeName+" could not be updated", err)
		}

		return ctx.String("ok")
	}

	return route, handler
}
