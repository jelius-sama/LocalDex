package api

import (
	"LocalDex/logger"
	"LocalDex/parser"
	"LocalDex/types"
	"LocalDex/util"
	"bytes"
	"fmt"
	"net/http"
	"os"
)

func ServePages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		MethodNotAllowed(w, r, getOnlyRoute)
		return
	}

	html, err := parser.GetHTML()
	if err != nil {
		InternalErrorPage(w, r, util.AddrOf("Something went wrong when getting html from FS!"))
		logger.Error("failed to get html shell:\n    " + err.Error())
		return
	}

	ssrData, err, status := parser.PerformSSR(r.URL.Path)
	if err != nil {
		if status == http.StatusNotFound {
			NotFoundPage(w, r, util.AddrOf("Content of dynamic page could not be found!"))
			return
		}

		// TODO: Implement dedicated error pages intead of generic 500 error.
		InternalErrorPage(w, r, util.AddrOf("Failed to perform SSR!"))
		logger.Error("performing SSR failed:\n    " + err.Error())
		return
	}

	if len(ssrData) == 0 {
		var metadata string

		defaultCase := func() error {
			metadata, err = parser.ParseMetadata(r.URL.Path, nil)
			if err != nil {
				return fmt.Errorf("metadata parsing failed: %w", err)
			}
			return nil
		}

		switch r.URL.Path {
		default:
			if err := defaultCase(); err != nil {
				InternalErrorPage(w, r, util.AddrOf("Failed to parse metadata of the page!"))
				logger.Error(err.Error())
				return
			}
		}

		// Replace marker in HTML
		html = bytes.Replace(html, []byte("<!-- Server Props -->"), []byte(metadata), 1)
	} else {
		html = bytes.Replace(html, []byte("<!-- SSR Data -->"), []byte(ssrData), 1)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if len(ssrData) == 0 {
		w.Header().Set("Cache-Control", "public, max-age=86400")
	} else {
		w.Header().Set("Cache-Control", "public, max-age=3600")
	}
	if os.Getenv("ENV") == types.ENV.Prod {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	}
	w.WriteHeader(http.StatusOK)
	w.Write(html)
}
