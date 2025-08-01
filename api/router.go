package api

import (
	vars "LocalDex"
	"LocalDex/logger"
	"LocalDex/types"
	"LocalDex/util"
	"net/http"
	"os"
	"runtime"
	"strings"

	"io/fs"
	"mime"
	"path/filepath"
)

var getOnlyRoute = util.AddrOf("Only Request with GET Method are allowed!")

func HandleRouting() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/", ServePages)

	router.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/") {
			_, file, line, ok := runtime.Caller(1)
			NotFoundAPI(w, r, util.AddrOf("Unreachable code reached!"))
			if ok {
				logger.Error("Reached unreachable code in `", file, "` at line: ", line)
			} else {
				logger.Error("Reached unreachable code (unknown location)")
			}
			return
		}

		method := r.Method
		path := strings.TrimPrefix(r.URL.Path, "/api/")
		lookupKey := method + " /" + path

		if handler, exists := ApiRoutes[lookupKey]; exists {
			NoCache(handler)(w, r)
			return
		}

		NotFoundAPI(w, r, util.AddrOf("API Route Not Found!"))
	})

	router.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			MethodNotAllowed(w, r, getOnlyRoute)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/")
		content, err := fs.ReadFile(vars.AssetsFS, path)
		if err != nil {
			NotFoundAPI(w, r, util.AddrOf("Requested Asset was not found"))
			return
		}

		ext := filepath.Ext(path)
		mimeType := mime.TypeByExtension(ext)

		switch path {
		case "assets/sw.js":
			mimeType = "application/javascript"
			break

		case "assets/manifest.json":
			mimeType = "application/manifest+json"
			w.Header().Set("Service-Worker-Allowed", "/")
			break

		default:
			if len(mimeType) == 0 {
				mimeType = "application/octet-stream"
			}
			break
		}

		w.Header().Set("Content-Type", mimeType)
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.Write(content)
	})

	router.HandleFunc("/src/", func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("ENV") == types.ENV.Dev {
			NotFoundAPI(w, r, util.AddrOf("/src/ route hit in development mode which is not permitted"))
			return
		}
		if r.Method != http.MethodGet {
			MethodNotAllowed(w, r, getOnlyRoute)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/")
		content, err := fs.ReadFile(vars.ViteFS, "client/dist/"+path)
		if err != nil {
			NotFoundAPI(w, r, util.AddrOf("Requested Source file was not found!"))
			return
		}

		ext := filepath.Ext(path)
		mimeType := mime.TypeByExtension(ext)
		if len(mimeType) == 0 {
			mimeType = "application/octet-stream"
		}

		w.Header().Set("Content-Type", mimeType)
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.Write(content)
	})

	return router
}
