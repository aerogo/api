package api

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"unicode"

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
	// Retrieve
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

		// Get
		collectionType := api.db.Type(table)
		collectionTypeName := collectionType.Name()

		route = api.root + strings.ToLower(collectionTypeName) + "/:id/get/:item"
		handler = func(ctx *aero.Context) string {
			var collectionObj, item interface{}
			var err error

			collectionID := ctx.Get("id")
			collectionObj, err = api.db.Get(collectionTypeName, collectionID)

			if err != nil {
				return ctx.Error(http.StatusNotFound, "Collection not found", err)
			}

			collection := collectionObj.(Collection)

			itemID := ctx.Get("item")
			item, err = collection.Get(itemID)

			if err != nil {
				return ctx.Error(http.StatusNotFound, "Item not found", err)
			}

			return ctx.JSON(item)
		}

		app.Get(route, handler)

		// Get property
		route = api.root + strings.ToLower(collectionTypeName) + "/:id/get/:item/:property"
		handler = func(ctx *aero.Context) string {
			var collectionObj, item interface{}
			var err error

			collectionID := ctx.Get("id")
			collectionObj, err = api.db.Get(collectionTypeName, collectionID)

			if err != nil {
				return ctx.Error(http.StatusNotFound, "Collection not found", err)
			}

			collection := collectionObj.(Collection)

			itemID := ctx.Get("item")
			item, err = collection.Get(itemID)

			if err != nil || item == nil {
				return ctx.Error(http.StatusNotFound, "Item not found", err)
			}

			propertyName := ctx.Get("property")
			itemType := reflect.TypeOf(item).Elem()
			itemValue := reflect.ValueOf(item).Elem()
			property := itemValue.FieldByName(propertyName)

			if !property.IsValid() {
				message := "Property '" + propertyName + "' does not exit on type " + itemType.Name() + "<br><br>Did you mean one of these?<ul>"

				for i := 0; i < itemType.NumField(); i++ {
					field := itemType.Field(i)

					if field.Anonymous || !unicode.IsUpper([]rune(field.Name)[0]) {
						continue
					}

					message += "<li>" + field.Name + "</li>"
				}
				message += "</ul>"

				return ctx.Error(http.StatusBadRequest, message, nil)
			}

			return ctx.JSON(property.Interface())
		}

		app.Get(route, handler)

		// Update
		route = api.root + strings.ToLower(collectionTypeName) + "/:id/update/:item"
		handler = func(ctx *aero.Context) string {
			var collectionObj, item interface{}
			var err error

			collectionID := ctx.Get("id")
			collectionObj, err = api.db.Get(collectionTypeName, collectionID)

			if err != nil {
				return ctx.Error(http.StatusNotFound, "Collection not found", err)
			}

			collection := collectionObj.(Collection)

			// Authorize
			err = collection.Authorize(ctx)

			if err != nil {
				return ctx.Error(http.StatusForbidden, "Not authorized", err)
			}

			itemID := ctx.Get("item")
			body := ctx.RequestBody()
			item = collection.PostBody(body)

			// Edit
			err = collection.Update(itemID, item)

			if err != nil {
				return ctx.Error(http.StatusBadRequest, "Item could not be edited", err)
			}

			// Save
			err = collection.Save()

			if err != nil {
				return ctx.Error(http.StatusInternalServerError, "Collection could not be saved", err)
			}

			return "ok"
		}

		app.Post(route, handler)
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

		collection := obj.(Collection)
		err = collection.Authorize(ctx)

		if err != nil {
			return ctx.Error(http.StatusForbidden, "Not authorized", err)
		}

		body := ctx.RequestBody()
		item := collection.PostBody(body)

		err = modify(collection, item)

		if err != nil {
			return ctx.Error(http.StatusInternalServerError, "Error adding item to list", err)
		}

		err = collection.Save()

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
	filterInterface := reflect.TypeOf((*Filter)(nil)).Elem()
	filterEnabled := reflect.PtrTo(objType).Implements(filterInterface)

	route := api.root + strings.ToLower(objTypeName) + "/:id"
	handler := func(ctx *aero.Context) string {
		objID := ctx.Get("id")
		obj, err := api.db.Get(objTypeName, objID)

		if err != nil {
			return ctx.Error(http.StatusNotFound, "Not found", err)
		}

		// Remove private data
		if filterEnabled {
			filter := obj.(Filter)

			if filter.ShouldFilter(ctx) {
				filter.Filter()
			}
		}

		return ctx.JSON(obj)
	}

	return route, handler
}

// Update ...
func (api *API) Update(table string) (string, aero.Handle) {
	objType := api.db.Type(table)
	objTypeName := objType.Name()
	editableInterface := reflect.TypeOf((*Editable)(nil)).Elem()

	if !reflect.PtrTo(objType).Implements(editableInterface) {
		return "", nil
	}

	route := api.root + strings.ToLower(objTypeName) + "/:id"
	handler := func(ctx *aero.Context) string {
		objID := ctx.Get("id")
		obj, err := api.db.Get(objTypeName, objID)

		if err != nil {
			return ctx.Error(http.StatusNotFound, "Not found", err)
		}

		// Authorize
		editable := obj.(Editable)
		err = editable.Authorize(ctx)

		if err != nil {
			return ctx.Error(http.StatusForbidden, "Not authorized", err)
		}

		// Parse body
		body := ctx.RequestBody()

		var data interface{}
		err = json.Unmarshal(body, &data)

		if err != nil {
			return ctx.Error(http.StatusBadRequest, "Invalid data format (expected JSON)", err)
		}

		// Edit
		err = editable.Update(ctx, data)

		if err != nil {
			return ctx.Error(http.StatusBadRequest, objTypeName+" could not be updated", err)
		}

		// Save
		err = editable.Save()

		if err != nil {
			return ctx.Error(http.StatusInternalServerError, objTypeName+" could not be saved", err)
		}

		return "ok"
	}

	return route, handler
}

// Create ...
func (api *API) Create(table string) (string, aero.Handle) {
	objType := api.db.Type(table)
	objTypeName := objType.Name()
	creatableInterface := reflect.TypeOf((*Creatable)(nil)).Elem()

	if !reflect.PtrTo(objType).Implements(creatableInterface) {
		return "", nil
	}

	route := api.root + strings.ToLower(objTypeName) + "/new"
	handler := func(ctx *aero.Context) string {
		obj := reflect.New(objType).Interface()
		creatable := obj.(Creatable)

		// Authorize
		err := creatable.Authorize(ctx)

		if err != nil {
			return ctx.Error(http.StatusForbidden, "Not authorized", err)
		}

		// Parse body
		body := ctx.RequestBody()

		var data interface{}
		err = json.Unmarshal(body, &data)

		if err != nil {
			return ctx.Error(http.StatusBadRequest, "Invalid data format (expected JSON)", err)
		}

		// Create
		err = creatable.Create(data, ctx)

		if err != nil {
			return ctx.Error(http.StatusBadRequest, objTypeName+" could not be created", err)
		}

		// Save
		err = creatable.Save()

		if err != nil {
			return ctx.Error(http.StatusInternalServerError, objTypeName+" could not be saved", err)
		}

		return "ok"
	}

	return route, handler
}
