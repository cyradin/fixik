package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cyradin/fixik/internal/container"
	"github.com/cyradin/fixik/internal/user"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type userLoginer interface {
	Login(ctx context.Context, username string, password string) (user.LoginResult, error)
}

type tokenRefresher interface {
	Refresh(ctx context.Context, refreshToken string) (user.LoginResult, error)
}

//nolint:gosec
type AuthLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (r AuthLoginRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Username, validation.Required),
		validation.Field(&r.Password, validation.Required),
	)
}

func authRoutes(c *container.Container) func(r chi.Router) {
	authService := c.AuthService()

	return func(r chi.Router) {
		r.Post("/login", authLogin(
			authService,
			c.Cfg().Auth.AccessTokenTTL,
			c.Cfg().Auth.RefreshTokenTTL,
			c.Cfg().Auth.SecureCookies,
		))

		r.Post("/logout", authLogout(
			c.Cfg().Auth.SecureCookies,
		))

		r.Post("/refresh", authRefresh(
			authService,
			c.Cfg().Auth.AccessTokenTTL,
			c.Cfg().Auth.RefreshTokenTTL,
			c.Cfg().Auth.SecureCookies,
		))

		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(c.AuthService()))
			r.Post("/password", changePassword(authService))
		})
	}
}

// @Summary Login
// @Description Login with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body AuthLoginRequest true "Login data"
// @Success 200 {object} UserResponse "Sets new access and refresh cookies"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
func authLogin(srv userLoginer, accessTTL, refreshTTL time.Duration, secureCookies bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handle(func(ctx context.Context, req AuthLoginRequest) (UserResponse, error) {
			result, err := srv.Login(ctx, req.Username, req.Password)
			if err != nil {
				if errors.Is(err, user.ErrUserNotFound) {
					return UserResponse{}, ErrInvalidLoginPassword
				}

				return UserResponse{}, err
			}

			setAccessToken(w, result.AccessToken, accessTTL, secureCookies)
			setRefreshToken(w, result.RefreshToken, refreshTTL, secureCookies)

			return toUserResponse(result.User), nil
		})(w, r)
	}
}

// @Summary Logout
// @Description Logout by clearing access and refresh cookies
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 "Clears cookies"
// @Failure 500 {object} ErrorResponse
// @Router /auth/logout [post]
func authLogout(secureCookies bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setAccessToken(w, "", -1*time.Second, secureCookies)
		setRefreshToken(w, "", -1*time.Second, secureCookies)

		w.WriteHeader(http.StatusOK)
	}
}

// @Summary Refresh
// @Description Refresh access and refresh tokens using refresh cookie
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} UserResponse "Sets new access and refresh cookies"
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/refresh [post]
func authRefresh(srv tokenRefresher, accessTTL, refreshTTL time.Duration, secureCookies bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handle(func(ctx context.Context, req NoBody) (UserResponse, error) {
			rtCookie, err := r.Cookie(refreshTokenCookie)
			if err != nil {
				return UserResponse{}, &UserMessageError{Err: fmt.Errorf("refresh token missing"), Status: http.StatusUnauthorized}
			}

			result, err := srv.Refresh(ctx, rtCookie.Value)
			if err != nil {
				return UserResponse{}, &UserMessageError{Err: err, Status: http.StatusUnauthorized}
			}

			setAccessToken(w, result.AccessToken, accessTTL, secureCookies)
			setRefreshToken(w, result.RefreshToken, refreshTTL, secureCookies)

			return toUserResponse(result.User), nil
		})(w, r)
	}
}

// @Summary Change user password
// @Description Change authorized user's password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ChangePasswordRequest true "Password data"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/password [post]
func changePassword(passwordChanger passwordChanger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handle(func(ctx context.Context, req ChangePasswordRequest) (NoBody, error) {
			u, ok := user.FromContext(r.Context())
			if !ok {
				return NoBody{}, ErrUnauthorized
			}

			if err := req.Validate(); err != nil {
				return NoBody{}, err
			}

			if err := passwordChanger.ChangePassword(ctx, u.ID, req.CurrentPassword, req.NewPassword); err != nil {
				if errors.Is(err, user.ErrUserNotFound) || errors.Is(err, user.ErrUserInvalidPassword) {
					return NoBody{}, ErrInvalidPassword
				}

				return NoBody{}, err
			}

			return NoBody{}, nil
		})(w, r)
	}
}

func setAccessToken(w http.ResponseWriter, token string, ttl time.Duration, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     accessTokenCookie,
		Value:    token,
		HttpOnly: true,
		Secure:   secure,
		Path:     "/",
		MaxAge:   int(ttl.Seconds()),
		SameSite: http.SameSiteLaxMode,
	})
}

func setRefreshToken(w http.ResponseWriter, token string, ttl time.Duration, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenCookie,
		Value:    token,
		HttpOnly: true,
		Secure:   secure,
		Path:     "/api/auth/refresh",
		MaxAge:   int(ttl.Seconds()),
		SameSite: http.SameSiteLaxMode,
	})
}

const (
	accessTokenCookie  = "access_token"
	refreshTokenCookie = "refresh_token"
)
