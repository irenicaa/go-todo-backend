package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

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
}

// Create ...
func (handler TodoRecord) Create(
	writer http.ResponseWriter,
	request *http.Request,
) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))

		return
	}

	var todo models.TodoRecord
	if err := json.Unmarshal(body, &todo); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))

		return
	}

	baseURL := &url.URL{Scheme: handler.URLScheme, Host: request.Host}
	presentationTodo, err := handler.UseCase.Create(baseURL, todo)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))

		return
	}

	response, err := json.Marshal(presentationTodo)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))

		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(response)
}
