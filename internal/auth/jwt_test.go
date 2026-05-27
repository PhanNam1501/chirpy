package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret-key"
	expiresIn := 1 * time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	if token == "" {
		t.Fatal("Expected non-empty token")
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret-key"
	expiresIn := 1 * time.Hour

	// Create token
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	// Validate token
	validatedID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}

	if validatedID != userID {
		t.Fatalf("Expected %v, got %v", userID, validatedID)
	}
}

func TestExpiredJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret-key"
	expiresIn := -1 * time.Hour // Expired token

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	// Try to validate expired token
	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Fatal("Expected error for expired token")
	}
}

func TestWrongSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret-key"
	wrongSecret := "wrong-secret-key"
	expiresIn := 1 * time.Hour

	// Create token with correct secret
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	// Try to validate with wrong secret
	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatal("Expected error when validating with wrong secret")
	}
}
