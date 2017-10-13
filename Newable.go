package api

// Newable defines an object type where new instances can be created by users and saved in the database.
type Newable interface {
	Savable
	Authorizable
	Creatable
}
