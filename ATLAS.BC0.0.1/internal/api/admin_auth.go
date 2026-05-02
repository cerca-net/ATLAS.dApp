package api

import (
	"crypto/subtle"
	"log"
	"net/http"
	"os"
)

// AdminAuthMiddleware protects admin-only endpoints with an API key.
//
// The key is read from the ADMIN_API_KEY environment variable.
// Clients must send it via the Authorization header:
//   Authorization: Bearer <key>
//
// If ADMIN_API_KEY is not set, admin endpoints are OPEN (dev mode only)
// and a warning is logged at startup.
func AdminAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	adminKey := os.Getenv("ADMIN_API_KEY")
	if adminKey == "" {
		log.Printf("⚠️  ADMIN_API_KEY not set — admin endpoints are UNPROTECTED. Set this in production!")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// If no key is configured, allow all requests (dev mode)
		if adminKey == "" {
			next(w, r)
			return
		}

		// Extract Bearer token
		auth := r.Header.Get("Authorization")
		const prefix = "Bearer "
		if len(auth) <= len(prefix) {
			http.Error(w, `{"error":"unauthorized","message":"Missing or invalid Authorization header. Use: Bearer <ADMIN_API_KEY>"}`, http.StatusUnauthorized)
			return
		}

		token := auth[len(prefix):]

		// Constant-time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(token), []byte(adminKey)) != 1 {
			http.Error(w, `{"error":"forbidden","message":"Invalid admin API key"}`, http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
