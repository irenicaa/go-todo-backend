package usecases

import (
	"github.com/irenicaa/go-todo-backend/models"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	InnerMock mock.Mock
}

func (mock *MockStorage) GetSingle(id int) (models.TodoRecord, error) {
	results := mock.InnerMock.Called(id)
	return results.Get(0).(models.TodoRecord), results.Error(1)
}

func (mock *MockStorage) Create(todo models.TodoRecord) (id int, err error) {
	results := mock.InnerMock.Called(todo)
	return results.Int(0), results.Error(1)
}

func (mock *MockStorage) Update(id int, todo models.TodoRecord) error {
	results := mock.InnerMock.Called(id, todo)
	return results.Error(0)
}

func (mock *MockStorage) Delete(id int) error {
	results := mock.InnerMock.Called(id)
	return results.Error(0)
}
