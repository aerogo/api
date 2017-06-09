package api

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/aerogo/aero"
)

// API ...
type API struct {
	root string
	db   Database
}

// New creates a new API.
func New(root string, db Database) *API {
	return &API{
		root: root,
		db:   db,
	}
}

// Get ...
func (api *API) Get(example interface{}) (string, aero.Handle) {
	listType := reflect.TypeOf(example).Elem()
	listTypeName := listType.Name()

	route := api.root + strings.ToLower(listTypeName) + "/:id"
	handler := func(ctx *aero.Context) string {
		listID := ctx.Get("id")
		list, err := api.db.Get(listTypeName, listID)

		if err != nil {
			return ctx.Error(http.StatusNotFound, "Not found", err)
		}

		return ctx.JSON(list)
	}

	return route, handler
}
