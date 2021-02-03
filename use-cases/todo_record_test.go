package usecases

import (
	"net/url"
	"testing"
	"testing/iotest"

	"github.com/irenicaa/go-todo-backend/models"
	"github.com/stretchr/testify/assert"
)

func TestTodoRecord_Create(t *testing.T) {
	type fields struct {
		Storage Storage
	}
	type args struct {
		baseURL *url.URL
		todo    models.TodoRecord
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.PresentationTodoRecord
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				Storage: func() Storage {
					todo := models.TodoRecord{
						ID:        23,
						Title:     "test",
						Completed: true,
						Order:     42,
					}

					storage := &MockStorage{}
					storage.InnerMock.On("Create", todo).Return(42, nil)

					return storage
				}(),
			},
			args: args{
				baseURL: &url.URL{Scheme: "https", Host: "example.com"},
				todo: models.TodoRecord{
					ID:        23,
					Title:     "test",
					Completed: true,
					Order:     42,
				},
			},
			want: models.PresentationTodoRecord{
				URL:       "https://example.com/api/v1/todos/42",
				Title:     "test",
				Completed: true,
				Order:     42,
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				Storage: func() Storage {
					todo := models.TodoRecord{
						ID:        23,
						Title:     "test",
						Completed: true,
						Order:     42,
					}

					storage := &MockStorage{}
					storage.InnerMock.On("Create", todo).Return(0, iotest.ErrTimeout)

					return storage
				}(),
			},
			args: args{
				baseURL: &url.URL{Scheme: "https", Host: "example.com"},
				todo: models.TodoRecord{
					ID:        23,
					Title:     "test",
					Completed: true,
					Order:     42,
				},
			},
			want:    models.PresentationTodoRecord{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := TodoRecord{
				Storage: tt.fields.Storage,
			}
			got, err := useCase.Create(tt.args.baseURL, tt.args.todo)

			tt.fields.Storage.(*MockStorage).InnerMock.AssertExpectations(t)
			assert.Equal(t, tt.want, got)
			tt.wantErr(t, err)
		})
	}
}
