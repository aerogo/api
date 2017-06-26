package api

// Collection interface for lists, sets, etc.
type Collection interface {
	Authorizable
	Savable

	Add(interface{}) error
	Remove(interface{}) bool
	Contains(interface{}) bool

	// PostBody reads the POST body and returns an object
	// that is passed to methods like Update, Add, Remove, etc.
	PostBody(body []byte) interface{}

	Get(id interface{}) (interface{}, error)
	Set(id interface{}, value interface{}) error
	Update(id interface{}, updates interface{}) error
}
