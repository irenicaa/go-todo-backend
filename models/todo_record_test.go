package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTodoRecord_Patch(t *testing.T) {
	type fields struct {
		ID        int
		Date      time.Time
		Title     string
		Completed bool
		Order     int
	}
	type args struct {
		patch TodoRecordPatch
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantTodo *TodoRecord
	}{
		{
			name: "updating of all fields",
			fields: fields{
				ID:        23,
				Date:      time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC),
				Title:     "test",
				Completed: true,
				Order:     42,
			},
			args: args{
				patch: TodoRecordPatch{
					Date: func() *time.Time {
						date := time.Date(2006, time.January, 3, 0, 0, 0, 0, time.UTC)
						return &date
					}(),
					Title: func() *string {
						title := "test2"
						return &title
					}(),
					Completed: func() *bool {
						completed := false
						return &completed
					}(),
					Order: func() *int {
						order := 43
						return &order
					}(),
				},
			},
			wantTodo: &TodoRecord{
				ID:        23,
				Date:      time.Date(2006, time.January, 3, 0, 0, 0, 0, time.UTC),
				Title:     "test2",
				Completed: false,
				Order:     43,
			},
		},
		{
			name: "updating of all fields except a date",
			fields: fields{
				ID:        23,
				Date:      time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC),
				Title:     "test",
				Completed: true,
				Order:     42,
			},
			args: args{
				patch: TodoRecordPatch{
					Date: nil,
					Title: func() *string {
						title := "test2"
						return &title
					}(),
					Completed: func() *bool {
						completed := false
						return &completed
					}(),
					Order: func() *int {
						order := 43
						return &order
					}(),
				},
			},
			wantTodo: &TodoRecord{
				ID:        23,
				Date:      time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC),
				Title:     "test2",
				Completed: false,
				Order:     43,
			},
		},
		{
			name: "updating of all fields except a title",
			fields: fields{
				ID:        23,
				Date:      time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC),
				Title:     "test",
				Completed: true,
				Order:     42,
			},
			args: args{
				patch: TodoRecordPatch{
					Date: func() *time.Time {
						date := time.Date(2006, time.January, 3, 0, 0, 0, 0, time.UTC)
						return &date
					}(),
					Title: nil,
					Completed: func() *bool {
						completed := false
						return &completed
					}(),
					Order: func() *int {
						order := 43
						return &order
					}(),
				},
			},
			wantTodo: &TodoRecord{
				ID:        23,
				Date:      time.Date(2006, time.January, 3, 0, 0, 0, 0, time.UTC),
				Title:     "test",
				Completed: false,
				Order:     43,
			},
		},
		{
			name: "updating of all fields except a completion flag",
			fields: fields{
				ID:        23,
				Date:      time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC),
				Title:     "test",
				Completed: true,
				Order:     42,
			},
			args: args{
				patch: TodoRecordPatch{
					Date: func() *time.Time {
						date := time.Date(2006, time.January, 3, 0, 0, 0, 0, time.UTC)
						return &date
					}(),
					Title: func() *string {
						title := "test2"
						return &title
					}(),
					Completed: nil,
					Order: func() *int {
						order := 43
						return &order
					}(),
				},
			},
			wantTodo: &TodoRecord{
				ID:        23,
				Date:      time.Date(2006, time.January, 3, 0, 0, 0, 0, time.UTC),
				Title:     "test2",
				Completed: true,
				Order:     43,
			},
		},
		{
			name: "updating of all fields except an order",
			fields: fields{
				ID:        23,
				Date:      time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC),
				Title:     "test",
				Completed: true,
				Order:     42,
			},
			args: args{
				patch: TodoRecordPatch{
					Date: func() *time.Time {
						date := time.Date(2006, time.January, 3, 0, 0, 0, 0, time.UTC)
						return &date
					}(),
					Title: func() *string {
						title := "test2"
						return &title
					}(),
					Completed: func() *bool {
						completed := false
						return &completed
					}(),
					Order: nil,
				},
			},
			wantTodo: &TodoRecord{
				ID:        23,
				Date:      time.Date(2006, time.January, 3, 0, 0, 0, 0, time.UTC),
				Title:     "test2",
				Completed: false,
				Order:     42,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo := &TodoRecord{
				ID:        tt.fields.ID,
				Date:      tt.fields.Date,
				Title:     tt.fields.Title,
				Completed: tt.fields.Completed,
				Order:     tt.fields.Order,
			}
			todo.Patch(tt.args.patch)

			assert.Equal(t, tt.wantTodo, todo)
		})
	}
}
