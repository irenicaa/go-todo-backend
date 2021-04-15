// +build integration

package tests

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	httputils "github.com/irenicaa/go-todo-backend/http-utils"
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

				_, err = sendRequest(http.MethodPut, todoURL, bytes.NewReader(requestBytes))
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

				_, err = sendRequest(http.MethodPatch, todoURL, bytes.NewReader(requestBytes))
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

func TestTodoRecord_withGetting(t *testing.T) {
	url := fmt.Sprintf("http://localhost:%d/api/v1/todos", *port)
	_, err := sendRequest(http.MethodDelete, url, nil)
	require.NoError(t, err)

	var createdTodos []models.PresentationTodoRecord
	for i := 0; i <= 10; i++ {
		originalTodo := models.TodoRecord{
			Title:     "test" + strconv.Itoa(i),
			Completed: true,
			Order:     i,
		}

		requestBytes, err := json.Marshal(originalTodo)
		require.NoError(t, err)

		response, err :=
			http.Post(url, "application/json", bytes.NewReader(requestBytes))
		require.NoError(t, err)

		createdTodo, err := unmarshalTodoRecord(response.Body)
		require.NoError(t, err)

		createdTodos = append(createdTodos, createdTodo)
	}

	response, err := http.Get(url)
	require.NoError(t, err)
	defer response.Body.Close()

	var gotTodos []models.PresentationTodoRecord
	err = httputils.GetJSONData(response.Body, &gotTodos)
	require.NoError(t, err)

	assert.Equal(t, createdTodos, gotTodos)
}

func TestTodoRecord_withDeleting(t *testing.T) {
	originalTodo := models.TodoRecord{
		Title:     "test",
		Completed: true,
		Order:     42,
	}

	requestBytes, err := json.Marshal(originalTodo)
	require.NoError(t, err)

	url := fmt.Sprintf("http://localhost:%d/api/v1/todos", *port)
	response, err :=
		http.Post(url, "application/json", bytes.NewReader(requestBytes))
	require.NoError(t, err)

	createdTodo, err := unmarshalTodoRecord(response.Body)
	require.NoError(t, err)

	_, err = sendRequest(http.MethodDelete, createdTodo.URL, nil)
	require.NoError(t, err)

	response, err = http.Get(createdTodo.URL)
	require.NoError(t, err)
	defer response.Body.Close()

	responseBytes, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	assert.Equal(
		t,
		"unable to get the to-do record: sql: no rows in result set",
		string(responseBytes),
	)
}

func sendRequest(method string, url string, body io.Reader) (
	*http.Response,
	error,
) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func unmarshalTodoRecord(reader io.ReadCloser) (
	models.PresentationTodoRecord,
	error,
) {
	defer reader.Close()

	var todo models.PresentationTodoRecord
	if err := httputils.GetJSONData(reader, &todo); err != nil {
		return models.PresentationTodoRecord{}, err
	}

	return todo, nil
}
