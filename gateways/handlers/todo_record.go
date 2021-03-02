package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	httputils "github.com/irenicaa/go-todo-backend/http-utils"
	"github.com/irenicaa/go-todo-backend/models"
)

// TodoRecordUseCase ...
type TodoRecordUseCase interface {
	GetAll(baseURL *url.URL) ([]models.PresentationTodoRecord, error)
	GetSingle(baseURL *url.URL, id int) (models.PresentationTodoRecord, error)
	Create(baseURL *url.URL, todo models.TodoRecord) (
		models.PresentationTodoRecord,
		error,
	)
	Update(baseURL *url.URL, id int, todo models.TodoRecord) (
		models.PresentationTodoRecord,
		error,
	)
	Patch(baseURL *url.URL, id int, todoPatch models.TodoRecordPatch) (
		models.PresentationTodoRecord,
		error,
	)
	Delete(id int) error
}

// TodoRecord ...
type TodoRecord struct {
	URLScheme string
	UseCase   TodoRecordUseCase
	Logger    httputils.Logger
}

// Create ...
func (handler TodoRecord) Create(
	writer http.ResponseWriter,
	request *http.Request,
) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		httputils.HandleError(
			writer,
			handler.Logger,
			http.StatusBadRequest,
			"unable to read the request body: %s",
			err,
		)

		return
	}

	var todo models.TodoRecord
	if err := json.Unmarshal(body, &todo); err != nil {
		httputils.HandleError(
			writer,
			handler.Logger,
			http.StatusBadRequest,
			"unable to unmarshal the request body: %s",
			err,
		)

		return
	}

	baseURL := &url.URL{Scheme: handler.URLScheme, Host: request.Host}
	presentationTodo, err := handler.UseCase.Create(baseURL, todo)
	if err != nil {
		httputils.HandleError(
			writer,
			handler.Logger,
			http.StatusInternalServerError,
			"%s",
			err,
		)

		return
	}

	httputils.HandleJSON(writer, handler.Logger, presentationTodo)
}

// Update ...
func (handler TodoRecord) Update(
	writer http.ResponseWriter,
	request *http.Request,
) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		httputils.HandleError(
			writer,
			handler.Logger,
			http.StatusBadRequest,
			"unable to read the request body: %s",
			err,
		)

		return
	}

	var todo models.TodoRecord
	if err := json.Unmarshal(body, &todo); err != nil {
		httputils.HandleError(
			writer,
			handler.Logger,
			http.StatusBadRequest,
			"unable to unmarshal the request body: %s",
			err,
		)

		return
	}

	baseURL := &url.URL{Scheme: handler.URLScheme, Host: request.Host}
	presentationTodo, err := handler.UseCase.Update(baseURL, 0, todo)
	if err != nil {
		httputils.HandleError(
			writer,
			handler.Logger,
			http.StatusInternalServerError,
			"%s",
			err,
		)

		return
	}

	httputils.HandleJSON(writer, handler.Logger, presentationTodo)
}

// Patch ...
func (handler TodoRecord) Patch(
	writer http.ResponseWriter,
	request *http.Request,
) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		httputils.HandleError(
			writer,
			handler.Logger,
			http.StatusBadRequest,
			"unable to read the request body: %s",
			err,
		)

		return
	}

	var todoPatch models.TodoRecordPatch
	if err := json.Unmarshal(body, &todoPatch); err != nil {
		httputils.HandleError(
			writer,
			handler.Logger,
			http.StatusBadRequest,
			"unable to unmarshal the request body: %s",
			err,
		)

		return
	}

	baseURL := &url.URL{Scheme: handler.URLScheme, Host: request.Host}
	presentationTodo, err := handler.UseCase.Patch(baseURL, 0, todoPatch)
	if err != nil {
		httputils.HandleError(
			writer,
			handler.Logger,
			http.StatusInternalServerError,
			"%s",
			err,
		)

		return
	}

	httputils.HandleJSON(writer, handler.Logger, presentationTodo)
}
