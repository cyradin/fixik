package web

import (
	"context"
	"errors"
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

			http.SetCookie(w, &http.Cookie{
				Name:     "access_token",
				Value:    result.AccessToken,
				HttpOnly: true,
				Secure:   secureCookies,
				Path:     "/",
				MaxAge:   int(accessTTL.Seconds()),
				SameSite: http.SameSiteLaxMode,
			})

			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    result.RefreshToken,
				HttpOnly: true,
				Secure:   secureCookies,
				Path:     "/api/auth/refresh",
				MaxAge:   int(refreshTTL.Seconds()),
				SameSite: http.SameSiteLaxMode,
			})

			return NoBody{}, nil
		})(w, r)
	}
}
