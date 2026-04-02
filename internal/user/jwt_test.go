package user

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTManager_GenerateAccessToken(t *testing.T) {
	t.Parallel()

	secret := "test-secret"

	tests := []struct {
		name   string
		userID int64
		ttl    time.Duration
	}{
		{
			name:   "success",
			userID: 42,
			ttl:    15 * time.Minute,
		},
		{
			name:   "another user",
			userID: 999,
			ttl:    1 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			manager := NewJWTManager(secret, tt.ttl, time.Hour)

			tokenStr, err := manager.GenerateAccessToken(tt.userID)
			require.NoError(t, err)
			require.NotEmpty(t, tokenStr)

			token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
				return []byte(secret), nil
			})
			require.NoError(t, err)
			require.True(t, token.Valid)

			claims, ok := token.Claims.(*Claims)
			require.True(t, ok)

			require.Equal(t, tt.userID, claims.UserID)

			require.WithinDuration(
				t,
				time.Now().Add(tt.ttl),
				claims.ExpiresAt.Time,
				time.Second,
			)
		})
	}
}

func TestJWTManager_GenerateRefreshToken(t *testing.T) {
	t.Parallel()

	secret := "test-secret"

	tests := []struct {
		name   string
		userID int64
		ttl    time.Duration
	}{
		{
			name:   "success",
			userID: 1,
			ttl:    24 * time.Hour,
		},
		{
			name:   "long ttl",
			userID: 2,
			ttl:    7 * 24 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			manager := NewJWTManager(secret, time.Minute, tt.ttl)

			tokenStr, err := manager.GenerateRefreshToken(tt.userID)
			require.NoError(t, err)
			require.NotEmpty(t, tokenStr)

			token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
				return []byte(secret), nil
			})
			require.NoError(t, err)
			require.True(t, token.Valid)

			claims, ok := token.Claims.(*Claims)
			require.True(t, ok)

			require.Equal(t, tt.userID, claims.UserID)

			require.WithinDuration(
				t,
				time.Now().Add(tt.ttl),
				claims.ExpiresAt.Time,
				time.Second,
			)
		})
	}
}
