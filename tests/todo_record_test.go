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
	"time"

	httputils "github.com/irenicaa/go-todo-backend/http-utils"
	"github.com/irenicaa/go-todo-backend/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var port = flag.Int("port", 8080, "server port")

func TestTodoRecord_withSingleModel(t *testing.T) {
	tests := []struct {
		name         string
		originalTodo models.PresentationTodoRecord
		action       func(t *testing.T, todoURL string)
		wantTodo     models.PresentationTodoRecord
	}{
		{
			name: "creation",
			originalTodo: models.PresentationTodoRecord{
				Date: models.Date(time.Date(
					2006, time.January, 2,
					0, 0, 0, 0,
					time.UTC,
				)),
				Title:     "test",
				Completed: true,
				Order:     42,
			},
			action: func(t *testing.T, todoURL string) {},
			wantTodo: models.PresentationTodoRecord{
				Date: models.Date(time.Date(
					2006, time.January, 2,
					0, 0, 0, 0,
					time.UTC,
				)),
				Title:     "test",
				Completed: true,
				Order:     42,
			},
		},
		{
			name: "updating",
			originalTodo: models.PresentationTodoRecord{
				Date: models.Date(time.Date(
					2006, time.January, 2,
					0, 0, 0, 0,
					time.UTC,
				)),
				Title:     "test",
				Completed: true,
				Order:     23,
			},
			action: func(t *testing.T, todoURL string) {
				newTodo := models.PresentationTodoRecord{
					Date: models.Date(time.Date(
						2006, time.January, 3,
						0, 0, 0, 0,
						time.UTC,
					)),
					Title:     "test2",
					Completed: true,
					Order:     42,
				}

				_, err := sendRequest(http.MethodPut, todoURL, newTodo)
				require.NoError(t, err)
			},
			wantTodo: models.PresentationTodoRecord{
				Date: models.Date(time.Date(
					2006, time.January, 3,
					0, 0, 0, 0,
					time.UTC,
				)),
				Title:     "test2",
				Completed: true,
				Order:     42,
			},
		},
		{
			name: "patching",
			originalTodo: models.PresentationTodoRecord{
				Date: models.Date(time.Date(
					2006, time.January, 2,
					0, 0, 0, 0,
					time.UTC,
				)),
				Title:     "test",
				Completed: true,
				Order:     42,
			},
			action: func(t *testing.T, todoURL string) {
				todoPatchTitle := "test2"
				todoPatch := models.TodoRecordPatch{Title: &todoPatchTitle}

				_, err := sendRequest(http.MethodPatch, todoURL, todoPatch)
				require.NoError(t, err)
			},
			wantTodo: models.PresentationTodoRecord{
				Date: models.Date(time.Date(
					2006, time.January, 2,
					0, 0, 0, 0,
					time.UTC,
				)),
				Title:     "test2",
				Completed: true,
				Order:     42,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("http://localhost:%d/api/v1/todos", *port)
			response, err := sendRequest(http.MethodPost, url, tt.originalTodo)
			require.NoError(t, err)

			createdTodo, err := unmarshalTodoRecord(response.Body)
			require.NoError(t, err)

			tt.action(t, createdTodo.URL)

			response, err = sendRequest(http.MethodGet, createdTodo.URL, nil)
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
		originalTodo := models.PresentationTodoRecord{
			Date: models.Date(time.Date(
				2006, time.January, 2+i,
				0, 0, 0, 0,
				time.UTC,
			)),
			Title:     "test" + strconv.Itoa(i),
			Completed: true,
			Order:     i,
		}

		response, err := sendRequest(http.MethodPost, url, originalTodo)
		require.NoError(t, err)

		createdTodo, err := unmarshalTodoRecord(response.Body)
		require.NoError(t, err)

		createdTodos = append(createdTodos, createdTodo)
	}
	// reversing
	for i, j := 0, len(createdTodos)-1; i < j; i, j = i+1, j-1 {
		createdTodos[i], createdTodos[j] = createdTodos[j], createdTodos[i]
	}

	response, err := sendRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	defer response.Body.Close()

	var gotTodos []models.PresentationTodoRecord
	err = httputils.GetJSONData(response.Body, &gotTodos)
	require.NoError(t, err)

	assert.Equal(t, createdTodos, gotTodos)
}

func TestTodoRecord_withDeleting(t *testing.T) {
	originalTodo := models.PresentationTodoRecord{
		Date: models.Date(time.Date(
			2006, time.January, 2,
			0, 0, 0, 0,
			time.UTC,
		)),
		Title:     "test",
		Completed: true,
		Order:     42,
	}

	url := fmt.Sprintf("http://localhost:%d/api/v1/todos", *port)
	response, err := sendRequest(http.MethodPost, url, originalTodo)
	require.NoError(t, err)

	createdTodo, err := unmarshalTodoRecord(response.Body)
	require.NoError(t, err)

	_, err = sendRequest(http.MethodDelete, createdTodo.URL, nil)
	require.NoError(t, err)

	response, err = sendRequest(http.MethodGet, createdTodo.URL, nil)
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

func sendRequest(method string, url string, data interface{}) (
	*http.Response,
	error,
) {
	var body io.Reader
	if data != nil {
		dataAsJSON, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		body = bytes.NewReader(dataAsJSON)
	}

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
