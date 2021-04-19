package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"testing/iotest"

	httputils "github.com/irenicaa/go-todo-backend/http-utils"
	"github.com/irenicaa/go-todo-backend/models"
	"github.com/stretchr/testify/assert"
)

func TestTodoRecord_GetAll(t *testing.T) {
	type fields struct {
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
			name: "success with to-do records",
			fields: fields{
				URLScheme: "http",
				UseCase: func() TodoRecordUseCase {
					baseURL := &url.URL{Scheme: "http", Host: "example.com"}
					presentationTodos := []models.PresentationTodoRecord{
						{
							URL:       "http://example.com/api/v1/todos/5",
							Title:     "test",
							Completed: true,
							Order:     12,
						},
						{
							URL:       "http://example.com/api/v1/todos/23",
							Title:     "test",
							Completed: true,
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
						`"title":"test",` +
						`"completed":true,` +
						`"order":12},` +
						`{"url":"http://example.com/api/v1/todos/23",` +
						`"title":"test",` +
						`"completed":true,` +
						`"order":42}]`,
				))),
				ContentLength: -1,
			},
		},
		{
			name: "success without to-do records",
			fields: fields{
				URLScheme: "http",
				UseCase: func() TodoRecordUseCase {
					baseURL := &url.URL{Scheme: "http", Host: "example.com"}
					presentationTodos := []models.PresentationTodoRecord(nil)

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
				StatusCode:    http.StatusOK,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Header:        http.Header{"Content-Type": {"application/json"}},
				Body:          ioutil.NopCloser(bytes.NewReader([]byte(`[]`))),
				ContentLength: -1,
			},
		},
		{
			name: "success with the title fragment",
			fields: fields{
				URLScheme: "http",
				UseCase: func() TodoRecordUseCase {
					baseURL := &url.URL{Scheme: "http", Host: "example.com"}
					presentationTodos := []models.PresentationTodoRecord{
						{
							URL:       "http://example.com/api/v1/todos/5",
							Title:     "test",
							Completed: true,
							Order:     12,
						},
						{
							URL:       "http://example.com/api/v1/todos/23",
							Title:     "test",
							Completed: true,
							Order:     42,
						},
					}

					useCase := &MockTodoRecordUseCase{}
					useCase.InnerMock.
						On("GetAll", baseURL, models.Query{TitleFragment: "test"}).
						Return(presentationTodos, nil)

					return useCase
				}(),
				Logger: &MockLogger{},
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/api/v1/todos?title_fragment=test",
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
						`"title":"test",` +
						`"completed":true,` +
						`"order":12},` +
						`{"url":"http://example.com/api/v1/todos/23",` +
						`"title":"test",` +
						`"completed":true,` +
						`"order":42}]`,
				))),
				ContentLength: -1,
			},
		},
		{
			name: "error",
			fields: fields{
				URLScheme: "http",
				UseCase: func() TodoRecordUseCase {
					baseURL := &url.URL{Scheme: "http", Host: "example.com"}

					useCase := &MockTodoRecordUseCase{}
					useCase.InnerMock.
						On("GetAll", baseURL, models.Query{}).
						Return([]models.PresentationTodoRecord(nil), iotest.ErrTimeout)

					return useCase
				}(),
				Logger: func() httputils.Logger {
					logger := &MockLogger{}
					logger.InnerMock.
						On("Print", []interface{}{"timeout"}).
						Return().
						Times(1)

					return logger
				}(),
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/api/v1/todos",
					nil,
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusInternalServerError) + " " +
					http.StatusText(http.StatusInternalServerError),
				StatusCode:    http.StatusInternalServerError,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Header:        http.Header{},
				Body:          ioutil.NopCloser(bytes.NewReader([]byte("timeout"))),
				ContentLength: -1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()
			handler := TodoRecord{
				URLScheme: tt.fields.URLScheme,
				UseCase:   tt.fields.UseCase,
				Logger:    tt.fields.Logger,
			}
			handler.GetAll(responseRecorder, tt.args.request)

			tt.fields.UseCase.(*MockTodoRecordUseCase).InnerMock.AssertExpectations(t)
			tt.fields.Logger.(*MockLogger).InnerMock.AssertExpectations(t)
			assert.Equal(t, tt.wantResponse, responseRecorder.Result())
		})
	}
}

