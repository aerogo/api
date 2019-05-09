package api

import (
	"reflect"

	"github.com/aerogo/aero"
)

// API represents your application's API configuration.
type API struct {
	// The base endpoint for all API requests.
	root string

	// The database used.
	db Database

	// An internal list of registered actions.
	actions []*Action
}

// New creates a new API.
func New(root string, db Database) *API {
	return &API{
		root: root,
		db:   db,
	}
}

// Install installs the REST & GraphQL API to your webserver at the given endpoint.
func (api *API) Install(app *aero.Application) {
	for collection, typ := range api.db.Types() {
		api.RegisterCollection(app, collection, typ)
	}

	for _, action := range api.actions {
		route, handler := api.ActionHandler(action)
		app.Post(route, handler)
	}
}

// RegisterCollection registers a single collection.
func (api *API) RegisterCollection(app *aero.Application, collection string, objType reflect.Type) {
	// Get
	route, handler := api.Get(collection)
	app.Get(route, handler)

	// Get property
	route, handler = api.GetField(collection)
	app.Get(route, handler)

	// Edit
	route, handler = api.Edit(collection)

	if route != "" && handler != nil {
		app.Post(route, handler)
	}

	// Edit property
	route, handler = api.EditField(collection)

	if route != "" && handler != nil {
		app.Post(route, handler)
	}

	// Delete
	route, handler = api.Delete(collection)

	if route != "" && handler != nil {
		app.Post(route, handler)
	}

	// Append array element
	route, handler = api.ArrayAppend(collection)

	if route != "" && handler != nil {
		app.Post(route, handler)
	}

	// Remove array element
	route, handler = api.ArrayRemove(collection)

	if route != "" && handler != nil {
		app.Post(route, handler)
	}

	// Create
	route, handler = api.Create(collection)

	if route != "" && handler != nil {
		app.Post(route, handler)
	}
}

// RegisterAction registers an action for a collection.
func (api *API) RegisterAction(action *Action) {
	api.actions = append(api.actions, action)
}

// RegisterActions registers actions for a collection.
func (api *API) RegisterActions(collection string, actions []*Action) {
	for _, action := range actions {
		action.Collection = collection
	}

	api.actions = append(api.actions, actions...)
}

// Type ...
func (api *API) Type(collection string) reflect.Type {
	return api.db.Types()[collection]
}
