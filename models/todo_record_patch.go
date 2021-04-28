package models

import "time"

// TodoRecordPatch ...
type TodoRecordPatch struct {
	Date      *time.Time
	Title     *string
	Completed *bool
	Order     *int
}
