package api

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/aerogo/aero"
)

// Actions ...
func (api *API) Actions(table string) (string, aero.Handle) {
	objType := api.db.Type(table)
	objTypeName := objType.Name()
	actionableInterface := reflect.TypeOf((*Actionable)(nil)).Elem()

	if !reflect.PtrTo(objType).Implements(actionableInterface) {
		return "", nil
	}

	route := api.root + strings.ToLower(objTypeName) + "/:id/:action"
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
		err = actionable.Action(ctx, ctx.Get("action"))

		if err != nil {
			return ctx.Error(http.StatusBadRequest, objTypeName+" could not be updated", err)
		}

		// Save
		err = actionable.Save()

		if err != nil {
			return ctx.Error(http.StatusInternalServerError, objTypeName+" could not be saved", err)
		}

		return "ok"
	}

	return route, handler
}
