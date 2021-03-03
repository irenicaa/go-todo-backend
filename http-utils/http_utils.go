package httputils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Logger ...
type Logger interface {
	Print(arguments ...interface{})
}

// GetRequestBody ...
func GetRequestBody(request *http.Request, requestData interface{}) error {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return fmt.Errorf("unable to read the request body: %s", err)
	}

	if err := json.Unmarshal(body, requestData); err != nil {
		return fmt.Errorf("unable to unmarshal the request body: %s", err)
	}

	return nil
}

// HandleError ...
func HandleError(
	writer http.ResponseWriter,
	logger Logger,
	status int,
	format string,
	arguments ...interface{},
) {
	message := fmt.Sprintf(format, arguments...)
	logger.Print(message)

	writer.WriteHeader(status)
	writer.Write([]byte(message))
}

// HandleJSON ...
func HandleJSON(writer http.ResponseWriter, logger Logger, data interface{}) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		HandleError(
			writer,
			logger,
			http.StatusInternalServerError,
			"unable to marshal the data: %v",
			err,
		)

		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(dataBytes)
}
