package parser

import (
	"net/http"
	"strings"
)

var DynamicRoutes []string = []string{"/blog/"}

type RouteData struct {
	FullPath string
	Route    string
}

func PerformSSR(fullPath string) (string, error, int) {
	var route RouteData

	for _, DynamicRoute := range DynamicRoutes {
		if strings.HasPrefix(fullPath, DynamicRoute) {
			route = RouteData{
				FullPath: fullPath,
				Route:    DynamicRoute,
			}
			break
		}
	}

	switch route.Route {
	default:
		return "", nil, http.StatusNoContent
	}
}
