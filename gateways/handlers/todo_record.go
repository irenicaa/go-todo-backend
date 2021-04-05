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
	DeleteAll() error
	DeleteSingle(id int) error
}

// TodoRecord ...
type TodoRecord struct {
	URLScheme string
	UseCase   TodoRecordUseCase
	Logger    httputils.Logger
}

// GetAll ...
//   @router /api/v1/todos [GET]
//   @summary get all to-do records
//   @produce json
//   @success 200 {array} models.PresentationTodoRecord
//   @failure 500 {string} string
func (handler TodoRecord) GetAll(
	writer http.ResponseWriter,
	request *http.Request,
) {
	baseURL := handler.getBaseURL(request)
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
//   @router /api/v1/todos/{id} [GET]
//   @summary get the single to-do record
//   @param id path integer true "to-do record ID"
//   @produce json
//   @success 200 {object} models.PresentationTodoRecord
//   @failure 400 {string} string
//   @failure 500 {string} string
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

	baseURL := handler.getBaseURL(request)
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
//   @router /api/v1/todos [POST]
//   @summary create a to-do record
//   @param body body models.TodoRecord true "to-do record data"
//   @accept json
//   @produce json
//   @success 200 {object} models.PresentationTodoRecord
//   @failure 400 {string} string
//   @failure 500 {string} string
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

	baseURL := handler.getBaseURL(request)
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
//   @router /api/v1/todos/{id} [PUT]
//   @summary update the to-do record
//   @param id path integer true "to-do record ID"
//   @param body body models.TodoRecord true "to-do record data"
//   @accept json
//   @produce json
//   @success 200 {object} models.PresentationTodoRecord
//   @failure 400 {string} string
//   @failure 500 {string} string
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

	baseURL := handler.getBaseURL(request)
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
//   @router /api/v1/todos/{id} [PATCH]
//   @summary patch the to-do record
//   @param id path integer true "to-do record ID"
//   @param body body models.TodoRecordPatch true "to-do record patch"
//   @accept json
//   @produce json
//   @success 200 {object} models.PresentationTodoRecord
//   @failure 400 {string} string
//   @failure 500 {string} string
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

	baseURL := handler.getBaseURL(request)
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

// DeleteAll ...
//   @router /api/v1/todos [DELETE]
//   @summary delete the to-do records
//   @produce json
//   @success 204 {string} string
//   @failure 400 {string} string
//   @failure 500 {string} string
func (handler TodoRecord) DeleteAll(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if err := handler.UseCase.DeleteAll(); err != nil {
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

// DeleteSingle ...
//   @router /api/v1/todos/{id} [DELETE]
//   @summary delete the to-do record
//   @param id path integer true "to-do record ID"
//   @produce json
//   @success 204 {string} string
//   @failure 400 {string} string
//   @failure 500 {string} string
func (handler TodoRecord) DeleteSingle(
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

	if err := handler.UseCase.DeleteSingle(id); err != nil {
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

func (handler TodoRecord) getBaseURL(request *http.Request) *url.URL {
	return &url.URL{Scheme: handler.URLScheme, Host: request.Host}
}
