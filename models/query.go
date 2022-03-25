package models

import utilmodels "github.com/irenicaa/go-http-utils/models"

// Query ...
type Query struct {
	MinimalDate   utilmodels.Date
	MaximalDate   utilmodels.Date
	TitleFragment string
	Pagination    Pagination
}

// Pagination ...
type Pagination struct {
	PageSize int
	Page     int
}
