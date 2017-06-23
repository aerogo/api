package api

// Updatable objects can be updated with new data via POST requests.
type Updatable interface {
	Update(data interface{}) error
}
