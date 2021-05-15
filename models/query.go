package models

// Query ...
type Query struct {
	MinimalDate   Date
	MaximalDate   Date
	TitleFragment string
	Pagination    Pagination
}

// Pagination ...
type Pagination struct {
	PageSize int
	Page     int
}
