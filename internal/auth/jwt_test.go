package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	const secret = "test-secret"
	const wrongSecret = "wrong-secret"
	userID := uuid.New()

	t.Run("Create and Validate JWT - success", func(t *testing.T) {
		tokenString, err := MakeJWT(userID, secret, time.Hour)
		if err != nil {
			t.Fatalf("unexpected error creating JWT: %v", err)
		}

		validatedUserID, err := ValidateJWT(tokenString, secret)
		if err != nil {
			t.Fatalf("unexpected error validating JWT: %v", err)
		}

		if validatedUserID != userID {
			t.Errorf("expected user ID %q, got %q", userID, validatedUserID)
		}
	})

	t.Run("Expired token", func(t *testing.T) {
		tokenString, err := MakeJWT(userID, secret, -time.Hour)
		if err != nil {
			t.Fatalf("unexpected error creating expired JWT: %v", err)
		}

		_, err = ValidateJWT(tokenString, secret)
		if !errors.Is(err, jwt.ErrTokenExpired) {
			t.Fatalf("expected error for expired token, but got: %v", err)
		}
	})

	t.Run("Wrong secret", func(t *testing.T) {
		tokenString, err := MakeJWT(userID, secret, time.Hour)
		if err != nil {
			t.Fatalf("unexpected error creating JWT: %v", err)
		}

		_, err = ValidateJWT(tokenString, wrongSecret)
		if !errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			t.Fatalf("expected error for wrong secret, but got: %v", err)
		}
	})
}
