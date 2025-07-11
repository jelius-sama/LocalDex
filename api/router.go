package api

import (
	m "LocalDex/middleware"
	"context"
	"net/http"
	"path"
	"sort"
	"strings"
)

type ctxKey string

const paramKey ctxKey = "routeParams"

func WithParam(ctx context.Context, key, value string) context.Context {
	m, ok := ctx.Value(paramKey).(map[string]string)
	if !ok {
		m = make(map[string]string)
	}
	m[key] = value
	return context.WithValue(ctx, paramKey, m)
}

func Param(r *http.Request, key string) string {
	if m, ok := r.Context().Value(paramKey).(map[string]string); ok {
		return m[key]
	}
	return ""
}

type HandlerFunc http.HandlerFunc

type Router struct {
	routes          []route
	middlewares     []m.Middleware
	NotFoundHandler http.HandlerFunc
}

type route struct {
	method  string
	pattern string
	handler http.HandlerFunc
}

// Creates a new router
func New() *Router {
	return &Router{routes: []route{}}
}

func cleanPath(p string) string {
	cleaned := path.Clean(p)
	if cleaned == "." {
		return "/"
	}
	return cleaned
}

func isDynamicSegment(part string) bool {
	return strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}")
}

func countDynamicSegments(pattern string) int {
	count := 0
	for _, part := range strings.Split(pattern, "/") {
		if isDynamicSegment(part) {
			count++
		}
	}
	return count
}

// Matches dynamic and static paths
func matchPattern(pattern, path string) (bool, map[string]string) {
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return false, nil
	}

	params := make(map[string]string)
	for i := range patternParts {
		if strings.HasPrefix(patternParts[i], "{") && strings.HasSuffix(patternParts[i], "}") {
			paramName := patternParts[i][1 : len(patternParts[i])-1]
			params[paramName] = pathParts[i]
		} else if patternParts[i] != pathParts[i] {
			return false, nil
		}
	}
	return true, params
}

// The actual HTTP handler that will be passed to http.Server
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	reqPath := cleanPath(req.URL.Path)

	for _, rt := range r.routes {
		if req.Method == rt.method {
			if ok, params := matchPattern(cleanPath(rt.pattern), reqPath); ok {
				ctx := req.Context()
				for k, v := range params {
					ctx = WithParam(ctx, k, v)
				}
				req = req.WithContext(ctx)
				rt.handler.ServeHTTP(w, req)
				return
			}
		}
	}

	// Call custom 404 handler if set
	if r.NotFoundHandler != nil {
		r.NotFoundHandler(w, req)
		return
	}

	// Fallback default
	http.NotFound(w, req)
}

// Register middlewares
func (r *Router) Use(mw ...m.Middleware) {
	r.middlewares = append(r.middlewares, mw...)
}

// Register a route
func (r *Router) Handle(method, pattern string, handler http.HandlerFunc) {
	pattern = cleanPath(pattern)

	finalHandler := http.Handler(handler)
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		finalHandler = r.middlewares[i](finalHandler)
	}
	r.routes = append(r.routes, route{
		method:  method,
		pattern: pattern,
		handler: finalHandler.ServeHTTP,
	})

	// Sort: static routes first (fewer dynamic segments)
	sort.Slice(r.routes, func(i, j int) bool {
		return countDynamicSegments(r.routes[i].pattern) < countDynamicSegments(r.routes[j].pattern)
	})
}

// Convenience methods
func (r *Router) Get(pattern string, handler http.HandlerFunc) {
	r.Handle(http.MethodGet, pattern, handler)
}
func (r *Router) Post(pattern string, handler http.HandlerFunc) {
	r.Handle(http.MethodPost, pattern, handler)
}
func (r *Router) Put(pattern string, handler http.HandlerFunc) {
	r.Handle(http.MethodPut, pattern, handler)
}
func (r *Router) Delete(pattern string, handler http.HandlerFunc) {
	r.Handle(http.MethodDelete, pattern, handler)
}

// Nesting
func (r *Router) Route(prefix string, fn func(sub *Router)) {
	sub := &Router{}
	fn(sub)
	for _, rt := range sub.routes {
		r.Handle(rt.method, prefix+rt.pattern, rt.handler)
	}
}
