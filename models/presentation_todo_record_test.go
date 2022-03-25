package models

import (
	"net/url"
	"testing"
	"time"

	utilmodels "github.com/irenicaa/go-http-utils/models"
	"github.com/stretchr/testify/assert"
)

func TestNewPresentationTodoRecord(t *testing.T) {
	type args struct {
		baseURL *url.URL
		todo    TodoRecord
	}

	tests := []struct {
		name string
		args args
		want PresentationTodoRecord
	}{
		{
			name: "success",
			args: args{
				baseURL: &url.URL{Scheme: "https", Host: "example.com"},
				todo: TodoRecord{
					ID:        23,
					Date:      time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC),
					Title:     "test",
					Completed: true,
					Order:     42,
				},
			},
			want: PresentationTodoRecord{
				URL:       "https://example.com/api/v1/todos/23",
				Date:      utilmodels.Date(time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC)),
				Title:     "test",
				Completed: true,
				Order:     42,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPresentationTodoRecord(tt.args.baseURL, tt.args.todo)

			assert.Equal(t, tt.want, got)
		})
	}
}
