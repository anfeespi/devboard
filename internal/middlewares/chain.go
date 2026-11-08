package middlewares

import (
	"net/http"
)

// Middleware signature of any middleware
type Middleware func(http.Handler) http.Handler

// Chain apply multi middleware
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}

	return h
}
