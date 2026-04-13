package web

import (
	"context"
	"net/http"

	"github.com/cyradin/fixik/internal/container"
	"github.com/cyradin/fixik/internal/role"
	"github.com/go-chi/chi/v5"
)

type ListRolesResponse struct {
	Items []Role `json:"items" validate:"required"`
}

type Role struct {
	Name        string       `json:"name" validate:"required"`
	Code        string       `json:"code" enums:"admin,manager,user" validate:"required"`
	Description string       `json:"description" validate:"required"`
	Permissions []Permission `json:"permissions" validate:"required"`
}

type Permission struct {
	Code string `json:"code" validate:"required"`
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
		roles := role.List()

		items := make([]Role, 0, len(roles))
		for _, role := range roles {
			items = append(items, Role{
				Name:        role.Name,
				Code:        role.Code,
				Description: role.Description,
				Permissions: toPermissionsResponse(role.Permissions),
			})
		}

		return ListRolesResponse{Items: items}, nil
	})
}

func toPermissionsResponse(permissions role.Permission) []Permission {
	codes := permissions.Codes()
	result := make([]Permission, 0, len(codes))

	for _, code := range codes {
		result = append(result, Permission{
			Code: code,
		})
	}

	return result
}
