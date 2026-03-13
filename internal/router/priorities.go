package router

import (
	"net/http"

	"github.com/cyradin/fixik/internal/container"
	"github.com/cyradin/fixik/internal/dict"
	"github.com/go-chi/chi/v5"
)

func priorityRoutes(c *container.Container) func(r chi.Router) {
	priorityManager := c.PriorityManager()

	return func(r chi.Router) {
		r.Get("/", listPriorities(priorityManager))
		r.Post("/", createPriority(priorityManager))
		r.Get("/{id}", getPriority(priorityManager))
		r.Put("/{id}", updatePriority(priorityManager))
		r.Delete("/{id}", deletePriority(priorityManager))
	}
}

// @Summary Create status
// @Description Create new status dictionary entry
// @Tags priorities
// @Accept json
// @Produce json
// @Param request body CreateDictEntityRequest true "Priority data"
// @Success 200 {object} DictEntity
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /priorities [post]
func createPriority(manager *dict.EntityManager) http.HandlerFunc {
	return createDictEntity(manager)
}

// @Summary Get status by ID
// @Description Get status dictionary entry by ID
// @Tags priorities
// @Accept json
// @Produce json
// @Param id path int true "Priority ID"
// @Success 200 {object} DictEntity
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /priorities/{id} [get]
func getPriority(manager *dict.EntityManager) http.HandlerFunc {
	return getDictEntity(manager)
}

// @Summary Update status
// @Description Update status dictionary entry by ID
// @Tags priorities
// @Accept json
// @Produce json
// @Param id path int true "Priority ID"
// @Param request body UpdateDictEntityRequest true "Priority data"
// @Success 200 {object} DictEntity
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /priorities/{id} [put]
func updatePriority(manager *dict.EntityManager) http.HandlerFunc {
	return updateDictEntity(manager)
}

// deletePriority godoc
// @Summary Delete status
// @Description Delete status dictionary entry by ID
// @Tags priorities
// @Accept json
// @Produce json
// @Param id path int true "Priority ID"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /priorities/{id} [delete]
func deletePriority(manager *dict.EntityManager) http.HandlerFunc {
	return deleteDictEntity(manager)
}

// listPriorities godoc
// @Summary List priorities
// @Description Get all priorities in dictionary
// @Tags priorities
// @Accept json
// @Produce json
// @Success 200 {object} ListDictEntitiesResponse
// @Failure 500 {object} ErrorResponse
// @Router /priorities [get]
func listPriorities(manager *dict.EntityManager) http.HandlerFunc {
	return listDictEntities(manager)
}
