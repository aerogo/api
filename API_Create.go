package api

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/aerogo/aero"
)

// Create ...
func (api *API) Create(collection string) (string, aero.Handle) {
	objType := api.Type(collection)
	objTypeName := objType.Name()
	creatableInDBInterface := reflect.TypeOf((*Newable)(nil)).Elem()

	if !reflect.PtrTo(objType).Implements(creatableInDBInterface) {
		return "", nil
	}

	route := api.root + "new/" + strings.ToLower(objTypeName)
	handler := func(ctx *aero.Context) (body string) {
		// Recover from panics
		defer func() {
			if r := recover(); r != nil {
				err := r.(error)
				fmt.Println("Recovered panic in", ctx.URI())
				body = ctx.Error(http.StatusInternalServerError, "Server error", err)
			}
		}()

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
		creatable.Save()

		return ctx.JSON(obj)
	}

	return route, handler
}
