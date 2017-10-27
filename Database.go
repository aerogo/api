package api

import "reflect"

// Database ...
type Database interface {
	Get(table string, id string) (interface{}, error)
	Set(table string, id string, obj interface{})
	Delete(table string, id string) bool
	Types() map[string]reflect.Type
}
