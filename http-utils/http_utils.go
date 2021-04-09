package httputils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

// Logger ...
type Logger interface {
	Print(arguments ...interface{})
}

var idPattern = regexp.MustCompile(`/\d+`)

// GetIDFromURL ...
func GetIDFromURL(request *http.Request) (int, error) {
	idAsStr := idPattern.FindString(request.URL.Path)
	if idAsStr == "" {
		return 0, errors.New("unable to find an ID")
	}

	id, err := strconv.Atoi(idAsStr[1:])
	if err != nil {
		return 0, fmt.Errorf("unable to parse the ID: %s", err)
	}

	return id, nil
}

// GetJSONData ...
func GetJSONData(reader io.Reader, data interface{}) error {
	dataAsJSON, err := ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("unable to read the JSON data: %s", err)
	}

	if err := json.Unmarshal(dataAsJSON, data); err != nil {
		return fmt.Errorf("unable to unmarshal the JSON data: %s", err)
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
