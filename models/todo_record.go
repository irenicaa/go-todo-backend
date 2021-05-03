package models

import "time"

// TodoRecord ...
type TodoRecord struct {
	ID        int
	Date      time.Time
	Title     string
	Completed bool
	Order     int
}

// NewTodoRecord ...
func NewTodoRecord(presentationTodo PresentationTodoRecord) TodoRecord {
	return TodoRecord{
		Date:      time.Time(presentationTodo.Date),
		Title:     presentationTodo.Title,
		Completed: presentationTodo.Completed,
		Order:     presentationTodo.Order,
	}
}

// Patch ...
func (todo *TodoRecord) Patch(patch TodoRecordPatch) {
	if patch.Date != nil {
		todo.Date = time.Time(*patch.Date)
	}
	if patch.Title != nil {
		todo.Title = *patch.Title
	}
	if patch.Completed != nil {
		todo.Completed = *patch.Completed
	}
	if patch.Order != nil {
		todo.Order = *patch.Order
	}
}
