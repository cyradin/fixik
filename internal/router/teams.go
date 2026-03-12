package router

import (
	"net/http"

	"github.com/cyradin/fixik/internal/container"
	"github.com/cyradin/fixik/internal/dict"
	"github.com/go-chi/chi/v5"
)

func teamRoutes(c *container.Container) func(r chi.Router) {
	teamManager := c.TeamManager()

	return func(r chi.Router) {
		r.Get("/", listTeams(teamManager))
		r.Post("/", createTeam(teamManager))
		r.Get("/{id}", getTeam(teamManager))
		r.Put("/{id}", updateTeam(teamManager))
		r.Delete("/{id}", deleteTeam(teamManager))
	}
}

// @Summary Create team
// @Description Create new team dictionary entry
// @Tags teams
// @Accept json
// @Produce json
// @Param request body CreateDictEntityRequest true "Team data"
// @Success 200 {object} DictEntityResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /teams [post]
func createTeam(manager *dict.EntityManager) http.HandlerFunc {
	return createDictEntity(manager)
}

// @Summary Get team by ID
// @Description Get team dictionary entry by ID
// @Tags teams
// @Accept json
// @Produce json
// @Param id path int true "Team ID"
// @Success 200 {object} DictEntityResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /teams/{id} [get]
func getTeam(manager *dict.EntityManager) http.HandlerFunc {
	return getDictEntity(manager)
}

// @Summary Update team
// @Description Update team dictionary entry by ID
// @Tags teams
// @Accept json
// @Produce json
// @Param id path int true "Team ID"
// @Param request body UpdateDictEntityRequest true "Team data"
// @Success 200 {object} DictEntityResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /teams/{id} [put]
func updateTeam(manager *dict.EntityManager) http.HandlerFunc {
	return updateDictEntity(manager)
}

// deleteTeam godoc
// @Summary Delete team
// @Description Delete team dictionary entry by ID
// @Tags teams
// @Accept json
// @Produce json
// @Param id path int true "Team ID"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /teams/{id} [delete]
func deleteTeam(manager *dict.EntityManager) http.HandlerFunc {
	return deleteDictEntity(manager)
}

// listTeams godoc
// @Summary List teams
// @Description Get all teams in dictionary
// @Tags teams
// @Accept json
// @Produce json
// @Success 200 {object} ListDictEntitiesResponse
// @Failure 500 {object} ErrorResponse
// @Router /teams [get]
func listTeams(manager *dict.EntityManager) http.HandlerFunc {
	return listDictEntities(manager)
}
