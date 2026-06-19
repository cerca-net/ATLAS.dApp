package api

import (
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// AdminAuthMiddleware protects admin-only endpoints with an API key or Supabase JWT.
//
// If SUPABASE_JWT_SECRET environment variable is set:
//   Clients must pass a Supabase User JWT in the Authorization header:
//     Authorization: Bearer <JWT>
//   The middleware verifies the JWT signature and checks that the user's role is "admin"
//   under app_metadata or user_metadata.
//
// If SUPABASE_JWT_SECRET is NOT set, it falls back to ADMIN_API_KEY (for development):
//   Authorization: Bearer <ADMIN_API_KEY>
//
// If ADMIN_API_KEY is also not set, it uses the default "cerca-dev-admin-secret-key".
func AdminAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// Log startup warning if both keys are missing in the environment
	adminKey := os.Getenv("ADMIN_API_KEY")
	supabaseSecret := os.Getenv("SUPABASE_JWT_SECRET")
	if adminKey == "" && supabaseSecret == "" {
		log.Printf("⚠️  ADMIN_API_KEY and SUPABASE_JWT_SECRET not set — using default dev key: cerca-dev-admin-secret-key")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Extract Bearer token
		auth := r.Header.Get("Authorization")
		const prefix = "Bearer "
		if len(auth) <= len(prefix) || !strings.HasPrefix(auth, prefix) {
			http.Error(w, `{"error":"unauthorized","message":"Missing or invalid Authorization header. Use: Bearer <token>"}`, http.StatusUnauthorized)
			return
		}

		tokenStr := auth[len(prefix):]

		// Fetch env keys dynamically inside the handler to adapt to config changes/tests
		currentSupabaseSecret := os.Getenv("SUPABASE_JWT_SECRET")

		if currentSupabaseSecret != "" {
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(currentSupabaseSecret), nil
			})
			if err != nil || !token.Valid {
				log.Printf("⚠️ Auth failed: Invalid token: %v", err)
				http.Error(w, fmt.Sprintf(`{"error":"forbidden","message":"Invalid token: %v"}`, err), http.StatusForbidden)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, `{"error":"forbidden","message":"Invalid token claims"}`, http.StatusForbidden)
				return
			}

			// Check role in app_metadata or user_metadata
			isAdmin := false

			// Check app_metadata
			if appMetadata, ok := claims["app_metadata"].(map[string]interface{}); ok {
				if role, ok := appMetadata["role"].(string); ok && role == "admin" {
					isAdmin = true
				}
			}

			// Check user_metadata as fallback
			if !isAdmin {
				if userMetadata, ok := claims["user_metadata"].(map[string]interface{}); ok {
					if role, ok := userMetadata["role"].(string); ok && role == "admin" {
						isAdmin = true
					}
				}
			}

			if !isAdmin {
				log.Printf("⚠️ Auth failed: Token does not have admin role")
				http.Error(w, `{"error":"forbidden","message":"Access denied. Admin role required."}`, http.StatusForbidden)
				return
			}

			// Authorized successfully!
			next(w, r)
			return
		}

		// Fallback to static ADMIN_API_KEY
		currentAdminKey := os.Getenv("ADMIN_API_KEY")
		if currentAdminKey == "" {
			currentAdminKey = "cerca-dev-admin-secret-key"
		}

		// Constant-time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(tokenStr), []byte(currentAdminKey)) != 1 {
			http.Error(w, `{"error":"forbidden","message":"Invalid admin API key"}`, http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
