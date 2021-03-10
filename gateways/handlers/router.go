package handlers

import (
	"net/http"
	"strings"

	httputils "github.com/irenicaa/go-todo-backend/http-utils"
)

// Router ...
type Router struct {
	BaseURL    string
	TodoRecord TodoRecord
	Logger     httputils.Logger
}

// ServeHTTP ...
func (router Router) ServeHTTP(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if strings.HasPrefix(request.URL.Path, router.BaseURL+"/todos") {
		switch request.Method {
		case http.MethodPost:
			router.TodoRecord.Create(writer, request)
		case http.MethodGet:
			if request.URL.Path == router.BaseURL+"/todos" {
				router.TodoRecord.GetAll(writer, request)
			} else {
				router.TodoRecord.GetSingle(writer, request)
			}
		case http.MethodPut:
			router.TodoRecord.Update(writer, request)
		case http.MethodPatch:
			router.TodoRecord.Patch(writer, request)
		case http.MethodDelete:
			router.TodoRecord.Delete(writer, request)
		}

		return
	}

	httputils.HandleError(
		writer,
		router.Logger,
		http.StatusNotFound,
		http.StatusText(http.StatusNotFound),
	)
}
