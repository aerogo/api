package api

import (
	"net/http"
	"strings"

	"github.com/aerogo/aero"
)

// Get ...
func (api *API) Get(table string) (string, aero.Handle) {
	objType := api.Type(table)
	objTypeName := objType.Name()
	// filterInterface := reflect.TypeOf((*Filter)(nil)).Elem()
	// filterEnabled := reflect.PtrTo(objType).Implements(filterInterface)

	route := api.root + strings.ToLower(objTypeName) + "/:id"
	handler := func(ctx *aero.Context) string {
		objID := ctx.Get("id")

		// Fetch object
		obj, err := api.db.Get(objTypeName, objID)

		if err != nil {
			return ctx.Error(http.StatusNotFound, "Not found", err)
		}

		// // Remove private data
		// if filterEnabled {
		// 	filter := obj.(Filter)

		// 	if filter.ShouldFilter(ctx) {
		// 		filter.Filter()
		// 	}
		// }

		// Allow CORS
		ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")

		return ctx.JSON(obj)
	}

	return route, handler
}
