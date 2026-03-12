package router

import (
	"net/http"

	"github.com/cyradin/fixik/internal/container"
	"github.com/cyradin/fixik/internal/dict"
	"github.com/go-chi/chi/v5"
)

func roleRoutes(c *container.Container) func(r chi.Router) {
	roleManager := c.RoleManager()

	return func(r chi.Router) {
		r.Get("/", listRoles(roleManager))
		r.Post("/", createRole(roleManager))
		r.Get("/{id}", getRole(roleManager))
		r.Put("/{id}", updateRole(roleManager))
		r.Delete("/{id}", deleteRole(roleManager))
	}
}

// @Summary Create role
// @Description Create new role dictionary entry
// @Tags roles
// @Accept json
// @Produce json
// @Param request body CreateDictEntityRequest true "Role data"
// @Success 200 {object} DictEntityResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /roles [post]
func createRole(manager *dict.EntityManager) http.HandlerFunc {
	return createDictEntity(manager)
}

// @Summary Get role by ID
// @Description Get role dictionary entry by ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} DictEntityResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /roles/{id} [get]
func getRole(manager *dict.EntityManager) http.HandlerFunc {
	return getDictEntity(manager)
}

// @Summary Update role
// @Description Update role dictionary entry by ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param request body UpdateDictEntityRequest true "Role data"
// @Success 200 {object} DictEntityResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /roles/{id} [put]
func updateRole(manager *dict.EntityManager) http.HandlerFunc {
	return updateDictEntity(manager)
}

// deleteRole godoc
// @Summary Delete role
// @Description Delete role dictionary entry by ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /roles/{id} [delete]
func deleteRole(manager *dict.EntityManager) http.HandlerFunc {
	return deleteDictEntity(manager)
}

// listRoles godoc
// @Summary List roles
// @Description Get all roles in dictionary
// @Tags roles
// @Accept json
// @Produce json
// @Success 200 {object} ListDictEntitiesResponse
// @Failure 500 {object} ErrorResponse
// @Router /roles [get]
func listRoles(manager *dict.EntityManager) http.HandlerFunc {
	return listDictEntities(manager)
}
