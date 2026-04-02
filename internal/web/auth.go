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
	}
}

func authLogin(srv userLoginer, accessTTL, refreshTTL time.Duration, secureCookies bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handle(func(ctx context.Context, req AuthLoginRequest) (NoBody, error) {
			result, err := srv.Login(ctx, req.Username, req.Password)
			if err != nil {
				if errors.Is(err, user.ErrUserNotFound) {
					return NoBody{}, UnauthorizedError{Err: err}
				}

				return NoBody{}, err
			}

			setAccessToken(w, result.AccessToken, accessTTL, secureCookies)
			setRefreshToken(w, result.RefreshToken, refreshTTL, secureCookies)

			return NoBody{}, nil
		})(w, r)
	}
}

func authLogout(secureCookies bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setAccessToken(w, "", -1*time.Second, secureCookies)
		setRefreshToken(w, "", -1*time.Second, secureCookies)

		w.WriteHeader(http.StatusOK)
	}
}

func authRefresh(srv tokenRefresher, accessTTL, refreshTTL time.Duration, secureCookies bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handle(func(ctx context.Context, req NoBody) (NoBody, error) {
			rtCookie, err := r.Cookie("refresh_token")
			if err != nil {
				return NoBody{}, UnauthorizedError{Err: fmt.Errorf("refresh token missing")}
			}

			result, err := srv.Refresh(ctx, rtCookie.Value)
			if err != nil {
				return NoBody{}, UnauthorizedError{Err: err}
			}

			setAccessToken(w, result.AccessToken, accessTTL, secureCookies)
			setRefreshToken(w, result.RefreshToken, refreshTTL, secureCookies)

			return NoBody{}, nil
		})(w, r)
	}
}

func setAccessToken(w http.ResponseWriter, token string, ttl time.Duration, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
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
		Name:     "refresh_token",
		Value:    token,
		HttpOnly: true,
		Secure:   secure,
		Path:     "/api/auth/refresh",
		MaxAge:   int(ttl.Seconds()),
		SameSite: http.SameSiteLaxMode,
	})
}
