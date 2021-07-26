package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	httputils "github.com/irenicaa/go-todo-backend/http-utils"
	"github.com/irenicaa/go-todo-backend/models"
	"github.com/stretchr/testify/assert"
)

func TestRouter_ServeHTTP(t *testing.T) {
	type fields struct {
		BaseURL   string
		URLScheme string
		UseCase   TodoRecordUseCase
		Logger    httputils.Logger
	}
	type args struct {
		request *http.Request
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		wantResponse *http.Response
	}{
		{
			name: "success with getting of all records",
			fields: fields{
				BaseURL:   "/api/v1",
				URLScheme: "http",
				UseCase: func() TodoRecordUseCase {
					baseURL := &url.URL{Scheme: "http", Host: "example.com"}
					presentationTodos := []models.PresentationTodoRecord{
						{
							URL: "http://example.com/api/v1/todos/5",
							Date: models.Date(time.Date(
								2006, time.January, 2,
								0, 0, 0, 0,
								time.UTC,
							)),
							Title:     "test",
							Completed: true,
							Order:     12,
						},
						{
							URL: "http://example.com/api/v1/todos/23",
							Date: models.Date(time.Date(
								2006, time.January, 3,
								0, 0, 0, 0,
								time.UTC,
							)),
							Title:     "test2",
							Completed: false,
							Order:     42,
						},
					}

					useCase := &MockTodoRecordUseCase{}
					useCase.InnerMock.
						On("GetAll", baseURL, models.Query{}).
						Return(presentationTodos, nil)

					return useCase
				}(),
				Logger: &MockLogger{},
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/api/v1/todos",
					nil,
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusOK) + " " +
					http.StatusText(http.StatusOK),
				StatusCode: http.StatusOK,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     http.Header{"Content-Type": {"application/json"}},
				Body: ioutil.NopCloser(bytes.NewReader([]byte(
					`[{"url":"http://example.com/api/v1/todos/5",` +
						`"date":"2006-01-02",` +
						`"title":"test",` +
						`"completed":true,` +
						`"order":12},` +
						`{"url":"http://example.com/api/v1/todos/23",` +
						`"date":"2006-01-03",` +
						`"title":"test2",` +
						`"completed":false,` +
						`"order":42}]`,
				))),
				ContentLength: -1,
			},
		},
		{
			name: "success with getting of all records by a date",
			fields: fields{
				BaseURL:   "/api/v1",
				URLScheme: "http",
				UseCase: func() TodoRecordUseCase {
					baseURL := &url.URL{Scheme: "http", Host: "example.com"}
					presentationTodos := []models.PresentationTodoRecord{
						{
							URL: "http://example.com/api/v1/todos/5",
							Date: models.Date(time.Date(
								2006, time.January, 2,
								0, 0, 0, 0,
								time.UTC,
							)),
							Title:     "test",
							Completed: true,
							Order:     12,
						},
						{
							URL: "http://example.com/api/v1/todos/23",
							Date: models.Date(time.Date(
								2006, time.January, 2,
								0, 0, 0, 0,
								time.UTC,
							)),
							Title:     "test2",
							Completed: false,
							Order:     42,
						},
					}

					useCase := &MockTodoRecordUseCase{}
					useCase.InnerMock.
						On("GetAll", baseURL, models.Query{
							MinimalDate: models.Date(time.Date(
								2006, time.January, 2,
								0, 0, 0, 0,
								time.UTC,
							)),
							MaximalDate: models.Date(time.Date(
								2006, time.January, 2,
								0, 0, 0, 0,
								time.UTC,
							)),
						}).
						Return(presentationTodos, nil)

					return useCase
				}(),
				Logger: &MockLogger{},
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/api/v1/todos/2006-01-02",
					nil,
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusOK) + " " +
					http.StatusText(http.StatusOK),
				StatusCode: http.StatusOK,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     http.Header{"Content-Type": {"application/json"}},
				Body: ioutil.NopCloser(bytes.NewReader([]byte(
					`[{"url":"http://example.com/api/v1/todos/5",` +
						`"date":"2006-01-02",` +
						`"title":"test",` +
						`"completed":true,` +
						`"order":12},` +
						`{"url":"http://example.com/api/v1/todos/23",` +
						`"date":"2006-01-02",` +
						`"title":"test2",` +
						`"completed":false,` +
						`"order":42}]`,
				))),
				ContentLength: -1,
			},
		},
		{
			name: "success with getting of a single record",
			fields: fields{
				BaseURL:   "/api/v1",
				URLScheme: "http",
				UseCase: func() TodoRecordUseCase {
					baseURL := &url.URL{Scheme: "http", Host: "example.com"}
					presentationTodo := models.PresentationTodoRecord{
						URL: "http://example.com/api/v1/todos/12",
						Date: models.Date(time.Date(
							2006, time.January, 2,
							0, 0, 0, 0,
							time.UTC,
						)),
						Title:     "test",
						Completed: true,
						Order:     23,
					}

					useCase := &MockTodoRecordUseCase{}
					useCase.InnerMock.
						On("GetSingle", baseURL, 12).
						Return(presentationTodo, nil)

					return useCase
				}(),
				Logger: &MockLogger{},
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/api/v1/todos/12",
					nil,
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusOK) + " " +
					http.StatusText(http.StatusOK),
				StatusCode: http.StatusOK,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     http.Header{"Content-Type": {"application/json"}},
				Body: ioutil.NopCloser(bytes.NewReader([]byte(
					`{"url":"http://example.com/api/v1/todos/12",` +
						`"date":"2006-01-02",` +
						`"title":"test",` +
						`"completed":true,` +
						`"order":23}`,
				))),
				ContentLength: -1,
			},
		},
		{
			name: "success with creating of a record",
			fields: fields{
				BaseURL:   "/api/v1",
				URLScheme: "http",
				UseCase: func() TodoRecordUseCase {
					baseURL := &url.URL{Scheme: "http", Host: "example.com"}
					presentationTodoIn := models.PresentationTodoRecord{
						Date: models.Date(time.Date(
							2006, time.January, 2,
							0, 0, 0, 0,
							time.UTC,
						)),
						Title:     "test",
						Completed: true,
						Order:     23,
					}
					presentationTodoOut := models.PresentationTodoRecord{
						URL: "http://example.com/api/v1/todos/12",
						Date: models.Date(time.Date(
							2006, time.January, 2,
							0, 0, 0, 0,
							time.UTC,
						)),
						Title:     "test",
						Completed: true,
						Order:     23,
					}

					useCase := &MockTodoRecordUseCase{}
					useCase.InnerMock.
						On("Create", baseURL, presentationTodoIn).
						Return(presentationTodoOut, nil)

					return useCase
				}(),
				Logger: &MockLogger{},
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodPost,
					"http://example.com/api/v1/todos/",
					bytes.NewReader([]byte(`{
						"date": "2006-01-02",
						"title": "test",
						"completed": true,
						"order": 23
					}`)),
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusOK) + " " +
					http.StatusText(http.StatusOK),
				StatusCode: http.StatusOK,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     http.Header{"Content-Type": {"application/json"}},
				Body: ioutil.NopCloser(bytes.NewReader([]byte(
					`{"url":"http://example.com/api/v1/todos/12",` +
						`"date":"2006-01-02",` +
						`"title":"test",` +
						`"completed":true,` +
						`"order":23}`,
				))),
				ContentLength: -1,
			},
		},
		{
			name: "success with updating of a record",
			fields: fields{
				BaseURL:   "/api/v1",
				URLScheme: "http",
				UseCase: func() TodoRecordUseCase {
					baseURL := &url.URL{Scheme: "http", Host: "example.com"}
					presentationTodoIn := models.PresentationTodoRecord{
						Date: models.Date(time.Date(
							2006, time.January, 2,
							0, 0, 0, 0,
							time.UTC,
						)),
						Title:     "test",
						Completed: true,
						Order:     23,
					}
					presentationTodoOut := models.PresentationTodoRecord{
						URL: "http://example.com/api/v1/todos/12",
						Date: models.Date(time.Date(
							2006, time.January, 2,
							0, 0, 0, 0,
							time.UTC,
						)),
						Title:     "test",
						Completed: true,
						Order:     23,
					}

					useCase := &MockTodoRecordUseCase{}
					useCase.InnerMock.
						On("Update", baseURL, 12, presentationTodoIn).
						Return(presentationTodoOut, nil)

					return useCase
				}(),
				Logger: &MockLogger{},
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodPut,
					"http://example.com/api/v1/todos/12",
					bytes.NewReader([]byte(`{
						"date": "2006-01-02",
						"title": "test",
						"completed": true,
						"order": 23
					}`)),
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusOK) + " " +
					http.StatusText(http.StatusOK),
				StatusCode: http.StatusOK,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     http.Header{"Content-Type": {"application/json"}},
				Body: ioutil.NopCloser(bytes.NewReader([]byte(
					`{"url":"http://example.com/api/v1/todos/12",` +
						`"date":"2006-01-02",` +
						`"title":"test",` +
						`"completed":true,` +
						`"order":23}`,
				))),
				ContentLength: -1,
			},
		},
		{
			name: "success with patching of a record",
			fields: fields{
				BaseURL:   "/api/v1",
				URLScheme: "http",
				UseCase: func() TodoRecordUseCase {
					baseURL := &url.URL{Scheme: "http", Host: "example.com"}
					todoPatchTitle := "test"
					todoPatch := models.TodoRecordPatch{Title: &todoPatchTitle}
					presentationTodo := models.PresentationTodoRecord{
						URL: "http://example.com/api/v1/todos/12",
						Date: models.Date(time.Date(
							2006, time.January, 2,
							0, 0, 0, 0,
							time.UTC,
						)),
						Title:     "test",
						Completed: true,
						Order:     23,
					}

					useCase := &MockTodoRecordUseCase{}
					useCase.InnerMock.
						On("Patch", baseURL, 12, todoPatch).
						Return(presentationTodo, nil)

					return useCase
				}(),
				Logger: &MockLogger{},
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodPatch,
					"http://example.com/api/v1/todos/12",
					bytes.NewReader([]byte(`{"title": "test"}`)),
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusOK) + " " +
					http.StatusText(http.StatusOK),
				StatusCode: http.StatusOK,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     http.Header{"Content-Type": {"application/json"}},
				Body: ioutil.NopCloser(bytes.NewReader([]byte(
					`{"url":"http://example.com/api/v1/todos/12",` +
						`"date":"2006-01-02",` +
						`"title":"test",` +
						`"completed":true,` +
						`"order":23}`,
				))),
				ContentLength: -1,
			},
		},
		{
			name: "success with deleting of all records",
			fields: fields{
				BaseURL:   "/api/v1",
				URLScheme: "http",
				UseCase: func() TodoRecordUseCase {
					useCase := &MockTodoRecordUseCase{}
					useCase.InnerMock.On("DeleteAll").Return(nil)

					return useCase
				}(),
				Logger: &MockLogger{},
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodDelete,
					"http://example.com/api/v1/todos",
					nil,
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusNoContent) + " " +
					http.StatusText(http.StatusNoContent),
				StatusCode:    http.StatusNoContent,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Header:        http.Header{},
				Body:          ioutil.NopCloser(bytes.NewReader(nil)),
				ContentLength: -1,
			},
		},
		{
			name: "success with deleting of a single record",
			fields: fields{
				BaseURL:   "/api/v1",
				URLScheme: "http",
				UseCase: func() TodoRecordUseCase {
					useCase := &MockTodoRecordUseCase{}
					useCase.InnerMock.On("DeleteSingle", 12).Return(nil)

					return useCase
				}(),
				Logger: &MockLogger{},
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodDelete,
					"http://example.com/api/v1/todos/12",
					nil,
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusNoContent) + " " +
					http.StatusText(http.StatusNoContent),
				StatusCode:    http.StatusNoContent,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Header:        http.Header{},
				Body:          ioutil.NopCloser(bytes.NewReader(nil)),
				ContentLength: -1,
			},
		},
		{
			name: "error",
			fields: fields{
				BaseURL:   "/api/v1",
				URLScheme: "http",
				UseCase:   &MockTodoRecordUseCase{},
				Logger: func() httputils.Logger {
					logger := &MockLogger{}
					logger.InnerMock.
						On("Print", []interface{}{http.StatusText(http.StatusNotFound)}).
						Return().
						Times(1)

					return logger
				}(),
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/api/v1/incorrect",
					nil,
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusNotFound) + " " +
					http.StatusText(http.StatusNotFound),
				StatusCode: http.StatusNotFound,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     http.Header{},
				Body: ioutil.NopCloser(bytes.NewReader([]byte(
					http.StatusText(http.StatusNotFound),
				))),
				ContentLength: -1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()
			router := Router{
				BaseURL: tt.fields.BaseURL,
				TodoRecord: TodoRecord{
					URLScheme: tt.fields.URLScheme,
					UseCase:   tt.fields.UseCase,
					Logger:    tt.fields.Logger,
				},
				Logger: tt.fields.Logger,
			}
			router.ServeHTTP(responseRecorder, tt.args.request)

			tt.fields.UseCase.(*MockTodoRecordUseCase).InnerMock.AssertExpectations(t)
			tt.fields.Logger.(*MockLogger).InnerMock.AssertExpectations(t)
			assert.Equal(t, tt.wantResponse, responseRecorder.Result())
		})
	}
}
