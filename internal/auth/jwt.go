package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ErrTokenExpired token expired
var ErrTokenExpired = errors.New("expired token")

// ErrTokenInvalid invalid token
var ErrTokenInvalid = errors.New("invalid token")

// Claims the content of the JWT
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTConfig contain the configuration of the manager
type JWTConfig struct {
	SecretKey            []byte
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

// JWTManager manage the creation and verification of tokens
type JWTManager struct {
	cfg JWTConfig
}

// NewJWTManager create a new manager
func NewJWTManager(cfg JWTConfig) (*JWTManager, error) {
	if len(cfg.SecretKey) < 32 {
		return nil, fmt.Errorf("secret key needs to have min 32 bytes, now have %d", len(cfg.SecretKey))
	}

	return &JWTManager{cfg}, nil
}

// GenerateAccessToken create a new access token
func (m *JWTManager) GenerateAccessToken(userID, email, role string) (string, error) {
	return m.generateToken(userID, email, role, m.cfg.AccessTokenDuration)
}

// GenerateRefreshToken create a new refresh token
func (m *JWTManager) GenerateRefreshToken(userID, email, role string) (string, error) {
	return m.generateToken(userID, email, role, m.cfg.RefreshTokenDuration)
}

// VerifyToken verify the signature and extract the claims
func (m *JWTManager) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("algorithm not expected: %v", t.Header["alg"])
			}

			return m.cfg.SecretKey, nil
		},
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithIssuer("devboard"),
	)

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

// generateToken generate jwt token
func (m *JWTManager) generateToken(userID, email, role string, duration time.Duration) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        generateTokenID(),
			Issuer:    "devboard",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.cfg.SecretKey)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return signed, nil
}

func generateTokenID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
