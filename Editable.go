package api

// An Editable can authorize changes, be changed and be saved in the database.
type Editable interface {
	Authorizable
	Savable
	Updatable
	PostBodyReader
}
