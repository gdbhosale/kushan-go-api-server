package middlewares

import "net/http"

func AllowCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Allow all origins (you may want to restrict this to specific origins)
		switch origin := r.Header.Get("Origin"); origin {
		case "http://localhost:8083", "http://localhost:8084", "https://kushant-go.dwij.in":
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// Allow specific methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Allow Credentials
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Allow specific headers
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		// Handle preflight requests (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continue processing the request
		next.ServeHTTP(w, r)
	})
}
