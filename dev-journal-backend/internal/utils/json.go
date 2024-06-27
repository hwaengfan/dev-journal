package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// @return json payload from request if there is any
func ParseJSON(request *http.Request, payload any) error {
	if request.Body == nil {
		return fmt.Errorf("Request body is empty")
	}

	return json.NewDecoder(request.Body).Decode(payload)
}

// @return json response with status code
func WriteJSON(writer http.ResponseWriter, status int, content any) error {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(status)

	return json.NewEncoder(writer).Encode(content)
}

// @return json error response with status code
func WriteError(writer http.ResponseWriter, status int, err error) {
	WriteJSON(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
}
