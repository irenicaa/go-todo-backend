package models

// TodoRecord ...
type TodoRecord struct {
	ID        int
	Title     string
	Completed bool
	Order     int
}

// Patch ...
func (todo *TodoRecord) Patch(patch TodoRecordPatch) {
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
