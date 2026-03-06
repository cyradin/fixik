package router

import (
	"net/http"

	"github.com/cyradin/fixik/internal/config"
	"github.com/go-chi/chi/v5"
)

func NewDebug(_ *config.Config) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/livez", livenessHandler)
	r.Get("/readyz", readinessHandler)

	return r
}

func livenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
