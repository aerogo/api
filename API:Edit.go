package api

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/aerogo/aero"
	"github.com/aerogo/mirror"
)

// Edit ...
func (api *API) Edit(table string) (string, aero.Handle) {
	objType := api.db.Type(table)
	objTypeName := objType.Name()
	editableInterface := reflect.TypeOf((*Editable)(nil)).Elem()

	if !reflect.PtrTo(objType).Implements(editableInterface) {
		return "", nil
	}

	customEditableInterface := reflect.TypeOf((*CustomEditable)(nil)).Elem()
	afterEditableInterface := reflect.TypeOf((*AfterEditable)(nil)).Elem()
	virtualEditableInterface := reflect.TypeOf((*VirtualEditable)(nil)).Elem()

	usesCustomEdits := reflect.PtrTo(objType).Implements(customEditableInterface)
	usesAfterEdits := reflect.PtrTo(objType).Implements(afterEditableInterface)
	usesVirtualEdits := reflect.PtrTo(objType).Implements(virtualEditableInterface)

	route := api.root + strings.ToLower(objTypeName) + "/:id"
	handler := func(ctx *aero.Context) string {
		objID := ctx.Get("id")
		obj, err := api.db.Get(objTypeName, objID)

		if err != nil {
			return ctx.Error(http.StatusNotFound, "Not found", err)
		}

		// Authorize
		editable := obj.(Editable)
		err = editable.Authorize(ctx, "edit")

		if err != nil {
			return ctx.Error(http.StatusForbidden, "Not authorized", err)
		}

		// Parse body
		edits, err := ctx.RequestBodyJSONObject()

		if err != nil {
			return ctx.Error(http.StatusBadRequest, "Invalid data format (expected JSON)", err)
		}

		// Apply changes
		for key, value := range edits {
			newValue := reflect.ValueOf(value)

			// Virtual properties
			if usesVirtualEdits {
				virtualEditable := editable.(VirtualEditable)
				consumed, err := virtualEditable.VirtualEdit(ctx, key, newValue)

				if err != nil {
					return ctx.Error(http.StatusBadRequest, objTypeName+" "+key+" could not be edited", err)
				}

				if consumed {
					continue
				}
			}

			field, _, v, err := mirror.GetProperty(editable, key)

			if err != nil {
				return ctx.Error(http.StatusBadRequest, objTypeName+" "+key+" could not be edited", err)
			}

			// Is somebody attempting to edit fields that aren't editable?
			if field.Tag.Get("editable") != "true" {
				return ctx.Error(http.StatusBadRequest, objTypeName+" "+key+" could not be edited", errors.New("Field "+key+" is not editable"))
			}

			if !v.CanSet() {
				return ctx.Error(http.StatusBadRequest, objTypeName+" "+key+" could not be edited", errors.New("Field "+key+" is not settable"))
			}

			// Special edit
			if usesCustomEdits {
				customEditable := editable.(CustomEditable)
				consumed, err := customEditable.Edit(ctx, key, v, newValue)

				if err != nil {
					return ctx.Error(http.StatusBadRequest, objTypeName+" "+key+" could not be edited", err)
				}

				if consumed {
					continue
				}
			}

			// Implement special data type cases here
			switch v.Kind() {
			case reflect.Int:
				x := int64(newValue.Float())

				if !v.OverflowInt(x) {
					v.SetInt(x)
				} else {
					return ctx.Error(http.StatusBadRequest, objTypeName+" "+key+" could not be edited", errors.New("Field "+key+" would cause an integer overflow"))
				}

			default:
				v.Set(newValue)
			}
		}

		// AfterEdit
		if usesAfterEdits {
			afterEditable := editable.(AfterEditable)
			err := afterEditable.AfterEdit(ctx)

			if err != nil {
				return ctx.Error(http.StatusInternalServerError, objTypeName+" could not be edited", err)
			}
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
