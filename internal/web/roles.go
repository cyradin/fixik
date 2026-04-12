package web

import (
	"context"
	"net/http"

	"github.com/cyradin/fixik/internal/container"
	"github.com/cyradin/fixik/internal/user"
	"github.com/go-chi/chi/v5"
)

type ListRolesResponse struct {
	Items []Role `json:"items" validate:"required"`
}

type Role struct {
	Name        string `json:"name" validate:"required"`
	Code        string `json:"code" enums:"admin,manager,user" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func roleRoutes(_ *container.Container) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", listRoles())
	}
}

// listRoles godoc
// @Summary List roles
// @Description Get all user roles
// @Tags roles
// @Accept json
// @Produce json
// @Success 200 {object} ListRolesResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /roles [get]
func listRoles() http.HandlerFunc {
	return handle(func(ctx context.Context, _ NoBody) (ListRolesResponse, error) {
		roles := user.Roles()

		items := make([]Role, 0, len(roles))
		for _, role := range roles {
			items = append(items, Role{
				Name:        role.Name,
				Code:        role.Code,
				Description: role.Description,
			})
		}

		return ListRolesResponse{Items: items}, nil
	})
}
