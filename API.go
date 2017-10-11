package api

import (
	"reflect"

	"github.com/aerogo/aero"
)

var collectionInterface = reflect.TypeOf((*Collection)(nil)).Elem()

// API ...
type API struct {
	root    string
	db      Database
	actions []*Action
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

	for _, action := range api.actions {
		route, handler := api.ActionHandler(action)
		app.Post(route, handler)
	}
}

// RegisterTable registers a single table.
func (api *API) RegisterTable(app *aero.Application, table string, objType reflect.Type) {
	// Get
	route, handler := api.Get(table)
	app.Get(route, handler)

	// Get property
	route, handler = api.GetField(table)
	app.Get(route, handler)

	// Edit
	route, handler = api.Edit(table)

	if route != "" && handler != nil {
		app.Post(route, handler)
	}

	// Edit property
	route, handler = api.EditField(table)

	if route != "" && handler != nil {
		app.Post(route, handler)
	}

	// Create
	route, handler = api.Create(table)

	if route != "" && handler != nil {
		app.Post(route, handler)
	}

	// // Collections
	// if reflect.PtrTo(objType).Implements(collectionInterface) {
	// 	api.RegisterCollection(app, table)
	// }
}

// RegisterAction registers an action for a table.
func (api *API) RegisterAction(action *Action) {
	api.actions = append(api.actions, action)
}

// RegisterActions registers actions for a table.
func (api *API) RegisterActions(actions []*Action) {
	api.actions = append(api.actions, actions...)
}
