package models

import utilmodels "github.com/irenicaa/go-http-utils/models"

// TodoRecordPatch ...
type TodoRecordPatch struct {
	Date      *utilmodels.Date `json:"date"`
	Title     *string          `json:"title"`
	Completed *bool            `json:"completed"`
	Order     *int             `json:"order"`
}
