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
					`{"URL":"http://example.com/api/v1/todos/12",` +
						`"Title":"test",` +
						`"Completed":true,` +
						`"Order":23}`,
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
						"unable to unmarshal the request body: " +
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
						"unable to unmarshal the request body: " +
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
