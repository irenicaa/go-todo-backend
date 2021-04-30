package models

import (
	"fmt"
	"net/url"
)

// PresentationTodoRecord ...
type PresentationTodoRecord struct {
	URL       string `json:"url"`
	Date      Date   `json:"date"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	Order     int    `json:"order"`
}

// NewPresentationTodoRecord ...
func NewPresentationTodoRecord(
	baseURL *url.URL,
	todo TodoRecord,
) PresentationTodoRecord {
	url :=
		fmt.Sprintf("%s://%s/api/v1/todos/%d", baseURL.Scheme, baseURL.Host, todo.ID)
	return PresentationTodoRecord{
		URL:       url,
		Date:      Date(todo.Date),
		Title:     todo.Title,
		Completed: todo.Completed,
		Order:     todo.Order,
	}
}
