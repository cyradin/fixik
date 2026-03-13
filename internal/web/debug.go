package web

import (
	"net/http"

	_ "github.com/cyradin/fixik/docs"
	"github.com/cyradin/fixik/internal/container"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewDebugRouter(container *container.Container) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/livez", livenessHandler)
	r.Get("/readyz", readinessHandler(container.PgPool()))
	r.Get("/docs/*", httpSwagger.WrapHandler)

	return r
}

func livenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(pgPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := pgPool.Ping(r.Context()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
