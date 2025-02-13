package utils

import (
	"fmt"
	"net/http"
)

func ErrorJsonResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	if statusCode == http.StatusOK {
		w.WriteHeader(statusCode)
	}

	w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, message)))
}
