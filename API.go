package api

import (
	"reflect"

	"github.com/aerogo/aero"
)

var collectionInterface = reflect.TypeOf((*Collection)(nil)).Elem()

// API ...
type API struct {
	root                string
	db                  Database
	actionRegistrations []*ActionRegistration
}

// ActionRegistration ...
type ActionRegistration struct {
	Table  string
	Action *Action
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

	for _, reg := range api.actionRegistrations {
		route, handler := api.Action(reg.Table, reg.Action)
		app.Post(route, handler)
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

	// Collections
	if reflect.PtrTo(objType).Implements(collectionInterface) {
		api.RegisterCollection(app, table)
	}
}

// RegisterAction registers an action for a table.
func (api *API) RegisterAction(table string, action *Action) {
	api.actionRegistrations = append(api.actionRegistrations, &ActionRegistration{
		Table:  table,
		Action: action,
	})
}
