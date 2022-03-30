package handlers

import (
	"math"
	"net/http"
	"net/url"

	httputils "github.com/irenicaa/go-http-utils"
	"github.com/irenicaa/go-todo-backend/v2/models"
)

// TodoRecordUseCase ...
type TodoRecordUseCase interface {
	GetAll(baseURL *url.URL, query models.Query) (
		[]models.PresentationTodoRecord,
		error,
	)
	GetSingle(baseURL *url.URL, id int) (models.PresentationTodoRecord, error)
	Create(baseURL *url.URL, presentationTodo models.PresentationTodoRecord) (
		models.PresentationTodoRecord,
		error,
	)
	Update(
		baseURL *url.URL,
		id int,
		presentationTodo models.PresentationTodoRecord,
	) (
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
//   @router /todos [GET]
//   @summary get all to-do records
//   @param minimal_date query string false "filtration by the minimal date in the RFC 3339 format"
//   @param maximal_date query string false "filtration by the maximal date in the RFC 3339 format"
//   @param title_fragment query string false "search by the title fragment"
//   @param page_size query integer false "specify the page size for pagination" minimum(1)
//   @param page query integer false "specify the page for pagination" minimum(1)
//   @produce json
//   @success 200 {array} models.PresentationTodoRecord
//   @failure 500 {string} string
func (handler TodoRecord) GetAll(
	writer http.ResponseWriter,
	request *http.Request,
) {
	minimalDate, err := httputils.GetDateFormValue(request, "minimal_date")
	if err != nil && err != httputils.ErrKeyIsMissed {
		status, message :=
			http.StatusBadRequest, "unable to get the minimal_date parameter: %v"
		httputils.HandleError(writer, handler.Logger, status, message, err)

		return
	}

	maximalDate, err := httputils.GetDateFormValue(request, "maximal_date")
	if err != nil && err != httputils.ErrKeyIsMissed {
		status, message :=
			http.StatusBadRequest, "unable to get the maximal_date parameter: %v"
		httputils.HandleError(writer, handler.Logger, status, message, err)

		return
	}

	pageSize, err :=
		httputils.GetIntFormValue(request, "page_size", 1, math.MaxInt32)
	if err != nil && err != httputils.ErrKeyIsMissed {
		status, message :=
			http.StatusBadRequest, "unable to get the page_size parameter: %v"
		httputils.HandleError(writer, handler.Logger, status, message, err)

		return
	}

	page, err := httputils.GetIntFormValue(request, "page", 1, math.MaxInt32)
	if err != nil && err != httputils.ErrKeyIsMissed {
		status, message :=
			http.StatusBadRequest, "unable to get the page parameter: %v"
		httputils.HandleError(writer, handler.Logger, status, message, err)

		return
	}

	baseURL := handler.getBaseURL(request)
	presentationTodos, err := handler.UseCase.GetAll(baseURL, models.Query{
		MinimalDate:   minimalDate,
		MaximalDate:   maximalDate,
		TitleFragment: request.FormValue("title_fragment"),
		Pagination:    models.Pagination{PageSize: pageSize, Page: page},
	})
	if err != nil {
		status, message := http.StatusInternalServerError, "%s"
		httputils.HandleError(writer, handler.Logger, status, message, err)

		return
	}

	httputils.HandleJSON(writer, handler.Logger, presentationTodos)
}

// GetAllByDate ...
//   @router /todos/{date} [GET]
//   @summary get all to-do records
//   @param date path string true "to-do record date in the RFC 3339 format"
//   @param title_fragment query string false "search by the title fragment"
//   @param page_size query integer false "specify the page size for pagination" minimum(1)
//   @param page query integer false "specify the page for pagination" minimum(1)
//   @produce json
//   @success 200 {array} models.PresentationTodoRecord
//   @failure 500 {string} string
func (handler TodoRecord) GetAllByDate(
	writer http.ResponseWriter,
	request *http.Request,
) {
	date, err := httputils.GetDateFromURL(request)
	if err != nil {
		status, message := http.StatusBadRequest, "unable to get a date: %s"
		httputils.HandleError(writer, handler.Logger, status, message, err)

		return
	}

	pageSize, err :=
		httputils.GetIntFormValue(request, "page_size", 1, math.MaxInt32)
	if err != nil && err != httputils.ErrKeyIsMissed {
		status, message :=
			http.StatusBadRequest, "unable to get the page_size parameter: %v"
		httputils.HandleError(writer, handler.Logger, status, message, err)

		return
	}

	page, err := httputils.GetIntFormValue(request, "page", 1, math.MaxInt32)
	if err != nil && err != httputils.ErrKeyIsMissed {
		status, message :=
			http.StatusBadRequest, "unable to get the page parameter: %v"
		httputils.HandleError(writer, handler.Logger, status, message, err)

		return
	}

	baseURL := handler.getBaseURL(request)
	presentationTodos, err := handler.UseCase.GetAll(baseURL, models.Query{
		MinimalDate:   date,
		MaximalDate:   date,
		TitleFragment: request.FormValue("title_fragment"),
		Pagination:    models.Pagination{PageSize: pageSize, Page: page},
	})
	if err != nil {
		status, message := http.StatusInternalServerError, "%s"
		httputils.HandleError(writer, handler.Logger, status, message, err)

		return
	}

	httputils.HandleJSON(writer, handler.Logger, presentationTodos)
}

// GetSingle ...
//   @router /todos/{id} [GET]
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
		status, message := http.StatusBadRequest, "unable to get an ID: %s"
		httputils.HandleError(writer, handler.Logger, status, message, err)

		return
	}

	baseURL := handler.getBaseURL(request)
	presentationTodo, err := handler.UseCase.GetSingle(baseURL, id)
	if err != nil {
		status, message := http.StatusInternalServerError, "%s"
		httputils.HandleError(writer, handler.Logger, status, message, err)

		return
	}

	httputils.HandleJSON(writer, handler.Logger, presentationTodo)
}