func TestTodoRecord_GetSingle(t *testing.T) {
	type fields struct {
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
			name: "success",
			fields: fields{
				URLScheme: "http",
				UseCase: func() TodoRecordUseCase {
					baseURL := &url.URL{Scheme: "http", Host: "example.com"}
					presentationTodo := models.PresentationTodoRecord{
						URL:       "http://example.com/api/v1/todos/12",
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
						`"title":"test",` +
						`"completed":true,` +
						`"order":23}`,
				))),
				ContentLength: -1,
			},
		},
		{
			name: "error on ID getting",
			fields: fields{
				URLScheme: "http",
				UseCase:   &MockTodoRecordUseCase{},
				Logger: func() httputils.Logger {
					message := "unable to get an ID: " +
						"unable to find an ID"
					logger := &MockLogger{}
					logger.InnerMock.
						On("Print", []interface{}{message}).
						Return().
						Times(1)

					return logger
				}(),
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/api/v1/todos/",
					nil,
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusBadRequest) + " " +
					http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     http.Header{},
				Body: ioutil.NopCloser(bytes.NewReader([]byte(
					"unable to get an ID: " +
						"unable to find an ID",
				))),
				ContentLength: -1,
			},
		},
		{
			name: "error on to-do record getting",
			fields: fields{
				URLScheme: "http",
				UseCase: func() TodoRecordUseCase {
					baseURL := &url.URL{Scheme: "http", Host: "example.com"}

					useCase := &MockTodoRecordUseCase{}
					useCase.InnerMock.
						On("GetSingle", baseURL, 12).
						Return(models.PresentationTodoRecord{}, iotest.ErrTimeout)

					return useCase
				}(),
				Logger: func() httputils.Logger {
					logger := &MockLogger{}
					logger.InnerMock.
						On("Print", []interface{}{"timeout"}).
						Return().
						Times(1)

					return logger
				}(),
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/api/v1/todos/12",
					nil,
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusInternalServerError) + " " +
					http.StatusText(http.StatusInternalServerError),
				StatusCode:    http.StatusInternalServerError,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Header:        http.Header{},
				Body:          ioutil.NopCloser(bytes.NewReader([]byte("timeout"))),
				ContentLength: -1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()
			handler := TodoRecord{
				URLScheme: tt.fields.URLScheme,
				UseCase:   tt.fields.UseCase,
				Logger:    tt.fields.Logger,
			}
			handler.GetSingle(responseRecorder, tt.args.request)

			tt.fields.UseCase.(*MockTodoRecordUseCase).InnerMock.AssertExpectations(t)
			tt.fields.Logger.(*MockLogger).InnerMock.AssertExpectations(t)
			assert.Equal(t, tt.wantResponse, responseRecorder.Result())
		})
	}
}

