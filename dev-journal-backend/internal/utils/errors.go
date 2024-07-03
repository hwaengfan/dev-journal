package utils

import (
	"fmt"
	"net/http"
)

// WriteError writes an error to the response
func WriteError(writer http.ResponseWriter, status int, err error) {
	WriteJSON(writer, status, map[string]string{"error": err.Error()})
}

// WriteInvalidPayload writes an invalid payload error to the response
func WriteInvalidPayload(writer http.ResponseWriter, errors error) {
	WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
}

// WritePermissionDenied writes a permission denied error to the response
func WritePermissionDenied(writer http.ResponseWriter) {
	WriteError(writer, http.StatusForbidden, fmt.Errorf("permission denied"))
}
