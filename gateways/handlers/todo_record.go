package handlers

import (
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
	UseCase TodoRecordUseCase
}