func TestTodoRecord_Create(t *testing.T) {
	type fields struct {
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
			name: "success",
			fields: fields{
				URLScheme: "http",
				UseCase: func() TodoRecordUseCase {
					baseURL := &url.URL{Scheme: "http", Host: "example.com"}
					todo := models.TodoRecord{Title: "test", Completed: true, Order: 23}
					presentationTodo := models.PresentationTodoRecord{
						URL:       "http://example.com/api/v1/todos/12",
						Title:     "test",
						Completed: true,
						Order:     23,
					}

					useCase := &MockTodoRecordUseCase{}
					useCase.InnerMock.On("Create", baseURL, todo).Return(presentationTodo, nil)

					return useCase
				}(),
				Logger: &MockLogger{},
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodPost,
					"http://example.com/api/v1/todos/",
					bytes.NewReader([]byte(`{
						"Title": "test",
						"Completed": true,
						"Order": 23
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
						`"title":"test",` +
						`"completed":true,` +
						`"order":23}`,
				))),
				ContentLength: -1,
			},
		},
		{
			name: "error on request body getting",
			fields: fields{
				URLScheme: "http",
				UseCase:   &MockTodoRecordUseCase{},
				Logger: func() httputils.Logger {
					message := "unable to get the request body: " +
						"unable to unmarshal the JSON data: " +
						"invalid character 'i' looking for beginning of value"
					logger := &MockLogger{}
					logger.InnerMock.
						On("Print", []interface{}{message}).
						Return().
						Times(1)

					return logger
				}(),
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodPost,
					"http://example.com/api/v1/todos/",
					bytes.NewReader([]byte("incorrect")),
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusBadRequest) + " " +
					http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     http.Header{},
				Body: ioutil.NopCloser(bytes.NewReader([]byte(
					"unable to get the request body: " +
						"unable to unmarshal the JSON data: " +
						"invalid character 'i' looking for beginning of value",
				))),
				ContentLength: -1,
			},
		},
		{
			name: "error on to-do record creating",
			fields: fields{
				URLScheme: "http",
				UseCase: func() TodoRecordUseCase {
					baseURL := &url.URL{Scheme: "http", Host: "example.com"}
					todo := models.TodoRecord{Title: "test", Completed: true, Order: 23}

					useCase := &MockTodoRecordUseCase{}
					useCase.InnerMock.
						On("Create", baseURL, todo).
						Return(models.PresentationTodoRecord{}, iotest.ErrTimeout)

					return useCase
				}(),
				Logger: func() httputils.Logger {
					logger := &MockLogger{}
					logger.InnerMock.
						On("Print", []interface{}{"timeout"}).
						Return().
						Times(1)

					return logger
				}(),
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodPost,
					"http://example.com/api/v1/todos/",
					bytes.NewReader([]byte(`{
						"Title": "test",
						"Completed": true,
						"Order": 23
					}`)),
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusInternalServerError) + " " +
					http.StatusText(http.StatusInternalServerError),
				StatusCode:    http.StatusInternalServerError,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Header:        http.Header{},
				Body:          ioutil.NopCloser(bytes.NewReader([]byte("timeout"))),
				ContentLength: -1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()
			handler := TodoRecord{
				URLScheme: tt.fields.URLScheme,
				UseCase:   tt.fields.UseCase,
				Logger:    tt.fields.Logger,
			}
			handler.Create(responseRecorder, tt.args.request)

			tt.fields.UseCase.(*MockTodoRecordUseCase).InnerMock.AssertExpectations(t)
			tt.fields.Logger.(*MockLogger).InnerMock.AssertExpectations(t)
			assert.Equal(t, tt.wantResponse, responseRecorder.Result())
		})
	}
}

func TestTodoRecord_Update(t *testing.T) {
	type fields struct {
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
			name: "success",
			fields: fields{
				URLScheme: "http",
				UseCase: func() TodoRecordUseCase {
					baseURL := &url.URL{Scheme: "http", Host: "example.com"}
					todo := models.TodoRecord{Title: "test", Completed: true, Order: 23}
					presentationTodo := models.PresentationTodoRecord{
						URL:       "http://example.com/api/v1/todos/12",
						Title:     "test",
						Completed: true,
						Order:     23,
					}

					useCase := &MockTodoRecordUseCase{}
					useCase.InnerMock.
						On("Update", baseURL, 12, todo).
						Return(presentationTodo, nil)

					return useCase
				}(),
				Logger: &MockLogger{},
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodPut,
					"http://example.com/api/v1/todos/12",
					bytes.NewReader([]byte(`{
						"Title": "test",
						"Completed": true,
						"Order": 23
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
						`"title":"test",` +
						`"completed":true,` +
						`"order":23}`,
				))),
				ContentLength: -1,
			},
		},
		{
			name: "error on ID getting",
			fields: fields{
				URLScheme: "http",
				UseCase:   &MockTodoRecordUseCase{},
				Logger: func() httputils.Logger {
					message := "unable to get an ID: " +
						"unable to find an ID"
					logger := &MockLogger{}
					logger.InnerMock.
						On("Print", []interface{}{message}).
						Return().
						Times(1)

					return logger
				}(),
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodPut,
					"http://example.com/api/v1/todos/",
					bytes.NewReader([]byte(`{
						"Title": "test",
						"Completed": true,
						"Order": 23
					}`)),
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusBadRequest) + " " +
					http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     http.Header{},
				Body: ioutil.NopCloser(bytes.NewReader([]byte(
					"unable to get an ID: " +
						"unable to find an ID",
				))),
				ContentLength: -1,
			},
		},
		{
			name: "error on request body getting",
			fields: fields{
				URLScheme: "http",
				UseCase:   &MockTodoRecordUseCase{},
				Logger: func() httputils.Logger {
					message := "unable to get the request body: " +
						"unable to unmarshal the JSON data: " +
						"invalid character 'i' looking for beginning of value"
					logger := &MockLogger{}
					logger.InnerMock.
						On("Print", []interface{}{message}).
						Return().
						Times(1)

					return logger
				}(),
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodPut,
					"http://example.com/api/v1/todos/12",
					bytes.NewReader([]byte("incorrect")),
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusBadRequest) + " " +
					http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     http.Header{},
				Body: ioutil.NopCloser(bytes.NewReader([]byte(
					"unable to get the request body: " +
						"unable to unmarshal the JSON data: " +
						"invalid character 'i' looking for beginning of value",
				))),
				ContentLength: -1,
			},
		},
		{
			name: "error on to-do record updating",
			fields: fields{
				URLScheme: "http",
				UseCase: func() TodoRecordUseCase {
					baseURL := &url.URL{Scheme: "http", Host: "example.com"}
					todo := models.TodoRecord{Title: "test", Completed: true, Order: 23}

					useCase := &MockTodoRecordUseCase{}
					useCase.InnerMock.
						On("Update", baseURL, 12, todo).
						Return(models.PresentationTodoRecord{}, iotest.ErrTimeout)

					return useCase
				}(),
				Logger: func() httputils.Logger {
					logger := &MockLogger{}
					logger.InnerMock.
						On("Print", []interface{}{"timeout"}).
						Return().
						Times(1)

					return logger
				}(),
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodPut,
					"http://example.com/api/v1/todos/12",
					bytes.NewReader([]byte(`{
						"Title": "test",
						"Completed": true,
						"Order": 23
					}`)),
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusInternalServerError) + " " +
					http.StatusText(http.StatusInternalServerError),
				StatusCode:    http.StatusInternalServerError,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Header:        http.Header{},
				Body:          ioutil.NopCloser(bytes.NewReader([]byte("timeout"))),
				ContentLength: -1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()
			handler := TodoRecord{
				URLScheme: tt.fields.URLScheme,
				UseCase:   tt.fields.UseCase,
				Logger:    tt.fields.Logger,
			}
			handler.Update(responseRecorder, tt.args.request)

			tt.fields.UseCase.(*MockTodoRecordUseCase).InnerMock.AssertExpectations(t)
			tt.fields.Logger.(*MockLogger).InnerMock.AssertExpectations(t)
			assert.Equal(t, tt.wantResponse, responseRecorder.Result())
		})
	}
}

