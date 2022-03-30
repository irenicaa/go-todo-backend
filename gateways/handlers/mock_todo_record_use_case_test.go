package handlers

import (
	"net/url"

	"github.com/irenicaa/go-todo-backend/v2/models"
	"github.com/stretchr/testify/mock"
)

type MockTodoRecordUseCase struct {
	InnerMock mock.Mock
}

func (mock *MockTodoRecordUseCase) GetAll(
	baseURL *url.URL,
	query models.Query,
) ([]models.PresentationTodoRecord, error) {
	results := mock.InnerMock.Called(baseURL, query)
	return results.Get(0).([]models.PresentationTodoRecord), results.Error(1)
}

func (mock *MockTodoRecordUseCase) GetSingle(
	baseURL *url.URL,
	id int,
) (models.PresentationTodoRecord, error) {
	results := mock.InnerMock.Called(baseURL, id)
	return results.Get(0).(models.PresentationTodoRecord), results.Error(1)
}

func (mock *MockTodoRecordUseCase) Create(
	baseURL *url.URL,
	presentationTodo models.PresentationTodoRecord,
) (models.PresentationTodoRecord, error) {
	results := mock.InnerMock.Called(baseURL, presentationTodo)
	return results.Get(0).(models.PresentationTodoRecord), results.Error(1)
}

func (mock *MockTodoRecordUseCase) Update(
	baseURL *url.URL,
	id int,
	presentationTodo models.PresentationTodoRecord,
) (models.PresentationTodoRecord, error) {
	results := mock.InnerMock.Called(baseURL, id, presentationTodo)
	return results.Get(0).(models.PresentationTodoRecord), results.Error(1)
}

func (mock *MockTodoRecordUseCase) Patch(
	baseURL *url.URL,
	id int,
	todoPatch models.TodoRecordPatch,
) (models.PresentationTodoRecord, error) {
	results := mock.InnerMock.Called(baseURL, id, todoPatch)
	return results.Get(0).(models.PresentationTodoRecord), results.Error(1)
}

func (mock *MockTodoRecordUseCase) DeleteAll() error {
	results := mock.InnerMock.Called()
	return results.Error(0)
}

func (mock *MockTodoRecordUseCase) DeleteSingle(id int) error {
	results := mock.InnerMock.Called(id)
	return results.Error(0)
}
