package middleware

import (
	"net/http"
	"os"
	"strings"
)

// CORS adds CORS headers with origin whitelisting.
// For production, set ALLOWED_ORIGINS env var with comma-separated origins.
func CORS(next http.Handler) http.Handler {
	// Build allowed origins map from environment variable
	allowedOrigins := make(map[string]bool)

	// Default allowed origins for development
	allowedOrigins["http://localhost:3000"] = true
	allowedOrigins["http://127.0.0.1:3000"] = true

	// Add origins from environment variable
	envOrigins := os.Getenv("ALLOWED_ORIGINS")
	if envOrigins != "" {
		for _, origin := range strings.Split(envOrigins, ",") {
			origin = strings.TrimSpace(origin)
			if origin != "" {
				allowedOrigins[origin] = true
			}
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Only set CORS headers if origin is in whitelist
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "3600")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
