package usecases

import (
	"github.com/irenicaa/go-todo-backend/models"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	InnerMock mock.Mock
}

func (mock *MockStorage) Create(todo models.TodoRecord) (id int, err error) {
	results := mock.InnerMock.Called(todo)
	return results.Int(0), results.Error(1)
}
