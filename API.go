package api

import (
	"net/http"
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

// Install ...
func (api *API) Install(app *aero.Application) {
	for table := range api.db.Types() {
		route, handler := api.Get(table)
		app.Get(route, handler)
	}
}

// Get ...
func (api *API) Get(table string) (string, aero.Handle) {
	objType := api.db.Type(table)
	objTypeName := objType.Name()

	route := api.root + strings.ToLower(objTypeName) + "/:id"
	handler := func(ctx *aero.Context) string {
		objID := ctx.Get("id")
		obj, err := api.db.Get(objTypeName, objID)

		if err != nil {
			return ctx.Error(http.StatusNotFound, "Not found", err)
		}

		return ctx.JSON(obj)
	}

	return route, handler
}
