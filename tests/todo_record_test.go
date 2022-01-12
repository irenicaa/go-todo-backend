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
	"net/url"
	"strconv"
	"testing"
	"time"

	httputils "github.com/irenicaa/go-todo-backend/http-utils"
	"github.com/irenicaa/go-todo-backend/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var port = flag.Int("port", 8080, "server port")

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

func TestTodoRecord_withQuerying(t *testing.T) {
	tests := []struct {
		name            string
		originalTodos   []models.PresentationTodoRecord
		queryParameters url.Values
		wantTodos       []models.PresentationTodoRecord
	}{
		{
			name: "with filtration by the minimal date",
			originalTodos: func() []models.PresentationTodoRecord {
				var originalTodos []models.PresentationTodoRecord
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
					originalTodos = append(originalTodos, originalTodo)
				}

				return originalTodos
			}(),
			queryParameters: url.Values{"minimal_date": {"2006-01-07"}},
			wantTodos: func() []models.PresentationTodoRecord {
				var wantTodos []models.PresentationTodoRecord
				for i := 10; i >= 5; i-- {
					wantTodo := models.PresentationTodoRecord{
						Date: models.Date(time.Date(
							2006, time.January, 2+i,
							0, 0, 0, 0,
							time.UTC,
						)),
						Title:     "test" + strconv.Itoa(i),
						Completed: true,
						Order:     i,
					}
					wantTodos = append(wantTodos, wantTodo)
				}

				return wantTodos
			}(),
		},
		{
			name: "with filtration by the maximal date",
			originalTodos: func() []models.PresentationTodoRecord {
				var originalTodos []models.PresentationTodoRecord
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
					originalTodos = append(originalTodos, originalTodo)
				}

				return originalTodos
			}(),
			queryParameters: url.Values{"maximal_date": {"2006-01-07"}},
			wantTodos: func() []models.PresentationTodoRecord {
				var wantTodos []models.PresentationTodoRecord
				for i := 5; i >= 0; i-- {
					wantTodo := models.PresentationTodoRecord{
						Date: models.Date(time.Date(
							2006, time.January, 2+i,
							0, 0, 0, 0,
							time.UTC,
						)),
						Title:     "test" + strconv.Itoa(i),
						Completed: true,
						Order:     i,
					}
					wantTodos = append(wantTodos, wantTodo)
				}

				return wantTodos
			}(),
		},
		{
			name: "with filtration by the date range",
			originalTodos: func() []models.PresentationTodoRecord {
				var originalTodos []models.PresentationTodoRecord
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
					originalTodos = append(originalTodos, originalTodo)
				}

				return originalTodos
			}(),
			queryParameters: url.Values{
				"minimal_date": {"2006-01-05"},
				"maximal_date": {"2006-01-09"},
			},
			wantTodos: func() []models.PresentationTodoRecord {
				var wantTodos []models.PresentationTodoRecord
				for i := 7; i >= 3; i-- {
					wantTodo := models.PresentationTodoRecord{
						Date: models.Date(time.Date(
							2006, time.January, 2+i,
							0, 0, 0, 0,
							time.UTC,
						)),
						Title:     "test" + strconv.Itoa(i),
						Completed: true,
						Order:     i,
					}
					wantTodos = append(wantTodos, wantTodo)
				}

				return wantTodos
			}(),
		},
		{
			name: "with search by the title fragment",
			originalTodos: func() []models.PresentationTodoRecord {
				var originalTodos []models.PresentationTodoRecord
				for i := 0; i <= 10; i++ {
					var mark string
					if i%2 == 0 {
						mark = "even"
					} else {
						mark = "odd"
					}

					originalTodo := models.PresentationTodoRecord{
						Date: models.Date(time.Date(
							2006, time.January, 2+i,
							0, 0, 0, 0,
							time.UTC,
						)),
						Title:     fmt.Sprintf("test%d (%s)", i, mark),
						Completed: true,
						Order:     i,
					}
					originalTodos = append(originalTodos, originalTodo)
				}

				return originalTodos
			}(),
			queryParameters: url.Values{"title_fragment": {"even"}},
			wantTodos: func() []models.PresentationTodoRecord {
				var wantTodos []models.PresentationTodoRecord
				for i := 10; i >= 0; i -= 2 {
					wantTodo := models.PresentationTodoRecord{
						Date: models.Date(time.Date(
							2006, time.January, 2+i,
							0, 0, 0, 0,
							time.UTC,
						)),
						Title:     fmt.Sprintf("test%d (even)", i),
						Completed: true,
						Order:     i,
					}
					wantTodos = append(wantTodos, wantTodo)
				}

				return wantTodos
			}(),
		},
		{
			name: "with pagination",
			originalTodos: func() []models.PresentationTodoRecord {
				var originalTodos []models.PresentationTodoRecord
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
					originalTodos = append(originalTodos, originalTodo)
				}

				return originalTodos
			}(),
			queryParameters: url.Values{"page_size": {"2"}, "page": {"3"}},
			wantTodos: func() []models.PresentationTodoRecord {
				var wantTodos []models.PresentationTodoRecord
				for i := 6; i >= 5; i-- {
					wantTodo := models.PresentationTodoRecord{
						Date: models.Date(time.Date(
							2006, time.January, 2+i,
							0, 0, 0, 0,
							time.UTC,
						)),
						Title:     "test" + strconv.Itoa(i),
						Completed: true,
						Order:     i,
					}
					wantTodos = append(wantTodos, wantTodo)
				}

				return wantTodos
			}(),
		},
		{
			name: "with all parameters",
			originalTodos: func() []models.PresentationTodoRecord {
				var originalTodos []models.PresentationTodoRecord
				for i := 0; i <= 20; i++ {
					var mark string
					if i%2 == 0 {
						mark = "even"
					} else {
						mark = "odd"
					}

					originalTodo := models.PresentationTodoRecord{
						Date: models.Date(time.Date(
							2006, time.January, 2+i,
							0, 0, 0, 0,
							time.UTC,
						)),
						Title:     fmt.Sprintf("test%d (%s)", i, mark),
						Completed: true,
						Order:     i,
					}
					originalTodos = append(originalTodos, originalTodo)
				}

				return originalTodos
			}(),
			queryParameters: url.Values{
				"minimal_date":   {"2006-01-05"},
				"maximal_date":   {"2006-01-19"},
				"title_fragment": {"even"},
				"page_size":      {"2"},
				"page":           {"3"},
			},
			wantTodos: func() []models.PresentationTodoRecord {
				var wantTodos []models.PresentationTodoRecord
				for i := 8; i >= 6; i -= 2 {
					wantTodo := models.PresentationTodoRecord{
						Date: models.Date(time.Date(
							2006, time.January, 2+i,
							0, 0, 0, 0,
							time.UTC,
						)),
						Title:     fmt.Sprintf("test%d (even)", i),
						Completed: true,
						Order:     i,
					}
					wantTodos = append(wantTodos, wantTodo)
				}

				return wantTodos
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("http://localhost:%d/api/v1/todos", *port)
			_, err := sendRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			for _, originalTodo := range tt.originalTodos {
				_, err2 := sendRequest(http.MethodPost, url, originalTodo)
				require.NoError(t, err2)
			}

			queryURL := url + "?" + tt.queryParameters.Encode()
			response, err := sendRequest(http.MethodGet, queryURL, nil)
			require.NoError(t, err)
			defer response.Body.Close()

			var gotTodos []models.PresentationTodoRecord
			err = httputils.GetJSONData(response.Body, &gotTodos)
			require.NoError(t, err)
			for index := range gotTodos {
				gotTodos[index].URL = ""
			}

			assert.Equal(t, tt.wantTodos, gotTodos)
		})
	}
}

func TestTodoRecord_withModifying(t *testing.T) {
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
