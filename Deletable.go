package api

// Deletable defines an object type that can be deleted from the database.
type Deletable interface {
	Authorizable
	Delete() error
}
