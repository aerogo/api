package api

// Savable objects can be saved in the database.
type Savable interface {
	// Save saves the object in the database.
	Save() error
}
