package api

import "reflect"

// Database ...
type Database interface {
	Get(table string, id string) (interface{}, error)
	Set(table string, id string, obj interface{}) error
	Delete(table string, id string) (existed bool, err error)
	Type(table string) reflect.Type
	Types() map[string]reflect.Type
}
