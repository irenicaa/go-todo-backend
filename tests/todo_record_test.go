package tests

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/irenicaa/go-todo-backend/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var port = flag.Int("port", 8080, "server port")

func TestTodoRecord_withSingleModel(t *testing.T) {
	tests := []struct {
		name         string
		originalTodo models.TodoRecord
		action       func(t *testing.T, todoURL string)
		wantTodo     models.PresentationTodoRecord
	}{
		{
			name: "creation",
			originalTodo: models.TodoRecord{
				Title:     "test",
				Completed: true,
				Order:     42,
			},
			action: func(t *testing.T, todoURL string) {},
			wantTodo: models.PresentationTodoRecord{
				Title:     "test",
				Completed: true,
				Order:     42,
			},
		},
		{
			name: "updating",
			originalTodo: models.TodoRecord{
				Title:     "test",
				Completed: true,
				Order:     23,
			},
			action: func(t *testing.T, todoURL string) {
				newTodo := models.TodoRecord{
					Title:     "test2",
					Completed: true,
					Order:     42,
				}

				requestBytes, err := json.Marshal(newTodo)
				require.NoError(t, err)

				request, err :=
					http.NewRequest(http.MethodPut, todoURL, bytes.NewReader(requestBytes))
				require.NoError(t, err)

				_, err = http.DefaultClient.Do(request)
				require.NoError(t, err)
			},
			wantTodo: models.PresentationTodoRecord{
				Title:     "test2",
				Completed: true,
				Order:     42,
			},
		},
		{
			name: "patching",
			originalTodo: models.TodoRecord{
				Title:     "test",
				Completed: true,
				Order:     42,
			},
			action: func(t *testing.T, todoURL string) {
				todoPatchTitle := "test2"
				todoPatch := models.TodoRecordPatch{Title: &todoPatchTitle}

				requestBytes, err := json.Marshal(todoPatch)
				require.NoError(t, err)

				request, err :=
					http.NewRequest(http.MethodPatch, todoURL, bytes.NewReader(requestBytes))
				require.NoError(t, err)

				_, err = http.DefaultClient.Do(request)
				require.NoError(t, err)
			},
			wantTodo: models.PresentationTodoRecord{
				Title:     "test2",
				Completed: true,
				Order:     42,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBytes, err := json.Marshal(tt.originalTodo)
			require.NoError(t, err)

			url := fmt.Sprintf("http://localhost:%d/api/v1/todos", *port)
			response, err :=
				http.Post(url, "application/json", bytes.NewReader(requestBytes))
			require.NoError(t, err)

			createdTodo, err := unmarshalTodoRecord(response.Body)
			require.NoError(t, err)

			tt.action(t, createdTodo.URL)

			response, err = http.Get(createdTodo.URL)
			require.NoError(t, err)

			gotTodo, err := unmarshalTodoRecord(response.Body)
			require.NoError(t, err)

			tt.wantTodo.URL = createdTodo.URL
			assert.Equal(t, tt.wantTodo, gotTodo)
		})
	}
}

func unmarshalTodoRecord(reader io.ReadCloser) (
	models.PresentationTodoRecord,
	error,
) {
	defer reader.Close()

	responseBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return models.PresentationTodoRecord{}, err
	}

	var todo models.PresentationTodoRecord
	if err := json.Unmarshal(responseBytes, &todo); err != nil {
		return models.PresentationTodoRecord{}, err
	}

	return todo, nil
}