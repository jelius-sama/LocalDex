package main

import (
	vars "LocalDex"
	"LocalDex/api"
	"LocalDex/db"
	"LocalDex/logger"
	"LocalDex/types"
	"LocalDex/util"
	"errors"
	"fmt"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	DevPort     = "6969"
	Environment = "development"

	Title       = "LocalDex"
	Description = "NAS Server"
	Version     string
	AppRoot     string
	ETCDir      string
	Host        = "http://localhost"

	ReverseProxy string
	ProxyPort    string
)

// INFO: This function should only be called after `init()` or during init after the `APP_ROOT` variable has been set
func createAppRoot() {
	path := os.Getenv("APP_ROOT")

	err := dirExists(path)
	if err != nil {
		if strings.Contains(err.Error(), "path exists but is not a directory") {
			// Fatal error: file exists with the same name as directory
			logger.Panic("Cannot proceed:", err)
		}

		if strings.Contains(err.Error(), "directory does not exist") {
			// Try to create the directory
			if mkErr := os.MkdirAll(path, 0755); mkErr != nil {
				logger.Panic("Failed to create directory", path, ":", mkErr)
			}
		} else {
			// Other unexpected errors (permission, etc.)
			logger.Panic("Error checking directory", path, ":", err)
		}
	}
}

// NOTE: Be sure to explicitly set the app root in the `../config/server.config.json` file.
//   - This is important for the server store files and data.
//   - Conversely, if the server does not run with enough privileges, it will be unable
//   - access the TLS Certificates.
func init() {
	if Environment == types.ENV.Dev && AppRoot == "" {
		if h, err := os.UserHomeDir(); err == nil {
			AppRoot = filepath.Join(h, Title)
		}
	}

	if ReverseProxy == "true" {
		if util.IsValidPort(ProxyPort) == false {
			logger.Panic("supplied port for reverse proxy is invalid.")
		}

		vars.ReverseProxy = types.BehindReverseProxy{
			StatementValid: true,
			Port:           ProxyPort,
		}
	} else {
		vars.ReverseProxy = types.BehindReverseProxy{
			StatementValid: false,
		}
	}

	os.Setenv("APP_ROOT", AppRoot)
	if Environment == types.ENV.Dev {
		exePath, err := os.Executable()
		if err != nil {
			logger.Panic("could not get executable path:", err)
		}
		os.Setenv("ETC_DIR", filepath.Dir(filepath.Dir(exePath)))
	} else {
		if err := dirExists(ETCDir); err != nil {
			logger.Panic("environment setup not done properly")
		}

		os.Setenv("ETC_DIR", ETCDir)
	}

	os.Setenv("TITLE", Title)
	os.Setenv("DESCRIPTION", Description)
	os.Setenv("VERSION", Version)
	os.Setenv("ENV", Environment)
	os.Setenv("HOST", Host)

	// NOTE: Don't move this function call, look at the info on this function for more details
	createAppRoot()
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	// Handle other potential errors (e.g., permissions)
	logger.Panic("Error checking file `"+filePath+"`\n    ", err)
	return false
}

func dirExists(dirPath string) error {
	info, err := os.Stat(dirPath)
	if err == nil {
		// Path exists; check if it's a directory
		if info.IsDir() {
			return nil // OK: directory exists
		}
		return fmt.Errorf("path exists but is not a directory: %s", dirPath)
	}
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("directory does not exist: %s", dirPath)
	}
	// Other errors (e.g., permission issues)
	return fmt.Errorf("error checking directory %s: %w", dirPath, err)
}

// TODO: Implement authentication/authorization
func main() {
	defer db.Conn.Close()

	// INFO:: startServer checks the current environment configuration.
	//         - In development mode, it starts the server on the DevPort.
	//         - In production mode:
	//           - If behind a reverse proxy, it starts the server on ReverseProxy.Port.
	//           - Otherwise, it starts the server with HTTPS on port 443.
	//           - If SSL certificates are not found, the server crashes.
	startServer := func() (error, string) {
		routeHandler := util.Chain(types.MiddlewareChain{
			Handler: api.HandleRouting(),
			Middlewares: []types.Middleware{
				api.RecoveryMiddleware,
				api.LoggingMiddleware,
			},
		})

		// INFO: Development Server
		if os.Getenv("ENV") == types.ENV.Dev {
			if util.IsValidPort(DevPort) == false {
				logger.Panic("supplied port for dev server is invalid, falling back to :6969")
				DevPort = "6969"
			}

			os.Setenv("port", DevPort)
			logger.Info("Development server started on port :" + DevPort)

			err := http.ListenAndServe(":"+DevPort, routeHandler)

			return err, DevPort
		}

		// INFO: Production server
		fullchain := os.Getenv("FULL_CHAIN")
		privkey := os.Getenv("PRIV_KEY")

		exists1 := fileExists(fullchain)
		exists2 := fileExists(privkey)

		if !exists1 && !exists2 {
			logger.Panic("TLS Certificate could not be found!")
		}

		if vars.ReverseProxy.StatementValid == true {
			os.Setenv("port", vars.ReverseProxy.Port)
			logger.Info("Production server started behind reverse proxy on port :" + vars.ReverseProxy.Port)

			err := http.ListenAndServeTLS(":"+vars.ReverseProxy.Port, fullchain, privkey, routeHandler)

			return err, vars.ReverseProxy.Port
		} else {
			portToStart := "443"
			os.Setenv("port", portToStart)
			logger.Info("Production server started on port :" + portToStart)

			err := http.ListenAndServeTLS(":"+portToStart, fullchain, privkey, routeHandler)

			return err, portToStart
		}
	}

	if err, port := startServer(); err != nil {
		logger.Panic("failed to start server on :" + port)
	}
}
