package user

import (
	"context"
	"errors"
	"testing"

	"github.com/cyradin/fixik/internal/db"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Login(t *testing.T) {
	t.Parallel()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	dbUser := db.User{
		ID:       1,
		Name:     "Алиса",
		Username: "alice",
		Password: string(hashedPassword),
	}

	tests := []struct {
		name string
		mock func(*userProviderMock, *tokenManagerMock)
		err  error
	}{
		{
			name: "success",
			mock: func(u *userProviderMock, t *tokenManagerMock) {
				u.getByUsernameFn = func(ctx context.Context, username string) (db.User, error) {
					return dbUser, nil
				}
				t.generateAccessTokenFn = func(userID int64) (string, error) {
					return "access-token", nil
				}
				t.generateRefreshTokenFn = func(userID int64) (string, error) {
					return "refresh-token", nil
				}
			},
		},
		{
			name: "user not found",
			mock: func(u *userProviderMock, t *tokenManagerMock) {
				u.getByUsernameFn = func(ctx context.Context, username string) (db.User, error) {
					return db.User{}, db.ErrNotFound
				}
			},
			err: ErrUserNotFound,
		},
		{
			name: "wrong password",
			mock: func(u *userProviderMock, t *tokenManagerMock) {
				u.getByUsernameFn = func(ctx context.Context, username string) (db.User, error) {
					return dbUser, nil
				}
			},
			err: ErrUserNotFound,
		},
		{
			name: "access token error",
			mock: func(u *userProviderMock, t *tokenManagerMock) {
				u.getByUsernameFn = func(ctx context.Context, username string) (db.User, error) {
					return dbUser, nil
				}
				t.generateAccessTokenFn = func(userID int64) (string, error) {
					return "", errors.New("token error")
				}
			},
			err: errors.New("token error"),
		},
		{
			name: "refresh token error",
			mock: func(u *userProviderMock, t *tokenManagerMock) {
				u.getByUsernameFn = func(ctx context.Context, username string) (db.User, error) {
					return dbUser, nil
				}
				t.generateAccessTokenFn = func(userID int64) (string, error) {
					return "access-token", nil
				}
				t.generateRefreshTokenFn = func(userID int64) (string, error) {
					return "", errors.New("refresh error")
				}
			},
			err: errors.New("refresh error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepo := &userProviderMock{}
			tokenManager := &tokenManagerMock{}

			tt.mock(userRepo, tokenManager)

			svc := NewAuthService(userRepo, tokenManager)

			password := "password"
			if tt.name == "wrong password" {
				password = "wrong"
			}

			res, err := svc.Login(t.Context(), "alice", password)

			if tt.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, int64(1), res.User.ID)
			require.Equal(t, "access-token", res.AccessToken)
			require.Equal(t, "refresh-token", res.RefreshToken)
		})
	}
}

func TestAuthService_Refresh(t *testing.T) {
	t.Parallel()

	dbUser := db.User{
		ID:       1,
		Name:     "Алиса",
		Username: "alice",
		Password: "hashed", // хэш не нужен, мы тут напрямую по ID
	}

	tests := []struct {
		name string
		mock func(*userProviderMock, *tokenManagerMock)
		err  error
	}{
		{
			name: "success",
			mock: func(u *userProviderMock, t *tokenManagerMock) {
				t.validateRefreshTokenFn = func(token string) (int64, error) {
					return 1, nil
				}
				u.getByIDFn = func(ctx context.Context, id int64) (db.User, error) {
					return dbUser, nil
				}
				t.generateAccessTokenFn = func(userID int64) (string, error) {
					return "new-access-token", nil
				}
				t.generateRefreshTokenFn = func(userID int64) (string, error) {
					return "new-refresh-token", nil
				}
			},
		},
		{
			name: "invalid refresh token",
			mock: func(u *userProviderMock, t *tokenManagerMock) {
				t.validateRefreshTokenFn = func(token string) (int64, error) {
					return 0, errors.New("invalid token")
				}
			},
			err: errors.New("invalid refresh token"),
		},
		{
			name: "user not found",
			mock: func(u *userProviderMock, t *tokenManagerMock) {
				t.validateRefreshTokenFn = func(token string) (int64, error) {
					return 1, nil
				}
				u.getByIDFn = func(ctx context.Context, id int64) (db.User, error) {
					return db.User{}, db.ErrNotFound
				}
			},
			err: ErrUserNotFound,
		},
		{
			name: "access token generation error",
			mock: func(u *userProviderMock, t *tokenManagerMock) {
				t.validateRefreshTokenFn = func(token string) (int64, error) {
					return 1, nil
				}
				u.getByIDFn = func(ctx context.Context, id int64) (db.User, error) {
					return dbUser, nil
				}
				t.generateAccessTokenFn = func(userID int64) (string, error) {
					return "", errors.New("access token error")
				}
			},
			err: errors.New("access token error"),
		},
		{
			name: "refresh token generation error",
			mock: func(u *userProviderMock, t *tokenManagerMock) {
				t.validateRefreshTokenFn = func(token string) (int64, error) {
					return 1, nil
				}
				u.getByIDFn = func(ctx context.Context, id int64) (db.User, error) {
					return dbUser, nil
				}
				t.generateAccessTokenFn = func(userID int64) (string, error) {
					return "access", nil
				}
				t.generateRefreshTokenFn = func(userID int64) (string, error) {
					return "", errors.New("refresh token error")
				}
			},
			err: errors.New("refresh token error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepo := &userProviderMock{}
			tokenManager := &tokenManagerMock{}

			tt.mock(userRepo, tokenManager)

			svc := NewAuthService(userRepo, tokenManager)

			res, err := svc.Refresh(context.Background(), "old-refresh-token")

			if tt.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, int64(1), res.User.ID)
			require.Equal(t, "new-access-token", res.AccessToken)
			require.Equal(t, "new-refresh-token", res.RefreshToken)
		})
	}
}

type userProviderMock struct {
	getByUsernameFn func(ctx context.Context, username string) (db.User, error)
	getByIDFn       func(ctx context.Context, id int64) (db.User, error)
}

func (m *userProviderMock) GetByUsername(ctx context.Context, username string) (db.User, error) {
	return m.getByUsernameFn(ctx, username)
}

func (m *userProviderMock) GetByID(ctx context.Context, id int64) (db.User, error) {
	return m.getByIDFn(ctx, id)
}

type tokenManagerMock struct {
	generateAccessTokenFn  func(userID int64) (string, error)
	generateRefreshTokenFn func(userID int64) (string, error)
	validateRefreshTokenFn func(tokenStr string) (int64, error)
}

func (m *tokenManagerMock) GenerateAccessToken(userID int64) (string, error) {
	if m.generateAccessTokenFn == nil {
		return "", nil
	}

	return m.generateAccessTokenFn(userID)
}

func (m *tokenManagerMock) GenerateRefreshToken(userID int64) (string, error) {
	if m.generateRefreshTokenFn == nil {
		return "", nil
	}

	return m.generateRefreshTokenFn(userID)
}

func (m *tokenManagerMock) ValidateAccessToken(tokenStr string) (int64, error) {
	return 0, nil
}

func (m *tokenManagerMock) ValidateRefreshToken(tokenStr string) (int64, error) {
	if m.validateRefreshTokenFn == nil {
		return 0, nil
	}

	return m.validateRefreshTokenFn(tokenStr)
}
