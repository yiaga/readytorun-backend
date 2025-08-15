package middleware

import "net/http"

func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get the origin from the request header
        origin := r.Header.Get("Origin")

        // Define your allowed frontend origins
        allowedOrigins := map[string]bool{
            "https://readytorun.vercel.app": true,
            "https://readytorunng.org": true,
            "http://localhost:8080": true,
        }

        // Check if the request origin is in the allowed list
        if allowedOrigins[origin] {
            w.Header().Set("Access-Control-Allow-Origin", origin)
        }

        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Access-Control-Allow-Credentials", "true")

        // Handle preflight requests (OPTIONS method)
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }

        next.ServeHTTP(w, r)
    })
}