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
		currentAdminKey := os.Getenv("ADMIN_API_KEY")
		if currentAdminKey == "" {
			currentAdminKey = "cerca-dev-admin-secret-key"
		}

		// 1. Try static admin API key first (constant-time comparison)
		if subtle.ConstantTimeCompare([]byte(tokenStr), []byte(currentAdminKey)) == 1 {
			next(w, r)
			return
		}

		// 2. Try Supabase JWT if configured
		currentSupabaseSecret := os.Getenv("SUPABASE_JWT_SECRET")
		if currentSupabaseSecret != "" {
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(currentSupabaseSecret), nil
			})
			if err == nil && token.Valid {
				claims, ok := token.Claims.(jwt.MapClaims)
				if ok {
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

					if isAdmin {
						next(w, r)
						return
					}
				}
			}
			
			// If we got here, it's either an invalid token or not an admin
			log.Printf("⚠️ Auth failed: Supabase JWT parsing or role check failed: %v", err)
		}

		// Deny access if neither succeeded
		http.Error(w, `{"error":"forbidden","message":"Invalid admin API key or Supabase JWT"}`, http.StatusForbidden)
	}
}
