package handlers

import (
	"net/url"

	"github.com/irenicaa/go-todo-backend/models"
	"github.com/stretchr/testify/mock"
)

type MockTodoRecordUseCase struct {
	InnerMock mock.Mock
}

func (mock *MockTodoRecordUseCase) GetAll(
	baseURL *url.URL,
) ([]models.PresentationTodoRecord, error) {
	results := mock.InnerMock.Called(baseURL)
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
	todo models.TodoRecord,
) (models.PresentationTodoRecord, error) {
	results := mock.InnerMock.Called(baseURL, todo)
	return results.Get(0).(models.PresentationTodoRecord), results.Error(1)
}

func (mock *MockTodoRecordUseCase) Update(
	baseURL *url.URL,
	id int,
	todo models.TodoRecord,
) (models.PresentationTodoRecord, error) {
	results := mock.InnerMock.Called(baseURL, id, todo)
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
