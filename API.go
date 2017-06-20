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

// Install ...
func (api *API) Install(app *aero.Application) {
	collection := reflect.TypeOf((*Collection)(nil)).Elem()

	for table, objType := range api.db.Types() {
		// GET route
		route, handler := api.Get(table)
		app.Get(route, handler)

		// Collections
		if reflect.PtrTo(objType).Implements(collection) {
			// Add
			route, handler = api.Collection(table, "add", func(coll Collection, item interface{}) error {
				return coll.Add(item)
			})
			app.Post(route, handler)

			// Remove
			route, handler = api.Collection(table, "remove", func(coll Collection, item interface{}) error {
				coll.Remove(item)
				return nil
			})
			app.Post(route, handler)
		}
	}
}

// Collection ...
func (api *API) Collection(table string, modificationName string, modify CollectionModification) (string, aero.Handle) {
	objType := api.db.Type(table)
	objTypeName := objType.Name()

	route := api.root + strings.ToLower(objTypeName) + "/:id/" + modificationName
	handler := api.CollectionHandler(objTypeName, modify)

	return route, handler
}

// CollectionModification ...
type CollectionModification func(Collection, interface{}) error

// CollectionHandler ...
func (api *API) CollectionHandler(objTypeName string, modify CollectionModification) aero.Handle {
	return func(ctx *aero.Context) string {
		objID := ctx.Get("id")
		obj, err := api.db.Get(objTypeName, objID)

		if err != nil {
			return ctx.Error(http.StatusNotFound, "Collection not found", err)
		}

		body := ctx.RequestBody()
		coll := obj.(Collection)
		err = coll.Authorize(ctx)

		if err != nil {
			return ctx.Error(http.StatusForbidden, "Not authorized", err)
		}

		item := coll.TransformBody(body)

		err = modify(coll, item)

		if err != nil {
			return ctx.Error(http.StatusInternalServerError, "Error adding item to list", err)
		}

		err = coll.Save()

		if err != nil {
			return ctx.Error(http.StatusInternalServerError, "Error saving list in database", err)
		}

		return "ok"
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
