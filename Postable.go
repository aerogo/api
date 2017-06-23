package api

// PostBodyReader objects can transform post body data into an object.
type PostBodyReader interface {
	// PostBody reads the POST body and returns an object
	// that is passed to methods like Update, Add, Remove, etc.
	PostBody(body []byte) interface{}
}
