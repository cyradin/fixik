package router

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cyradin/fixik/internal/container"
	"github.com/cyradin/fixik/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func New(c *container.Container, allowedOriginsCORS []string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOriginsCORS,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, //nolint:mnd
	}))

	r.Route("/api", func(r chi.Router) {
		r.Post("/users", createUser)
		r.Get("/users", listUsers)
		r.Get("/users/{id}", getUser)
		r.Patch("/users/{id}", updateUser)
		r.Delete("/users/{id}", deleteUser)

		r.Route("/statuses", statusRoutes(c))

		r.Post("/incidents", createIncident)
		r.Get("/incidents", listIncidents)
		r.Get("/incidents/{id}", getIncident)
		r.Put("/incidents/{id}", updateIncident)
		r.Delete("/incidents/{id}", deleteIncident)

		r.Post("/incidents/{id}/comments", addComment)
		r.Get("/incidents/{id}/comments", listComments)

		r.Post("/incidents/{id}/assign", assignIncident)
		r.Post("/incidents/{id}/status", changeStatus)
	})

	return r
}

func decodeJSON(r *http.Request, v any) error {
	defer func() { _ = r.Body.Close() }()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	return dec.Decode(v)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, ErrorResponse{
		Error: err.Error(),
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(v)
}

func handle[Req any, Resp any](fn func(ctx context.Context, req Req) (Resp, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Req

		if _, ok := any(req).(NoBody); !ok {
			if err := decodeJSON(r, &req); err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
		}

		// optional validation
		if v, ok := any(req).(Validatable); ok {
			if err := v.Validate(); err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
		}

		ctx := r.Context()

		resp, err := fn(ctx, req)
		if err != nil {
			logger.FromContext(ctx).Error("request error", logger.Error(err))
			writeError(w, http.StatusInternalServerError, err)

			return
		}

		if _, ok := any(resp).(NoBody); !ok {
			writeJSON(w, http.StatusOK, resp)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

type Validatable interface {
	Validate() error
}

type NoBody struct{}
