package router

import (
	"context"
	"net/http"

	"github.com/cyradin/fixik/internal/container"
	"github.com/cyradin/fixik/internal/user"
	"github.com/go-chi/chi/v5"
)

type ListRolesResponse struct {
	Items []Role `json:"items"`
}

type Role struct {
	Name        string `json:"name"`
	Code        string `json:"code" enums:"admin,manager,user"`
	Description string `json:"description"`
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
