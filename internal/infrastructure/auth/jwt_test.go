package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/infrastructure/config"
)

func TestNewJWTSigner(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:     "test-secret",
		JWTExpiration: time.Hour,
	}

	signer := NewJWTSigner(cfg)
	assert.NotNil(t, signer)

	// Test that it implements the interface
	jwtSignerImpl, ok := signer.(*jwtSigner)
	assert.True(t, ok)
	assert.Equal(t, []byte("test-secret"), jwtSignerImpl.secret)
	assert.Equal(t, time.Hour, jwtSignerImpl.exp)
}

func TestJWTSigner_Sign(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:     "test-secret",
		JWTExpiration: time.Hour,
	}
	signer := NewJWTSigner(cfg)

	tests := []struct {
		name   string
		userID int64
	}{
		{
			name:   "should sign token successfully",
			userID: 123,
		},
		{
			name:   "should sign token for user ID 1",
			userID: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := signer.Sign(tt.userID)
			assert.NoError(t, err)
			assert.NotEmpty(t, token)

			// Verify the token can be parsed back
			userID, err := signer.Verify(token)
			assert.NoError(t, err)
			assert.Equal(t, tt.userID, userID)
		})
	}
}

func TestJWTSigner_Verify(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:     "test-secret",
		JWTExpiration: time.Hour,
	}
	signer := NewJWTSigner(cfg)

	tests := []struct {
		name        string
		setupToken  func() string
		expectedID  int64
		expectError bool
	}{
		{
			name: "should verify valid token successfully",
			setupToken: func() string {
				token, _ := signer.Sign(123)
				return token
			},
			expectedID:  123,
			expectError: false,
		},
		{
			name: "should return error for invalid token",
			setupToken: func() string {
				return "invalid.token.here"
			},
			expectedID:  0,
			expectError: true,
		},
		{
			name: "should return error for empty token",
			setupToken: func() string {
				return ""
			},
			expectedID:  0,
			expectError: true,
		},
		{
			name: "should return error for token with wrong signing method",
			setupToken: func() string {
				claims := Claims{
					UserID: "123",
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
						IssuedAt:  jwt.NewNumericDate(time.Now()),
					},
				}
				// Use non-HMAC signing method (should trigger signing method validation error)
				token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
				tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
				return tokenString
			},
			expectedID:  0,
			expectError: true,
		},
		{
			name: "should return error for token with wrong secret",
			setupToken: func() string {
				wrongSigner := &jwtSigner{
					secret: []byte("wrong-secret"),
					exp:    time.Hour,
				}
				token, _ := wrongSigner.Sign(123)
				return token
			},
			expectedID:  0,
			expectError: true,
		},
		{
			name: "should return error for expired token",
			setupToken: func() string {
				expiredSigner := &jwtSigner{
					secret: []byte("test-secret"),
					exp:    -time.Hour, // Expired
				}
				token, _ := expiredSigner.Sign(123)
				return token
			},
			expectedID:  0,
			expectError: true,
		},
		{
			name: "should return error for token with invalid user_id",
			setupToken: func() string {
				claims := Claims{
					UserID: "invalid-user-id",
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
						IssuedAt:  jwt.NewNumericDate(time.Now()),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte("test-secret"))
				return tokenString
			},
			expectedID:  0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.setupToken()
			userID, err := signer.Verify(token)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedID, userID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, userID)
			}
		})
	}
}
