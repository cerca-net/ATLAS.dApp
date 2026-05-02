package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// Validation constants
const (
	MaxAddressLength    = 128
	MaxHashLength       = 128
	MaxMnemonicLength   = 512
	MaxDataFieldLength  = 4096
	MaxTransactionField = 1024
)

// Validation patterns
var (
	hexPattern     = regexp.MustCompile(`^(0x)?[0-9a-fA-F]+$`)
	addressPattern = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
)

// ValidationError represents a request validation failure.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateAddress checks that an address looks like a valid hex address.
func ValidateAddress(address string) *ValidationError {
	if address == "" {
		return &ValidationError{Field: "address", Message: "address is required"}
	}
	if len(address) > MaxAddressLength {
		return &ValidationError{Field: "address", Message: "address exceeds maximum length"}
	}
	if !hexPattern.MatchString(address) {
		return &ValidationError{Field: "address", Message: "address must be a hex string"}
	}
	return nil
}

// ValidateHash checks that a hash is a valid hex string.
func ValidateHash(hash string) *ValidationError {
	if hash == "" {
		return &ValidationError{Field: "hash", Message: "hash is required"}
	}
	if len(hash) > MaxHashLength {
		return &ValidationError{Field: "hash", Message: "hash exceeds maximum length"}
	}
	if !hexPattern.MatchString(hash) {
		return &ValidationError{Field: "hash", Message: "hash must be a hex string"}
	}
	return nil
}

// ValidateAmount checks that an amount is positive.
func ValidateAmount(amount int64) *ValidationError {
	if amount < 0 {
		return &ValidationError{Field: "amount", Message: "amount must be non-negative"}
	}
	return nil
}

// ValidateStringField checks a generic string field for length and dangerous content.
func ValidateStringField(name, value string, maxLen int) *ValidationError {
	if len(value) > maxLen {
		return &ValidationError{Field: name, Message: fmt.Sprintf("%s exceeds maximum length of %d", name, maxLen)}
	}
	// Basic injection prevention — reject strings with script tags or SQL-like patterns
	lower := strings.ToLower(value)
	dangerousPatterns := []string{"<script", "javascript:", "onclick=", "onerror=", "'; drop", "\" or 1=1", "union select"}
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lower, pattern) {
			return &ValidationError{Field: name, Message: fmt.Sprintf("%s contains disallowed content", name)}
		}
	}
	return nil
}

// WriteValidationError sends a 400 response with the validation error.
func WriteValidationError(w http.ResponseWriter, ve *ValidationError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":      "validation_error",
		"field":      ve.Field,
		"message":    ve.Message,
	})
}

// RequireMethod validates the HTTP method. Returns true if the method matches.
func RequireMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		http.Error(w, fmt.Sprintf(`{"error":"method_not_allowed","message":"Expected %s, got %s"}`, method, r.Method), http.StatusMethodNotAllowed)
		return false
	}
	return true
}

// RequireQueryParam extracts and validates a required query parameter.
// Returns empty string and writes error if missing.
func RequireQueryParam(w http.ResponseWriter, r *http.Request, name string) string {
	val := r.URL.Query().Get(name)
	if val == "" {
		WriteValidationError(w, &ValidationError{Field: name, Message: fmt.Sprintf("query parameter '%s' is required", name)})
		return ""
	}
	return val
}
