package api

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/mohae/deepcopy"

	"github.com/aerogo/aero"
)

// Get returns the route and the handler for the given collection.
func (api *API) Get(collection string) (string, aero.Handler) {
	objType := api.Type(collection)
	typeName := objType.Name()
	filterInterface := reflect.TypeOf((*Filter)(nil)).Elem()
	filterEnabled := reflect.PtrTo(objType).Implements(filterInterface)
	route := api.root + strings.ToLower(typeName) + "/:id"

	handler := func(ctx aero.Context) error {
		id := ctx.Get("id")

		// Fetch object
		obj, err := api.db.Get(typeName, id)

		if err != nil {
			return ctx.Error(http.StatusNotFound, "Not found", err)
		}

		// Remove private data
		if filterEnabled {
			obj = deepcopy.Copy(obj)
			filter := obj.(Filter)

			if filter.ShouldFilter(ctx) {
				filter.Filter()
			}
		}

		// Allow CORS
		ctx.Response().SetHeader("Access-Control-Allow-Origin", "*")

		// Respond with JSON
		return ctx.JSON(obj)
	}

	return route, handler
}
