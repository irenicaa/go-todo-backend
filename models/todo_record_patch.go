package models

// TodoRecordPatch ...
type TodoRecordPatch struct {
	Date      *Date
	Title     *string
	Completed *bool
	Order     *int
}
