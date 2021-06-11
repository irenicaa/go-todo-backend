package models

// TodoRecordPatch ...
type TodoRecordPatch struct {
	Date      *Date   `json:"date"`
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
	Order     *int    `json:"order"`
}
