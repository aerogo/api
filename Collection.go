package api

// Collection interface for lists, sets, etc.
type Collection interface {
	Authorizable
	Savable
	PostBodyReader

	Add(interface{}) error
	Remove(interface{}) bool
	Contains(interface{}) bool

	Get(id interface{}) (interface{}, error)
	Set(id interface{}, value interface{}) error
	Update(id interface{}, updates interface{}) error
}
