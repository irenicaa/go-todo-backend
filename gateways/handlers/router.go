package handlers

import "net/http"

// Router ...
type Router struct {
	BaseURL    string
	TodoRecord TodoRecord
}

// ServeHTTP ...
func (router Router) ServeHTTP(
	writer http.ResponseWriter,
	request *http.Request,
) {
}
