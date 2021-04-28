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

// Patch ...
func (todo *TodoRecord) Patch(patch TodoRecordPatch) {
	if patch.Date != nil {
		todo.Date = *patch.Date
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
