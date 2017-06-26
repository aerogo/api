package api

// Creatable defines an object type where new instances can be created by users.
type Creatable interface {
	Savable
	Authorizable
	Create(data interface{}) error
}
