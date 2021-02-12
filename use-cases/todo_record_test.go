package usecases

import (
	"net/url"
	"testing"
	"testing/iotest"

	"github.com/irenicaa/go-todo-backend/models"
	"github.com/stretchr/testify/assert"
)

func TestTodoRecord_GetAll(t *testing.T) {
	type fields struct {
		Storage Storage
	}
	type args struct {
		baseURL *url.URL
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.PresentationTodoRecord
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success without todos",
			fields: fields{
				Storage: func() Storage {
					storage := &MockStorage{}
					storage.InnerMock.On("GetAll").Return([]models.TodoRecord(nil), nil)

					return storage
				}(),
			},
			args: args{
				baseURL: &url.URL{Scheme: "https", Host: "example.com"},
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "success with todos",
			fields: fields{
				Storage: func() Storage {
					todos := []models.TodoRecord{
						{
							ID:        5,
							Title:     "test",
							Completed: true,
							Order:     12,
						},
						{
							ID:        23,
							Title:     "test",
							Completed: true,
							Order:     42,
						},
					}

					storage := &MockStorage{}
					storage.InnerMock.On("GetAll").Return(todos, nil)

					return storage
				}(),
			},
			args: args{
				baseURL: &url.URL{Scheme: "https", Host: "example.com"},
			},
			want: []models.PresentationTodoRecord{
				{
					URL:       "https://example.com/api/v1/todos/5",
					Title:     "test",
					Completed: true,
					Order:     12,
				},
				{
					URL:       "https://example.com/api/v1/todos/23",
					Title:     "test",
					Completed: true,
					Order:     42,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				Storage: func() Storage {
					storage := &MockStorage{}
					storage.InnerMock.
						On("GetAll").
						Return([]models.TodoRecord(nil), iotest.ErrTimeout)

					return storage
				}(),
			},
			args: args{
				baseURL: &url.URL{Scheme: "https", Host: "example.com"},
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := TodoRecord{
				Storage: tt.fields.Storage,
			}
			got, err := useCase.GetAll(tt.args.baseURL)

			tt.fields.Storage.(*MockStorage).InnerMock.AssertExpectations(t)
			assert.Equal(t, tt.want, got)
			tt.wantErr(t, err)
		})
	}
}

func TestTodoRecord_GetSingle(t *testing.T) {
	type fields struct {
		Storage Storage
	}
	type args struct {
		baseURL *url.URL
		id      int
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
					storage.InnerMock.On("GetSingle", 23).Return(todo, nil)

					return storage
				}(),
			},
			args: args{
				baseURL: &url.URL{Scheme: "https", Host: "example.com"},
				id:      23,
			},
			want: models.PresentationTodoRecord{
				URL:       "https://example.com/api/v1/todos/23",
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
					storage := &MockStorage{}
					storage.InnerMock.
						On("GetSingle", 23).
						Return(models.TodoRecord{}, iotest.ErrTimeout)

					return storage
				}(),
			},
			args: args{
				baseURL: &url.URL{Scheme: "https", Host: "example.com"},
				id:      23,
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
			got, err := useCase.GetSingle(tt.args.baseURL, tt.args.id)

			tt.fields.Storage.(*MockStorage).InnerMock.AssertExpectations(t)
			assert.Equal(t, tt.want, got)
			tt.wantErr(t, err)
		})
	}
}

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

func TestTodoRecord_Update(t *testing.T) {
	type fields struct {
		Storage Storage
	}
	type args struct {
		baseURL *url.URL
		id      int
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
					storage.InnerMock.On("Update", 42, todo).Return(nil)

					return storage
				}(),
			},
			args: args{
				baseURL: &url.URL{Scheme: "https", Host: "example.com"},
				id:      42,
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
					storage.InnerMock.On("Update", 42, todo).Return(iotest.ErrTimeout)

					return storage
				}(),
			},
			args: args{
				baseURL: &url.URL{Scheme: "https", Host: "example.com"},
				id:      42,
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
			got, err := useCase.Update(tt.args.baseURL, tt.args.id, tt.args.todo)

			tt.fields.Storage.(*MockStorage).InnerMock.AssertExpectations(t)
			assert.Equal(t, tt.want, got)
			tt.wantErr(t, err)
		})
	}
}

