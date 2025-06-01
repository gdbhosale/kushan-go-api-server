package middlewares

import (
	"net/http"
)

// Type Middleware to create Stack
type Middleware func(http.Handler) http.Handler

// Method to to create stack of middlewares
func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}

		return next
	}
}