// Create ...
//   @router /todos [POST]
//   @summary create a to-do record
//   @param body body models.PresentationTodoRecord true "to-do record data"
//   @accept json
//   @produce json
//   @success 200 {object} models.PresentationTodoRecord
//   @failure 400 {string} string
//   @failure 500 {string} string
func (handler TodoRecord) Create(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var presentationTodo models.PresentationTodoRecord
	if err := httputils.ReadJSONData(request.Body, &presentationTodo); err != nil {
		status, message := http.StatusBadRequest, "unable to get the request body: %s"
		httputils.HandleError(writer, handler.Logger, status, message, err)

		return
	}

	baseURL := handler.getBaseURL(request)
	presentationTodo, err := handler.UseCase.Create(baseURL, presentationTodo)
	if err != nil {
		status, message := http.StatusInternalServerError, "%s"
		httputils.HandleError(writer, handler.Logger, status, message, err)

		return
	}

	httputils.HandleJSON(writer, handler.Logger, presentationTodo)
}

// Update ...
//   @router /todos/{id} [PUT]
//   @summary update the to-do record
//   @param id path integer true "to-do record ID"
//   @param body body models.PresentationTodoRecord true "to-do record data"
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
		status, message := http.StatusBadRequest, "unable to get an ID: %s"
		httputils.HandleError(writer, handler.Logger, status, message, err)

		return
	}

	var presentationTodo models.PresentationTodoRecord
	if err := httputils.ReadJSONData(request.Body, &presentationTodo); err != nil {
		status, message := http.StatusBadRequest, "unable to get the request body: %s"
		httputils.HandleError(writer, handler.Logger, status, message, err)

		return
	}

	baseURL := handler.getBaseURL(request)
	presentationTodo, err = handler.UseCase.Update(baseURL, id, presentationTodo)
	if err != nil {
		status, message := http.StatusInternalServerError, "%s"
		httputils.HandleError(writer, handler.Logger, status, message, err)

		return
	}

	httputils.HandleJSON(writer, handler.Logger, presentationTodo)
}

// Patch ...
//   @router /todos/{id} [PATCH]
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
		status, message := http.StatusBadRequest, "unable to get an ID: %s"
		httputils.HandleError(writer, handler.Logger, status, message, err)

		return
	}

	var todoPatch models.TodoRecordPatch
	if err := httputils.ReadJSONData(request.Body, &todoPatch); err != nil {
		status, message := http.StatusBadRequest, "unable to get the request body: %s"
		httputils.HandleError(writer, handler.Logger, status, message, err)

		return
	}

	baseURL := handler.getBaseURL(request)
	presentationTodo, err := handler.UseCase.Patch(baseURL, id, todoPatch)
	if err != nil {
		status, message := http.StatusInternalServerError, "%s"
		httputils.HandleError(writer, handler.Logger, status, message, err)

		return
	}

	httputils.HandleJSON(writer, handler.Logger, presentationTodo)
}

// DeleteAll ...
//   @router /todos [DELETE]
//   @summary delete the to-do records
//   @success 204 {string} string
//   @failure 500 {string} string
func (handler TodoRecord) DeleteAll(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if err := handler.UseCase.DeleteAll(); err != nil {
		status, message := http.StatusInternalServerError, "%s"
		httputils.HandleError(writer, handler.Logger, status, message, err)

		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

// DeleteSingle ...
//   @router /todos/{id} [DELETE]
//   @summary delete the to-do record
//   @param id path integer true "to-do record ID"
//   @success 204 {string} string
//   @failure 400 {string} string
//   @failure 500 {string} string
func (handler TodoRecord) DeleteSingle(
	writer http.ResponseWriter,
	request *http.Request,
) {
	id, err := httputils.GetIDFromURL(request)
	if err != nil {
		status, message := http.StatusBadRequest, "unable to get an ID: %s"
		httputils.HandleError(writer, handler.Logger, status, message, err)

		return
	}

	if err := handler.UseCase.DeleteSingle(id); err != nil {
		status, message := http.StatusInternalServerError, "%s"
		httputils.HandleError(writer, handler.Logger, status, message, err)

		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (handler TodoRecord) getBaseURL(request *http.Request) *url.URL {
	return &url.URL{Scheme: handler.URLScheme, Host: request.Host}
}
