package middlewares

import (
	"goat/internal"

	"net/http"
	"strings"
)

func EnsureAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		internal.Debug("Middlewares::EnsureAdmin", "Checking if user is admin")
		if !strings.Contains(r.Header.Get("Authorization"), "Admin") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func LoadUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		internal.Debug("Middlewares::EnsureAdmin", "Loading user")
		next.ServeHTTP(w, r)
	})
}

func CheckPermissions(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		internal.Debug("Middlewares::EnsureAdmin", "Checking Permissions")
		next.ServeHTTP(w, r)
	})
}
