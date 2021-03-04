package handlers

import (
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

// GetAll ...
func (handler TodoRecord) GetAll(
	writer http.ResponseWriter,
	request *http.Request,
) {
	baseURL := &url.URL{Scheme: handler.URLScheme, Host: request.Host}
	presentationTodo, err := handler.UseCase.GetAll(baseURL)
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

// GetSingle ...
func (handler TodoRecord) GetSingle(
	writer http.ResponseWriter,
	request *http.Request,
) {
	id, err := httputils.GetIDFromURL(request)
	if err != nil {
		httputils.HandleError(
			writer,
			handler.Logger,
			http.StatusBadRequest,
			"unable to get an ID: %s",
			err,
		)

		return
	}

	baseURL := &url.URL{Scheme: handler.URLScheme, Host: request.Host}
	presentationTodo, err := handler.UseCase.GetSingle(baseURL, id)
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

// Create ...
func (handler TodoRecord) Create(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var todo models.TodoRecord
	if err := httputils.GetRequestBody(request, &todo); err != nil {
		httputils.HandleError(
			writer,
			handler.Logger,
			http.StatusBadRequest,
			"unable to get the request body: %s",
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
	id, err := httputils.GetIDFromURL(request)
	if err != nil {
		httputils.HandleError(
			writer,
			handler.Logger,
			http.StatusBadRequest,
			"unable to get an ID: %s",
			err,
		)

		return
	}

	var todo models.TodoRecord
	if err := httputils.GetRequestBody(request, &todo); err != nil {
		httputils.HandleError(
			writer,
			handler.Logger,
			http.StatusBadRequest,
			"unable to get the request body: %s",
			err,
		)

		return
	}

	baseURL := &url.URL{Scheme: handler.URLScheme, Host: request.Host}
	presentationTodo, err := handler.UseCase.Update(baseURL, id, todo)
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
	id, err := httputils.GetIDFromURL(request)
	if err != nil {
		httputils.HandleError(
			writer,
			handler.Logger,
			http.StatusBadRequest,
			"unable to get an ID: %s",
			err,
		)

		return
	}

	var todoPatch models.TodoRecordPatch
	if err := httputils.GetRequestBody(request, &todoPatch); err != nil {
		httputils.HandleError(
			writer,
			handler.Logger,
			http.StatusBadRequest,
			"unable to get the request body: %s",
			err,
		)

		return
	}

	baseURL := &url.URL{Scheme: handler.URLScheme, Host: request.Host}
	presentationTodo, err := handler.UseCase.Patch(baseURL, id, todoPatch)
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

// Delete ...
func (handler TodoRecord) Delete(
	writer http.ResponseWriter,
	request *http.Request,
) {
	id, err := httputils.GetIDFromURL(request)
	if err != nil {
		httputils.HandleError(
			writer,
			handler.Logger,
			http.StatusBadRequest,
			"unable to get an ID: %s",
			err,
		)

		return
	}

	if err := handler.UseCase.Delete(id); err != nil {
		httputils.HandleError(
			writer,
			handler.Logger,
			http.StatusInternalServerError,
			"%s",
			err,
		)

		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