func TestTodoRecord_Patch(t *testing.T) {
	type fields struct {
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
			name: "success",
			fields: fields{
				URLScheme: "http",
				UseCase: func() TodoRecordUseCase {
					baseURL := &url.URL{Scheme: "http", Host: "example.com"}
					todoPatchTitle := "test"
					todoPatch := models.TodoRecordPatch{Title: &todoPatchTitle}
					presentationTodo := models.PresentationTodoRecord{
						URL:       "http://example.com/api/v1/todos/12",
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
					bytes.NewReader([]byte(`{"Title": "test"}`)),
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
						`"title":"test",` +
						`"completed":true,` +
						`"order":23}`,
				))),
				ContentLength: -1,
			},
		},
		{
			name: "error on ID getting",
			fields: fields{
				URLScheme: "http",
				UseCase:   &MockTodoRecordUseCase{},
				Logger: func() httputils.Logger {
					message := "unable to get an ID: " +
						"unable to find an ID"
					logger := &MockLogger{}
					logger.InnerMock.
						On("Print", []interface{}{message}).
						Return().
						Times(1)

					return logger
				}(),
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodPatch,
					"http://example.com/api/v1/todos/",
					bytes.NewReader([]byte(`{"Title": "test"}`)),
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusBadRequest) + " " +
					http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     http.Header{},
				Body: ioutil.NopCloser(bytes.NewReader([]byte(
					"unable to get an ID: " +
						"unable to find an ID",
				))),
				ContentLength: -1,
			},
		},
		{
			name: "error on request body getting",
			fields: fields{
				URLScheme: "http",
				UseCase:   &MockTodoRecordUseCase{},
				Logger: func() httputils.Logger {
					message := "unable to get the request body: " +
						"unable to unmarshal the JSON data: " +
						"invalid character 'i' looking for beginning of value"
					logger := &MockLogger{}
					logger.InnerMock.
						On("Print", []interface{}{message}).
						Return().
						Times(1)

					return logger
				}(),
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodPatch,
					"http://example.com/api/v1/todos/12",
					bytes.NewReader([]byte("incorrect")),
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusBadRequest) + " " +
					http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     http.Header{},
				Body: ioutil.NopCloser(bytes.NewReader([]byte(
					"unable to get the request body: " +
						"unable to unmarshal the JSON data: " +
						"invalid character 'i' looking for beginning of value",
				))),
				ContentLength: -1,
			},
		},
		{
			name: "error on to-do record patching",
			fields: fields{
				URLScheme: "http",
				UseCase: func() TodoRecordUseCase {
					baseURL := &url.URL{Scheme: "http", Host: "example.com"}
					todoPatchTitle := "test"
					todoPatch := models.TodoRecordPatch{Title: &todoPatchTitle}

					useCase := &MockTodoRecordUseCase{}
					useCase.InnerMock.
						On("Patch", baseURL, 12, todoPatch).
						Return(models.PresentationTodoRecord{}, iotest.ErrTimeout)

					return useCase
				}(),
				Logger: func() httputils.Logger {
					logger := &MockLogger{}
					logger.InnerMock.
						On("Print", []interface{}{"timeout"}).
						Return().
						Times(1)

					return logger
				}(),
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodPatch,
					"http://example.com/api/v1/todos/12",
					bytes.NewReader([]byte(`{"Title": "test"}`)),
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusInternalServerError) + " " +
					http.StatusText(http.StatusInternalServerError),
				StatusCode:    http.StatusInternalServerError,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Header:        http.Header{},
				Body:          ioutil.NopCloser(bytes.NewReader([]byte("timeout"))),
				ContentLength: -1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()
			handler := TodoRecord{
				URLScheme: tt.fields.URLScheme,
				UseCase:   tt.fields.UseCase,
				Logger:    tt.fields.Logger,
			}
			handler.Patch(responseRecorder, tt.args.request)

			tt.fields.UseCase.(*MockTodoRecordUseCase).InnerMock.AssertExpectations(t)
			tt.fields.Logger.(*MockLogger).InnerMock.AssertExpectations(t)
			assert.Equal(t, tt.wantResponse, responseRecorder.Result())
		})
	}
}

