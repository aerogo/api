package api

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/aerogo/aero"
	"github.com/aerogo/mirror"
)

// GetProperty ...
func (api *API) GetProperty(table string) (string, aero.Handle) {
	objType := api.db.Type(table)
	objTypeName := objType.Name()
	filterInterface := reflect.TypeOf((*Filter)(nil)).Elem()
	filterEnabled := reflect.PtrTo(objType).Implements(filterInterface)

	route := api.root + strings.ToLower(objTypeName) + "/:id/:property"
	handler := func(ctx *aero.Context) string {
		objID := ctx.Get("id")
		property := ctx.Get("property")

		// Fetch object
		obj, err := api.db.Get(objTypeName, objID)

		if err != nil {
			return ctx.Error(http.StatusNotFound, "Not found", err)
		}

		// Remove private data
		if filterEnabled {
			filter := obj.(Filter)

			if filter.ShouldFilter(ctx) {
				filter.Filter()
			}
		}

		// Allow CORS
		ctx.SetResponseHeader("Access-Control-Allow-Origin", "*")

		// Get property
		_, _, value, err := mirror.GetField(obj, property)

		if err != nil {
			return ctx.Error(http.StatusBadRequest, "This property does not exist in type "+objTypeName, err)
		}

		return ctx.JSON(value.Interface())
	}

	return route, handler
}
