package util

import (
	"LocalDex/types"
	"encoding/json"
	"net/http"
	"strconv"
)

// AddrOf takes a value and returns its address as a pointer.
func AddrOf[T any](literal T) *T { return &literal }

func IsValidPort(portStr string) bool {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return false // not a number
	}
	return port >= 1 && port <= 65535
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, statusCode int, errorText string, msg *string) {
	resp := map[string]any{
		"error":   errorText,
		"Message": msg,
	}

	WriteJSON(w, statusCode, resp)
}

func WriteSuccess(w http.ResponseWriter, statusCode int, successText string, msg *string) {
	resp := map[string]any{
		"success": successText,
		"Message": msg,
	}

	WriteJSON(w, statusCode, resp)
}

func Chain(cm types.MiddlewareChain) http.Handler {
	for i := len(cm.Middlewares) - 1; i >= 0; i-- {
		cm.Handler = cm.Middlewares[i](cm.Handler)
	}
	return cm.Handler
}