func TestTodoRecord_Patch(t *testing.T) {
	type fields struct {
		Storage Storage
	}
	type args struct {
		baseURL   *url.URL
		id        int
		todoPatch models.TodoRecordPatch
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
					patchedTodo := models.TodoRecord{
						ID:        23,
						Title:     "test2",
						Completed: true,
						Order:     42,
					}

					storage := &MockStorage{}
					storage.InnerMock.On("GetSingle", 23).Return(todo, nil)
					storage.InnerMock.On("Update", 23, patchedTodo).Return(nil)

					return storage
				}(),
			},
			args: args{
				baseURL: &url.URL{Scheme: "https", Host: "example.com"},
				id:      23,
				todoPatch: models.TodoRecordPatch{
					Title: func() *string {
						title := "test2"
						return &title
					}(),
				},
			},
			want: models.PresentationTodoRecord{
				URL:       "https://example.com/api/v1/todos/23",
				Title:     "test2",
				Completed: true,
				Order:     42,
			},
			wantErr: assert.NoError,
		},
		{
			name: "error with getting",
			fields: fields{
				Storage: func() Storage {
					storage := &MockStorage{}
					storage.InnerMock.
						On("GetSingle", 23).
						Return(models.TodoRecord{}, iotest.ErrTimeout)

					return storage
				}(),
			},
			args: args{
				baseURL: &url.URL{Scheme: "https", Host: "example.com"},
				id:      23,
				todoPatch: models.TodoRecordPatch{
					Title: func() *string {
						title := "test2"
						return &title
					}(),
				},
			},
			want:    models.PresentationTodoRecord{},
			wantErr: assert.Error,
		},
		{
			name: "error with updating",
			fields: fields{
				Storage: func() Storage {
					todo := models.TodoRecord{
						ID:        23,
						Title:     "test",
						Completed: true,
						Order:     42,
					}
					patchedTodo := models.TodoRecord{
						ID:        23,
						Title:     "test2",
						Completed: true,
						Order:     42,
					}

					storage := &MockStorage{}
					storage.InnerMock.On("GetSingle", 23).Return(todo, nil)
					storage.InnerMock.On("Update", 23, patchedTodo).Return(iotest.ErrTimeout)

					return storage
				}(),
			},
			args: args{
				baseURL: &url.URL{Scheme: "https", Host: "example.com"},
				id:      23,
				todoPatch: models.TodoRecordPatch{
					Title: func() *string {
						title := "test2"
						return &title
					}(),
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
			got, err := useCase.Patch(tt.args.baseURL, tt.args.id, tt.args.todoPatch)

			tt.fields.Storage.(*MockStorage).InnerMock.AssertExpectations(t)
			assert.Equal(t, tt.want, got)
			tt.wantErr(t, err)
		})
	}
}

func TestTodoRecord_Delete(t *testing.T) {
	type fields struct {
		Storage Storage
	}
	type args struct {
		id int
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				Storage: func() Storage {
					storage := &MockStorage{}
					storage.InnerMock.On("Delete", 42).Return(nil)

					return storage
				}(),
			},
			args:    args{id: 42},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				Storage: func() Storage {
					storage := &MockStorage{}
					storage.InnerMock.On("Delete", 42).Return(iotest.ErrTimeout)

					return storage
				}(),
			},
			args:    args{id: 42},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := TodoRecord{
				Storage: tt.fields.Storage,
			}
			err := useCase.Delete(tt.args.id)

			tt.fields.Storage.(*MockStorage).InnerMock.AssertExpectations(t)
			tt.wantErr(t, err)
		})
	}
}