func TestTodoRecord_DeleteAll(t *testing.T) {
	type fields struct {
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
			name: "success",
			fields: fields{
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
			name: "error",
			fields: fields{
				URLScheme: "http",
				UseCase: func() TodoRecordUseCase {
					useCase := &MockTodoRecordUseCase{}
					useCase.InnerMock.On("DeleteAll").Return(iotest.ErrTimeout)

					return useCase
				}(),
				Logger: func() httputils.Logger {
					logger := &MockLogger{}
					logger.InnerMock.
						On("Print", []interface{}{"timeout"}).
						Return().
						Times(1)

					return logger
				}(),
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodDelete,
					"http://example.com/api/v1/todos",
					nil,
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusInternalServerError) + " " +
					http.StatusText(http.StatusInternalServerError),
				StatusCode:    http.StatusInternalServerError,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Header:        http.Header{},
				Body:          ioutil.NopCloser(bytes.NewReader([]byte("timeout"))),
				ContentLength: -1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()
			handler := TodoRecord{
				URLScheme: tt.fields.URLScheme,
				UseCase:   tt.fields.UseCase,
				Logger:    tt.fields.Logger,
			}
			handler.DeleteAll(responseRecorder, tt.args.request)

			tt.fields.UseCase.(*MockTodoRecordUseCase).InnerMock.AssertExpectations(t)
			tt.fields.Logger.(*MockLogger).InnerMock.AssertExpectations(t)
			assert.Equal(t, tt.wantResponse, responseRecorder.Result())
		})
	}
}

