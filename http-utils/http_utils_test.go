package httputils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"testing/iotest"

	"github.com/irenicaa/go-todo-backend/models"
	"github.com/stretchr/testify/assert"
)

func TestGetRequestBody(t *testing.T) {
	type args struct {
		request     *http.Request
		requestData interface{}
	}

	tests := []struct {
		name            string
		args            args
		wantRequestData interface{}
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				request: httptest.NewRequest(
					http.MethodPost,
					"http://example.com/",
					bytes.NewReader([]byte(`{
						"Title": "test",
						"Completed": true,
						"Order": 23
					}`)),
				),
				requestData: &models.TodoRecord{},
			},
			wantRequestData: &models.TodoRecord{
				Title:     "test",
				Completed: true,
				Order:     23,
			},
			wantErr: assert.NoError,
		},
		{
			name: "error on reading",
			args: args{
				request: httptest.NewRequest(
					http.MethodPost,
					"http://example.com/",
					iotest.TimeoutReader(bytes.NewReader([]byte(`{
						"Title": "test",
						"Completed": true,
						"Order": 23
					}`))),
				),
				requestData: &models.TodoRecord{},
			},
			wantRequestData: &models.TodoRecord{},
			wantErr:         assert.Error,
		},
		{
			name: "error on unmarshalling",
			args: args{
				request: httptest.NewRequest(
					http.MethodPost,
					"http://example.com/",
					bytes.NewReader([]byte("incorrect")),
				),
				requestData: &models.TodoRecord{},
			},
			wantRequestData: &models.TodoRecord{},
			wantErr:         assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GetRequestBody(tt.args.request, tt.args.requestData)

			assert.Equal(t, tt.wantRequestData, tt.args.requestData)
			tt.wantErr(t, err)
		})
	}
}

func TestHandleError(t *testing.T) {
	type args struct {
		logger    Logger
		status    int
		format    string
		arguments []interface{}
	}

	tests := []struct {
		name         string
		args         args
		wantResponse *http.Response
	}{
		{
			name: "succes",
			args: args{
				logger: func() Logger {
					logger := &MockLogger{}
					logger.InnerMock.
						On("Print", []interface{}{"test: 23 one"}).
						Return().
						Times(1)

					return logger
				}(),
				status:    http.StatusNotFound,
				format:    "test: %d %s",
				arguments: []interface{}{23, "one"},
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusNotFound) + " " +
					http.StatusText(http.StatusNotFound),
				StatusCode:    http.StatusNotFound,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Header:        http.Header{},
				Body:          ioutil.NopCloser(bytes.NewReader([]byte("test: 23 one"))),
				ContentLength: -1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()
			HandleError(
				responseRecorder,
				tt.args.logger,
				tt.args.status,
				tt.args.format,
				tt.args.arguments...,
			)

			tt.args.logger.(*MockLogger).InnerMock.AssertExpectations(t)
			assert.Equal(t, tt.wantResponse, responseRecorder.Result())
		})
	}
}

func TestHandleJSON(t *testing.T) {
	type testData struct {
		FieldOne int
		FieldTwo int
	}
	type incorrectTestData struct {
		FieldOne   int
		FieldTwo   int
		FieldThree func()
	}
	type args struct {
		logger Logger
		data   interface{}
	}

	tests := []struct {
		name         string
		args         args
		wantResponse *http.Response
	}{
		{
			name: "success",
			args: args{
				logger: &MockLogger{},
				data:   testData{FieldOne: 23, FieldTwo: 42},
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusOK) + " " +
					http.StatusText(http.StatusOK),
				StatusCode: http.StatusOK,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     http.Header{"Content-Type": {"application/json"}},
				Body: ioutil.NopCloser(bytes.NewReader(
					[]byte(`{"FieldOne":23,"FieldTwo":42}`),
				)),
				ContentLength: -1,
			},
		},
		{
			name: "error",
			args: args{
				logger: func() Logger {
					logger := &MockLogger{}
					logger.InnerMock.
						On("Print", []interface{}{
							"unable to marshal the data: json: unsupported type: func()",
						}).
						Return().
						Times(1)

					return logger
				}(),
				data: incorrectTestData{
					FieldOne:   23,
					FieldTwo:   42,
					FieldThree: func() {},
				},
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusInternalServerError) + " " +
					http.StatusText(http.StatusInternalServerError),
				StatusCode: http.StatusInternalServerError,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     http.Header{},
				Body: ioutil.NopCloser(bytes.NewReader(
					[]byte(`unable to marshal the data: json: unsupported type: func()`),
				)),
				ContentLength: -1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()
			HandleJSON(responseRecorder, tt.args.logger, tt.args.data)

			tt.args.logger.(*MockLogger).InnerMock.AssertExpectations(t)
			assert.Equal(t, tt.wantResponse, responseRecorder.Result())
		})
	}
}
