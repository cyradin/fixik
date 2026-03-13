package router

import (
	"net/http"

	"github.com/cyradin/fixik/internal/container"
	"github.com/cyradin/fixik/internal/dict"
	"github.com/go-chi/chi/v5"
)

func impactRoutes(c *container.Container) func(r chi.Router) {
	impactManager := c.ImpactManager()

	return func(r chi.Router) {
		r.Get("/", listImpacts(impactManager))
		r.Post("/", createImpact(impactManager))
		r.Get("/{id}", getImpact(impactManager))
		r.Put("/{id}", updateImpact(impactManager))
		r.Delete("/{id}", deleteImpact(impactManager))
	}
}

// @Summary Create impact
// @Description Create new impact dictionary entry
// @Tags impacts
// @Accept json
// @Produce json
// @Param request body CreateDictEntityRequest true "Impact data"
// @Success 200 {object} DictEntity
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /impacts [post]
func createImpact(manager *dict.EntityManager) http.HandlerFunc {
	return createDictEntity(manager)
}

// @Summary Get impact by ID
// @Description Get impact dictionary entry by ID
// @Tags impacts
// @Accept json
// @Produce json
// @Param id path int true "Impact ID"
// @Success 200 {object} DictEntity
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /impacts/{id} [get]
func getImpact(manager *dict.EntityManager) http.HandlerFunc {
	return getDictEntity(manager)
}

// @Summary Update impact
// @Description Update impact dictionary entry by ID
// @Tags impacts
// @Accept json
// @Produce json
// @Param id path int true "Impact ID"
// @Param request body UpdateDictEntityRequest true "Impact data"
// @Success 200 {object} DictEntity
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /impacts/{id} [put]
func updateImpact(manager *dict.EntityManager) http.HandlerFunc {
	return updateDictEntity(manager)
}

// deleteImpact godoc
// @Summary Delete impact
// @Description Delete impact dictionary entry by ID
// @Tags impacts
// @Accept json
// @Produce json
// @Param id path int true "Impact ID"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /impacts/{id} [delete]
func deleteImpact(manager *dict.EntityManager) http.HandlerFunc {
	return deleteDictEntity(manager)
}

// listImpacts godoc
// @Summary List impacts
// @Description Get all impacts in dictionary
// @Tags impacts
// @Accept json
// @Produce json
// @Success 200 {object} ListDictEntitiesResponse
// @Failure 500 {object} ErrorResponse
// @Router /impacts [get]
func listImpacts(manager *dict.EntityManager) http.HandlerFunc {
	return listDictEntities(manager)
}
