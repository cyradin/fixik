package web

import (
	"net/http"

	"github.com/cyradin/fixik/internal/container"
	"github.com/cyradin/fixik/internal/dict"
	"github.com/go-chi/chi/v5"
)

func statusRoutes(c *container.Container) func(r chi.Router) {
	statusManager := c.StatusManager()

	return func(r chi.Router) {
		r.Get("/", listStatuses(statusManager))
		r.Post("/", createStatus(statusManager))
		r.Get("/{id}", getStatus(statusManager))
		r.Put("/{id}", updateStatus(statusManager))
		r.Delete("/{id}", deleteStatus(statusManager))
	}
}

// @Summary Create status
// @Description Create new status dictionary entry
// @Tags statuses
// @Accept json
// @Produce json
// @Param request body CreateDictEntityRequest true "Status data"
// @Success 200 {object} DictEntity
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /statuses [post]
func createStatus(manager *dict.EntityManager) http.HandlerFunc {
	return createDictEntity(manager)
}

// @Summary Get status by ID
// @Description Get status dictionary entry by ID
// @Tags statuses
// @Accept json
// @Produce json
// @Param id path int true "Status ID"
// @Success 200 {object} DictEntity
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /statuses/{id} [get]
func getStatus(manager *dict.EntityManager) http.HandlerFunc {
	return getDictEntity(manager)
}

// @Summary Update status
// @Description Update status dictionary entry by ID
// @Tags statuses
// @Accept json
// @Produce json
// @Param id path int true "Status ID"
// @Param request body UpdateDictEntityRequest true "Status data"
// @Success 200 {object} DictEntity
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /statuses/{id} [put]
func updateStatus(manager *dict.EntityManager) http.HandlerFunc {
	return updateDictEntity(manager)
}

// deleteStatus godoc
// @Summary Delete status
// @Description Delete status dictionary entry by ID
// @Tags statuses
// @Accept json
// @Produce json
// @Param id path int true "Status ID"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /statuses/{id} [delete]
func deleteStatus(manager *dict.EntityManager) http.HandlerFunc {
	return deleteDictEntity(manager)
}

// listStatuses godoc
// @Summary List statuses
// @Description Get all statuses in dictionary
// @Tags statuses
// @Accept json
// @Produce json
// @Success 200 {object} ListDictEntitiesResponse
// @Failure 500 {object} ErrorResponse
// @Router /statuses [get]
func listStatuses(manager *dict.EntityManager) http.HandlerFunc {
	return listDictEntities(manager)
}
