package middleware

import (
	"LocalDex/logger"
	"net/http"
)

func TestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.TimedInfo("TestMiddleware: before handler")
		next.ServeHTTP(w, r)
	})
}

func TestMiddlewareSecond(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.TimedInfo("TestMiddlewareSecond: before handler")
		next.ServeHTTP(w, r)
	})
}
