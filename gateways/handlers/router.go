package handlers

import (
	"net/http"
	"strings"

	httputils "github.com/irenicaa/go-http-utils"
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
			return
		case http.MethodGet:
			if request.URL.Path == router.BaseURL+"/todos" {
				router.TodoRecord.GetAll(writer, request)
			} else if httputils.DatePattern.MatchString(request.URL.Path) {
				router.TodoRecord.GetAllByDate(writer, request)
			} else {
				router.TodoRecord.GetSingle(writer, request)
			}

			return
		case http.MethodPut:
			router.TodoRecord.Update(writer, request)
			return
		case http.MethodPatch:
			router.TodoRecord.Patch(writer, request)
			return
		case http.MethodDelete:
			if request.URL.Path == router.BaseURL+"/todos" {
				router.TodoRecord.DeleteAll(writer, request)
			} else {
				router.TodoRecord.DeleteSingle(writer, request)
			}

			return
		}
	}

	status, message := http.StatusNotFound, http.StatusText(http.StatusNotFound)
	httputils.HandleError(writer, router.Logger, status, message)
}