func TestTodoRecord_DeleteSingle(t *testing.T) {
	type fields struct {
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
			name: "success",
			fields: fields{
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
			name: "error on ID getting",
			fields: fields{
				URLScheme: "http",
				UseCase:   &MockTodoRecordUseCase{},
				Logger: func() httputils.Logger {
					message := "unable to get an ID: " +
						"unable to find an ID"
					logger := &MockLogger{}
					logger.InnerMock.
						On("Print", []interface{}{message}).
						Return().
						Times(1)

					return logger
				}(),
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodDelete,
					"http://example.com/api/v1/todos/",
					nil,
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusBadRequest) + " " +
					http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     http.Header{},
				Body: ioutil.NopCloser(bytes.NewReader([]byte(
					"unable to get an ID: " +
						"unable to find an ID",
				))),
				ContentLength: -1,
			},
		},
		{
			name: "error on to-do record deleting",
			fields: fields{
				URLScheme: "http",
				UseCase: func() TodoRecordUseCase {
					useCase := &MockTodoRecordUseCase{}
					useCase.InnerMock.On("DeleteSingle", 12).Return(iotest.ErrTimeout)

					return useCase
				}(),
				Logger: func() httputils.Logger {
					logger := &MockLogger{}
					logger.InnerMock.
						On("Print", []interface{}{"timeout"}).
						Return().
						Times(1)

					return logger
				}(),
			},
			args: args{
				request: httptest.NewRequest(
					http.MethodDelete,
					"http://example.com/api/v1/todos/12",
					nil,
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusInternalServerError) + " " +
					http.StatusText(http.StatusInternalServerError),
				StatusCode:    http.StatusInternalServerError,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Header:        http.Header{},
				Body:          ioutil.NopCloser(bytes.NewReader([]byte("timeout"))),
				ContentLength: -1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()
			handler := TodoRecord{
				URLScheme: tt.fields.URLScheme,
				UseCase:   tt.fields.UseCase,
				Logger:    tt.fields.Logger,
			}
			handler.DeleteSingle(responseRecorder, tt.args.request)

			tt.fields.UseCase.(*MockTodoRecordUseCase).InnerMock.AssertExpectations(t)
			tt.fields.Logger.(*MockLogger).InnerMock.AssertExpectations(t)
			assert.Equal(t, tt.wantResponse, responseRecorder.Result())
		})
	}
}
