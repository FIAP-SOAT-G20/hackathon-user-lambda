package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/core/port"
	"github.com/FIAP-SOAT-G20/hackathon-user-lambda/internal/infrastructure/config"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type jwtSigner struct {
	secret []byte
	exp    time.Duration
}

func NewJWTSigner(cfg *config.Config) port.JWTSigner {
	return &jwtSigner{secret: []byte(cfg.JWTSecret), exp: cfg.JWTExpiration}
}

// ensure implementation
var _ port.JWTSigner = (*jwtSigner)(nil)

func (j *jwtSigner) Sign(userID int64) (string, error) {
	claims := Claims{
		UserID: strconv.FormatInt(userID, 10),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.exp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *jwtSigner) Verify(tokenStr string) (int64, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return j.secret, nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}
	if userID, err := strconv.ParseInt(claims.UserID, 10, 64); err == nil {
		return userID, nil
	}
	return 0, errors.New("invalid user_id in token")
}
