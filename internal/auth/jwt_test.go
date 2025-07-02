package auth

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	const secret = "test-secret"
	const wrongSecret = "wrong-secret"
	userID := uuid.New()

	tests := []struct {
		name              string
		expiresIn         time.Duration
		validationSecret  string
		expectedErr       error
		expectUserIDMatch bool
	}{
		{
			name:              "Create and Validate JWT - success",
			expiresIn:         time.Hour,
			validationSecret:  secret,
			expectedErr:       nil,
			expectUserIDMatch: true,
		},
		{
			name:             "Expired token",
			expiresIn:        -time.Hour,
			validationSecret: secret,
			expectedErr:      jwt.ErrTokenExpired,
		},
		{
			name:             "Wrong secret",
			expiresIn:        time.Hour,
			validationSecret: wrongSecret,
			expectedErr:      jwt.ErrTokenSignatureInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := MakeJWT(userID, secret, tt.expiresIn)
			if err != nil {
				t.Fatalf("unexpected error creating JWT: %v", err)
			}

			validatedUserID, err := ValidateJWT(tokenString, tt.validationSecret)

			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Fatalf("expected error '%v', but got '%v'", tt.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error validating JWT: %v", err)
			}

			if tt.expectUserIDMatch && validatedUserID != userID {
				t.Errorf("expected user ID %q, got %q", userID, validatedUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name          string
		header        http.Header
		expectedToken string
		expectError   bool
	}{
		{
			name: "Valid token",
			header: http.Header{
				"Authorization": {"Bearer my-token"},
			},
			expectedToken: "my-token",
			expectError:   false,
		},
		{
			name:          "No auth header",
			header:        http.Header{},
			expectedToken: "",
			expectError:   true,
		},
		{
			name: "Malformed header - no Bearer prefix",
			header: http.Header{
				"Authorization": {"my-token"},
			},
			expectedToken: "",
			expectError:   true,
		},
		{
			name: "Valid header - empty token",
			header: http.Header{
				"Authorization": {"Bearer "},
			},
			expectedToken: "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.header)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			if token != tt.expectedToken {
				t.Errorf("expected token '%s', got '%s'", tt.expectedToken, token)
			}
		})
	}
}
