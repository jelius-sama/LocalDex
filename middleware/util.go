package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

type MiddlewareChain struct {
	Handler     http.Handler
	Middlewares []Middleware
}

func Chain(cm MiddlewareChain) http.Handler {
	for i := len(cm.Middlewares) - 1; i >= 0; i-- {
		cm.Handler = cm.Middlewares[i](cm.Handler)
	}
	return cm.Handler
}
