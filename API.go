package api

import (
	"reflect"

	"github.com/aerogo/aero"
)

var collectionInterface = reflect.TypeOf((*Collection)(nil)).Elem()

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
	for table, objType := range api.db.Types() {
		api.RegisterTable(app, table, objType)
	}
}

// RegisterTable registers a single table.
func (api *API) RegisterTable(app *aero.Application, table string, objType reflect.Type) {
	// Get
	route, handler := api.Get(table)
	app.Get(route, handler)

	// Update
	route, handler = api.Update(table)

	if route != "" && handler != nil {
		app.Post(route, handler)
	}

	// Create
	route, handler = api.Create(table)

	if route != "" && handler != nil {
		app.Post(route, handler)
	}

	// Actions
	route, handler = api.Actions(table)

	if route != "" && handler != nil {
		app.Post(route, handler)
	}

	// Collections
	if reflect.PtrTo(objType).Implements(collectionInterface) {
		api.RegisterCollection(app, table)
	}
}
