package middleware

import (
	"net/http"
)

// CORSMiddleware applies CORS headers to HTTP responses
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set allowed origins
		w.Header().Set("Access-Control-Allow-Origin", "*") // ⚠️ change "*" to your frontend domain in production

		// Set allowed methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Set allowed headers
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
