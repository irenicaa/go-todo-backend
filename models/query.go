package models

// Query ...
type Query struct {
	TitleFragment string
	Pagination    Pagination
}

// Pagination ...
type Pagination struct {
	PageSize int
	Page     int
}
