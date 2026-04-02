package user

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret          string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewJWTManager(
	secret string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *JWTManager {
	return &JWTManager{
		secret:          secret,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (j *JWTManager) GenerateAccessToken(userID int64) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secret))
}

func (j *JWTManager) GenerateRefreshToken(userID int64) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secret))
}

func (j *JWTManager) ValidateAccessToken(tokenStr string) (int64, error) {
	return j.validateToken(tokenStr)
}

func (j *JWTManager) ValidateRefreshToken(tokenStr string) (int64, error) {
	return j.validateToken(tokenStr)
}

func (j *JWTManager) validateToken(tokenStr string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}

		return []byte(j.secret), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return 0, jwt.ErrTokenInvalidClaims
	}

	if claims.ExpiresAt == nil || claims.ExpiresAt.Before(time.Now()) {
		return 0, jwt.ErrTokenExpired
	}

	return claims.UserID, nil
}
