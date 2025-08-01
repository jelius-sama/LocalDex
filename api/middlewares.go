package api

import (
	"LocalDex/logger"
	"fmt"
	"net/http"
	"strings"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.TimedInfo(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.TimedError("Encountered a panic, returning 500 to client and recovering the server!")
				if strings.HasPrefix(r.URL.Path, "/api/") {
					InternalErrorAPI(w, r, nil)
					return
				}

				if strings.HasPrefix(r.URL.Path, "/assets/") {
					InternalErrorAPI(w, r, nil)
					return
				}

				if strings.HasPrefix(r.URL.Path, "/src/") {
					InternalErrorAPI(w, r, nil)
					return
				}

				InternalErrorPage(w, r, nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func NoCache(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache") // for HTTP/1.0 proxies
		w.Header().Set("Expires", "0")       // for older caches
		next(w, r)
	}
}

func Cacheable(maxAge int, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
		next(w, r)
	}
}
