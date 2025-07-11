package api

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string  `json:"error"`
	Message *string `json:"message"`
}

func writeError(w http.ResponseWriter, statusCode int, errorText string, msg *string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := ErrorResponse{
		Error:   errorText,
		Message: msg,
	}

	json.NewEncoder(w).Encode(resp)
}

func BadRequest(w http.ResponseWriter, r *http.Request, msg *string) {
	writeError(w, http.StatusBadRequest, "400 Bad Request", msg)
}

func Unauthorized(w http.ResponseWriter, r *http.Request, msg *string) {
	writeError(w, http.StatusUnauthorized, "401 Unauthorized", msg)
}

func Forbidden(w http.ResponseWriter, r *http.Request, msg *string) {
	writeError(w, http.StatusForbidden, "403 Forbidden", msg)
}

func NotFound(w http.ResponseWriter, r *http.Request, msg *string) {
	writeError(w, http.StatusNotFound, "404 Not Found", msg)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request, msg *string) {
	writeError(w, http.StatusMethodNotAllowed, "405 Method Not Allowed", msg)
}

func RequestTimeout(w http.ResponseWriter, r *http.Request, msg *string) {
	writeError(w, http.StatusRequestTimeout, "408 Request Timeout", msg)
}
