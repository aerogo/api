package api

// Collection interface for lists, sets, etc.
type Collection interface {
	Add(interface{}) error
	Remove(interface{}) error
	Contains(interface{}) bool
}
