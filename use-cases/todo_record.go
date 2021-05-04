package usecases

import (
	"fmt"
	"net/url"

	"github.com/irenicaa/go-todo-backend/models"
)

// TodoRecordStorage ...
type TodoRecordStorage interface {
	GetAll(query models.Query) ([]models.TodoRecord, error)
	GetSingle(id int) (models.TodoRecord, error)
	Create(todo models.TodoRecord) (id int, err error)
	Update(id int, todo models.TodoRecord) error
	DeleteAll() error
	DeleteSingle(id int) error
}

// TodoRecord ...
type TodoRecord struct {
	Storage TodoRecordStorage
}

// GetAll ...
func (useCase TodoRecord) GetAll(baseURL *url.URL, query models.Query) (
	[]models.PresentationTodoRecord,
	error,
) {
	todos, err := useCase.Storage.GetAll(query)
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
func (useCase TodoRecord) Create(
	baseURL *url.URL,
	presentationTodo models.PresentationTodoRecord,
) (
	models.PresentationTodoRecord,
	error,
) {
	todo := models.NewTodoRecord(presentationTodo)
	id, err := useCase.Storage.Create(todo)
	if err != nil {
		return models.PresentationTodoRecord{},
			fmt.Errorf("unable to create a to-do record: %v", err)
	}

	todo.ID = id

	presentationTodo = models.NewPresentationTodoRecord(baseURL, todo)
	return presentationTodo, nil
}

// Update ...
func (useCase TodoRecord) Update(
	baseURL *url.URL,
	id int,
	presentationTodo models.PresentationTodoRecord,
) (
	models.PresentationTodoRecord,
	error,
) {
	todo := models.NewTodoRecord(presentationTodo)
	if err := useCase.Storage.Update(id, todo); err != nil {
		return models.PresentationTodoRecord{},
			fmt.Errorf("unable to update the to-do record: %v", err)
	}

	todo.ID = id

	presentationTodo = models.NewPresentationTodoRecord(baseURL, todo)
	return presentationTodo, nil
}

// Patch ...
func (useCase TodoRecord) Patch(
	baseURL *url.URL,
	id int,
	todoPatch models.TodoRecordPatch,
) (
	models.PresentationTodoRecord,
	error,
) {
	todo, err := useCase.Storage.GetSingle(id)
	if err != nil {
		return models.PresentationTodoRecord{},
			fmt.Errorf("unable to get the to-do record: %v", err)
	}

	todo.Patch(todoPatch)

	presentationTodo := models.NewPresentationTodoRecord(baseURL, todo)
	return useCase.Update(baseURL, id, presentationTodo)
}

// DeleteAll ...
func (useCase TodoRecord) DeleteAll() error {
	if err := useCase.Storage.DeleteAll(); err != nil {
		return fmt.Errorf("unable to delete the to-do records: %v", err)
	}

	return nil
}

// DeleteSingle ...
func (useCase TodoRecord) DeleteSingle(id int) error {
	if err := useCase.Storage.DeleteSingle(id); err != nil {
		return fmt.Errorf("unable to delete the to-do record: %v", err)
	}

	return nil
}
