package api

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/mohae/deepcopy"

	"github.com/aerogo/aero"
	"github.com/aerogo/mirror"
)

// GetField ...
func (api *API) GetField(collection string) (string, aero.Handler) {
	objType := api.Type(collection)
	objTypeName := objType.Name()
	filterInterface := reflect.TypeOf((*Filter)(nil)).Elem()
	filterEnabled := reflect.PtrTo(objType).Implements(filterInterface)

	route := api.root + strings.ToLower(objTypeName) + "/:id/field/:field"
	handler := func(ctx aero.Context) error {
		objID := ctx.Get("id")
		field := ctx.Get("field")

		// Fetch object
		obj, err := api.db.Get(objTypeName, objID)

		if err != nil {
			return ctx.Error(http.StatusNotFound, "Not found", err)
		}

		// // Remove private data
		if filterEnabled {
			obj = deepcopy.Copy(obj)
			filter := obj.(Filter)

			if filter.ShouldFilter(ctx) {
				filter.Filter()
			}
		}

		// Allow CORS
		ctx.Response().SetHeader("Access-Control-Allow-Origin", "*")

		// Get field
		_, _, value, err := mirror.GetField(obj, field)

		if err != nil {
			return ctx.Error(http.StatusBadRequest, "Could not find"+field+" in type "+objTypeName, err)
		}

		return ctx.JSON(value.Interface())
	}

	return route, handler
}
