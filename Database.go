package api

// Database ...
type Database interface {
	Get(table string, id string) (interface{}, error)
	Set(table string, id string, obj interface{}) error
}
