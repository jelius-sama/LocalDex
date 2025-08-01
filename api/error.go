package api

import (
	"LocalDex/logger"
	"LocalDex/parser"
	"LocalDex/types"
	"LocalDex/util"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func BadRequest(w http.ResponseWriter, r *http.Request, msg *string) {
	util.WriteError(w, http.StatusBadRequest, "400 Bad Request", msg)
}

func Unauthorized(w http.ResponseWriter, r *http.Request, msg *string) {
	util.WriteError(w, http.StatusUnauthorized, "401 Unauthorized", msg)
}

func Forbidden(w http.ResponseWriter, r *http.Request, msg *string) {
	util.WriteError(w, http.StatusForbidden, "403 Forbidden", msg)
}

func NotFoundAPI(w http.ResponseWriter, r *http.Request, msg *string) {
	util.WriteError(w, http.StatusNotFound, "404 Not Found", msg)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request, msg *string) {
	util.WriteError(w, http.StatusMethodNotAllowed, "405 Method Not Allowed", msg)
}

func RequestTimeout(w http.ResponseWriter, r *http.Request, msg *string) {
	util.WriteError(w, http.StatusRequestTimeout, "408 Request Timeout", msg)
}

func InternalErrorAPI(w http.ResponseWriter, r *http.Request, msg *string) {
	util.WriteError(w, http.StatusInternalServerError, "500 Internal Server Error", msg)
}

func InternalErrorPage(w http.ResponseWriter, r *http.Request, msg *string) {
	// Attempt to get the HTML shell
	html, err := parser.GetHTML()
	if err != nil {
		// If even the shell fails, fall back to basic response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		logger.Error("failed to get html shell in error handler:\n    " + err.Error())
		return
	}

	// Create the SSR data as JSON
	payload := map[string]any{
		"status":  500,
		"message": "Internal Server Error",
	}
	if msg != nil {
		payload["message"] = *msg
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		logger.Error("failed to marshal SSR error JSON:\n    " + err.Error())
		return
	}

	serverProps, err := parser.ParseStaticMetadataForPaths([]string{"*", "#internal_server_error"})
	if err != nil {
		http.Error(w, "Something went wrong!", http.StatusInternalServerError)
		logger.Error("failed to parse metadata for server props:\n    " + err.Error())
		return
	}

	// Build the <script> tag with ID __SERVER_DATA__
	ssrScript := fmt.Sprintf(`<script id="__SERVER_DATA__" type="application/json">%s</script>`, jsonData)

	// Inject into the HTML shell at the SSR Data marker
	ssrInjectedHTML := bytes.Replace(html, []byte("<!-- SSR Data -->"), []byte(ssrScript), 1)
	finalHTML := bytes.Replace(ssrInjectedHTML, []byte("<!-- Server Props -->"), []byte(serverProps), 1)

	// Write the response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if os.Getenv("ENV") == types.ENV.Prod {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(finalHTML)
}

func NotFoundPage(w http.ResponseWriter, r *http.Request, msg *string) {
	// Attempt to get the HTML shell
	html, err := parser.GetHTML()
	if err != nil {
		InternalErrorPage(w, r, util.AddrOf("Something went wrong!"))
		logger.Error("failed to get html shell in error handler:\n    " + err.Error())
		return
	}

	// Create the SSR data as JSON
	payload := map[string]any{
		"status":  404,
		"message": "Page not found",
		"path":    r.URL.Path,
	}
	if msg != nil {
		payload["message"] = *msg
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		InternalErrorPage(w, r, util.AddrOf("Internal Server Error"))
		logger.Error("failed to marshal SSR error JSON:\n    " + err.Error())
		return
	}

	serverProps, err := parser.ParseStaticMetadataForPaths([]string{"*", "#not_found"})
	if err != nil {
		InternalErrorPage(w, r, util.AddrOf("Something went wrong!"))
		logger.Error("failed to parse metadata for server props:\n    " + err.Error())
		return
	}

	// Build the <script> tag with ID __SERVER_DATA__
	ssrScript := fmt.Sprintf(`<script id="__SERVER_DATA__" type="application/json">%s</script>`, jsonData)

	// Inject into the HTML shell at the SSR Data marker
	ssrInjectedHTML := bytes.Replace(html, []byte("<!-- SSR Data -->"), []byte(ssrScript), 1)
	finalHTML := bytes.Replace(ssrInjectedHTML, []byte("<!-- Server Props -->"), []byte(serverProps), 1)

	// Write the response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if os.Getenv("ENV") == types.ENV.Prod {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.WriteHeader(http.StatusNotFound)
	w.Write(finalHTML)
}
