package router

import (
	"net/http"

	"github.com/cyradin/fixik/internal/container"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDebug(container *container.Container) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/livez", livenessHandler)
	r.Get("/readyz", readinessHandler(container.PgPool()))

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
