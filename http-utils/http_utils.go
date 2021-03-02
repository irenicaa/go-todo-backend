package httputils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Logger ...
type Logger interface {
	Print(arguments ...interface{})
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
