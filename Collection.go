package api

// Collection interface for lists, sets, etc.
type Collection interface {
	Add(interface{}) error
	Remove(interface{}) bool
	Contains(interface{}) bool

	// TransformBody returns an item that is passed to methods like Add, Remove, etc.
	TransformBody(body []byte) interface{}

	// Save saves the object in the database.
	Save() error
}
