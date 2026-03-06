package router

import (
	"github.com/cyradin/fixik/internal/config"
	"github.com/go-chi/chi/v5"
)

func New(_ *config.Config) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/users", createUser)
	r.Get("/users", listUsers)
	r.Get("/users/{id}", getUser)
	r.Patch("/users/{id}", updateUser)
	r.Delete("/users/{id}", deleteUser)

	r.Post("/incidents", createIncident)
	r.Get("/incidents", listIncidents)
	r.Get("/incidents/{id}", getIncident)
	r.Put("/incidents/{id}", updateIncident)
	r.Delete("/incidents/{id}", deleteIncident)

	r.Post("/incidents/{id}/comments", addComment)
	r.Get("/incidents/{id}/comments", listComments)

	r.Post("/incidents/{id}/assign", assignIncident)
	r.Post("/incidents/{id}/status", changeStatus)

	return r
}
