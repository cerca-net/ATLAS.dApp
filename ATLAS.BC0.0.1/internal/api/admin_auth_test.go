package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestAdminAuthMiddlewareFallback(t *testing.T) {
	// Set up environment variables
	os.Setenv("SUPABASE_JWT_SECRET", "")
	os.Setenv("ADMIN_API_KEY", "test-secret-key")
	defer os.Setenv("ADMIN_API_KEY", "")

	// Create handler wrapped in middleware
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})
	authHandler := AdminAuthMiddleware(dummyHandler)

	// Case 1: Missing Authorization Header
	req := httptest.NewRequest("GET", "/admin/faucet", nil)
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status Unauthorized (401), got %v", w.Code)
	}

	// Case 2: Valid API Key
	req = httptest.NewRequest("GET", "/admin/faucet", nil)
	req.Header.Set("Authorization", "Bearer test-secret-key")
	w = httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK (200), got %v", w.Code)
	}
	if w.Body.String() != "success" {
		t.Errorf("Expected body 'success', got '%s'", w.Body.String())
	}

	// Case 3: Invalid API Key
	req = httptest.NewRequest("GET", "/admin/faucet", nil)
	req.Header.Set("Authorization", "Bearer wrong-secret-key")
	w = httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status Forbidden (403), got %v", w.Code)
	}
}

func TestAdminAuthMiddlewareSupabaseJWT(t *testing.T) {
	jwtSecret := "supabase-test-jwt-secret-key-12345"
	os.Setenv("SUPABASE_JWT_SECRET", jwtSecret)
	defer os.Setenv("SUPABASE_JWT_SECRET", "")

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})
	authHandler := AdminAuthMiddleware(dummyHandler)

	createToken := func(role string, claimsLocation string, expired bool) string {
		claims := jwt.MapClaims{
			"exp": time.Now().Add(time.Hour).Unix(),
			"sub": "user-uuid",
		}
		if expired {
			claims["exp"] = time.Now().Add(-time.Hour).Unix()
		}

		metadata := map[string]interface{}{
			"role": role,
		}

		if claimsLocation == "app_metadata" {
			claims["app_metadata"] = metadata
		} else if claimsLocation == "user_metadata" {
			claims["user_metadata"] = metadata
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, _ := token.SignedString([]byte(jwtSecret))
		return tokenStr
	}

	// Case 1: Valid JWT with app_metadata role=admin
	validAppToken := createToken("admin", "app_metadata", false)
	req := httptest.NewRequest("GET", "/admin/faucet", nil)
	req.Header.Set("Authorization", "Bearer "+validAppToken)
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected OK (200) for app_metadata admin token, got %v", w.Code)
	}

	// Case 2: Valid JWT with user_metadata role=admin
	validUserToken := createToken("admin", "user_metadata", false)
	req = httptest.NewRequest("GET", "/admin/faucet", nil)
	req.Header.Set("Authorization", "Bearer "+validUserToken)
	w = httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected OK (200) for user_metadata admin token, got %v", w.Code)
	}

	// Case 3: Valid JWT with non-admin role
	invalidRoleToken := createToken("user", "app_metadata", false)
	req = httptest.NewRequest("GET", "/admin/faucet", nil)
	req.Header.Set("Authorization", "Bearer "+invalidRoleToken)
	w = httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected Forbidden (403) for user role token, got %v", w.Code)
	}

	// Case 4: Expired token
	expiredToken := createToken("admin", "app_metadata", true)
	req = httptest.NewRequest("GET", "/admin/faucet", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)
	w = httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected Forbidden (403) for expired token, got %v", w.Code)
	}

	// Case 5: Incorrect signing key
	badToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
		"app_metadata": map[string]interface{}{
			"role": "admin",
		},
	})
	badTokenStr, _ := badToken.SignedString([]byte("wrong-key-secret"))
	req = httptest.NewRequest("GET", "/admin/faucet", nil)
	req.Header.Set("Authorization", "Bearer "+badTokenStr)
	w = httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected Forbidden (403) for bad signature token, got %v", w.Code)
	}
}
