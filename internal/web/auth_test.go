package web

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cyradin/fixik/internal/user"
	"github.com/stretchr/testify/require"
)

func TestAuthLogin(t *testing.T) {
	t.Parallel()

	accessTTL := 15 * time.Minute
	refreshTTL := 24 * time.Hour

	tests := []struct {
		name   string
		req    AuthLoginRequest
		mock   func(*mockAuthService)
		status int
		check  func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			req: AuthLoginRequest{
				Username: "alice",
				Password: "pass",
			},
			mock: func(m *mockAuthService) {
				m.loginFn = func(ctx context.Context, username, password string) (user.LoginResult, error) {
					return user.LoginResult{
						User:         user.User{ID: 1},
						AccessToken:  "access-token",
						RefreshToken: "refresh-token",
					}, nil
				}
			},
			status: http.StatusOK,
			check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				t.Helper()

				cookies := rr.Result().Cookies()
				require.Len(t, cookies, 2)

				var access, refresh *http.Cookie

				for _, c := range cookies {
					if c.Name == "access_token" {
						access = c
					}

					if c.Name == "refresh_token" {
						refresh = c
					}
				}

				require.NotNil(t, access)
				require.Equal(t, "access-token", access.Value)
				require.True(t, access.HttpOnly)
				require.Equal(t, "/", access.Path)

				require.NotNil(t, refresh)
				require.Equal(t, "refresh-token", refresh.Value)
				require.True(t, refresh.HttpOnly)
				require.Equal(t, "/api/auth/refresh", refresh.Path)
			},
		},
		{
			name:   "validation error",
			req:    AuthLoginRequest{},
			mock:   func(m *mockAuthService) {},
			status: http.StatusBadRequest,
		},
		{
			name: "user not found",
			req: AuthLoginRequest{
				Username: "alice",
				Password: "wrong",
			},
			mock: func(m *mockAuthService) {
				m.loginFn = func(ctx context.Context, username, password string) (user.LoginResult, error) {
					return user.LoginResult{}, user.ErrUserNotFound
				}
			},
			status: http.StatusUnauthorized,
		},
		{
			name: "internal error",
			req: AuthLoginRequest{
				Username: "alice",
				Password: "pass",
			},
			mock: func(m *mockAuthService) {
				m.loginFn = func(ctx context.Context, username, password string) (user.LoginResult, error) {
					return user.LoginResult{}, errors.New("fail")
				}
			},
			status: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := &mockAuthService{}
			tt.mock(m)

			handler := authLogin(m, accessTTL, refreshTTL, false)

			rr := testRequest(t, handler, http.MethodPost, "/auth/login", tt.req)

			require.Equal(t, tt.status, rr.Code)

			if tt.check != nil {
				tt.check(t, rr)
			}
		})
	}
}

type mockAuthService struct {
	loginFn func(ctx context.Context, username, password string) (user.LoginResult, error)
}

func (m *mockAuthService) Login(ctx context.Context, username string, password string) (user.LoginResult, error) {
	return m.loginFn(ctx, username, password)
}
