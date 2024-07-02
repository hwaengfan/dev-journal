package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// For validating payload
var Validate = validator.New()

// ParseJSON parses a JSON request
func ParseJSON(request *http.Request, payload any) error {
	if request.Body == nil {
		return fmt.Errorf("request body is empty")
	}

	return json.NewDecoder(request.Body).Decode(payload)
}

// WriteJSON writes a JSON response
func WriteJSON(writer http.ResponseWriter, status int, content any) error {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(status)

	return json.NewEncoder(writer).Encode(content)
}
