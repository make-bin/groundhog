// @AI_GENERATED
package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/make-bin/groundhog/pkg/utils/config"
)

// Sentinel errors for JWT validation failures.
var (
	ErrTokenExpired = errors.New("token has expired")
	ErrTokenInvalid = errors.New("token is invalid")
)

// Claims represents the JWT claims payload.
type Claims struct {
	PrincipalID string `json:"principal_id"`
	jwt.RegisteredClaims
}

// JWTService handles JWT token generation and validation.
type JWTService struct {
	secret []byte
	ttl    time.Duration
}

// NewJWTService creates a new JWTService from the provided JWT configuration.
func NewJWTService(cfg *config.JWTConfig) *JWTService {
	return &JWTService{
		secret: []byte(cfg.Secret),
		ttl:    cfg.AccessTokenTTL,
	}
}

// GenerateToken creates a signed JWT for the given principal ID.
// If ttl is 0, the service's default TTL is used.
func (s *JWTService) GenerateToken(principalID string, ttl time.Duration) (string, error) {
	if ttl == 0 {
		ttl = s.ttl
	}

	now := time.Now()
	claims := Claims{
		PrincipalID: principalID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			Issuer:    "openclaw",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// ValidateToken parses and validates the given token string.
// Returns the claims on success, ErrTokenExpired for expired tokens,
// and ErrTokenInvalid for any other validation failure.
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return s.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

// @AI_GENERATED: end
