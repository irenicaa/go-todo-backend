package usecases

import (
	"fmt"
	"net/url"

	"github.com/irenicaa/go-todo-backend/models"
)

// Storage ...
type Storage interface {
	Create(todo models.TodoRecord) (id int, err error)
	Update(id int, todo models.TodoRecord) error
}

// TodoRecord ...
type TodoRecord struct {
	Storage Storage
}

// Create ...
func (useCase TodoRecord) Create(baseURL *url.URL, todo models.TodoRecord) (
	models.PresentationTodoRecord,
	error,
) {
	id, err := useCase.Storage.Create(todo)
	if err != nil {
		return models.PresentationTodoRecord{},
			fmt.Errorf("unable to create a to-do record: %v", err)
	}

	todo.ID = id

	presentationTodo := models.NewPresentationTodoRecord(baseURL, todo)
	return presentationTodo, nil
}

// Update ...
func (useCase TodoRecord) Update(
	baseURL *url.URL,
	id int,
	todo models.TodoRecord,
) (
	models.PresentationTodoRecord,
	error,
) {
	err := useCase.Storage.Update(id, todo)
	if err != nil {
		return models.PresentationTodoRecord{},
			fmt.Errorf("unable to update a to-do record: %v", err)
	}

	todo.ID = id

	presentationTodo := models.NewPresentationTodoRecord(baseURL, todo)
	return presentationTodo, nil
}
