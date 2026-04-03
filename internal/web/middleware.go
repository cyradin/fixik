package web

import (
	"context"
	"net/http"

	"github.com/cyradin/fixik/internal/user"
)

type userAuthProvider interface {
	GetUserFromAccessToken(ctx context.Context, token string) (user.User, error)
}

func AuthMiddleware(userAuthProvider userAuthProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie(accessTokenCookie)
			if err != nil || c.Value == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			u, err := userAuthProvider.GetUserFromAccessToken(r.Context(), c.Value)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			ctx := user.WithContext(r.Context(), u)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
