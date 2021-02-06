package usecases

import (
	"fmt"
	"net/url"

	"github.com/irenicaa/go-todo-backend/models"
)

// Storage ...
type Storage interface {
	GetAll() ([]models.TodoRecord, error)
	GetSingle(id int) (models.TodoRecord, error)
	Create(todo models.TodoRecord) (id int, err error)
	Update(id int, todo models.TodoRecord) error
	Delete(id int) error
}

// TodoRecord ...
type TodoRecord struct {
	Storage Storage
}

// GetAll ...
func (useCase TodoRecord) GetAll(baseURL *url.URL) (
	[]models.PresentationTodoRecord,
	error,
) {
	todos, err := useCase.Storage.GetAll()
	if err != nil {
		return nil, fmt.Errorf("unable to get the to-do records: %v", err)
	}

	var presentationTodos []models.PresentationTodoRecord
	for _, record := range todos {
		presentationTodo := models.NewPresentationTodoRecord(baseURL, record)
		presentationTodos = append(presentationTodos, presentationTodo)
	}

	return presentationTodos, nil
}

// GetSingle ...
func (useCase TodoRecord) GetSingle(baseURL *url.URL, id int) (
	models.PresentationTodoRecord,
	error,
) {
	todo, err := useCase.Storage.GetSingle(id)
	if err != nil {
		return models.PresentationTodoRecord{},
			fmt.Errorf("unable to get the to-do record: %v", err)
	}

	presentationTodo := models.NewPresentationTodoRecord(baseURL, todo)
	return presentationTodo, nil
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
	if err := useCase.Storage.Update(id, todo); err != nil {
		return models.PresentationTodoRecord{},
			fmt.Errorf("unable to update the to-do record: %v", err)
	}

	todo.ID = id

	presentationTodo := models.NewPresentationTodoRecord(baseURL, todo)
	return presentationTodo, nil
}

// Delete ...
func (useCase TodoRecord) Delete(id int) error {
	if err := useCase.Storage.Delete(id); err != nil {
		return fmt.Errorf("unable to delete the to-do record: %v", err)
	}

	return nil
}
