package web

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/cyradin/fixik/internal/container"
	"github.com/cyradin/fixik/internal/role"
	"github.com/cyradin/fixik/internal/team"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type teamManager interface {
	Create(ctx context.Context, e team.Team) (team.Team, error)
	GetByID(ctx context.Context, id int64) (team.Team, error)
	List(ctx context.Context) ([]team.Team, error)
	Update(ctx context.Context, s team.Team) (team.Team, error)
	Delete(ctx context.Context, id int64) error
}

type CreateTeamRequest struct {
	Name        string `json:"name" validate:"required"`
	Code        string `json:"code" validate:"required"`
	Description string `json:"description" validate:"required"`
	Sort        int    `json:"sort" validate:"required"`
}

func (r CreateTeamRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, maxNameLen)),
		validation.Field(
			&r.Code,
			validation.Required,
			validation.Length(1, maxCodeLen),
			validation.Match(codeRegexp),
		),
		validation.Field(&r.Sort, validation.Required),
	)
}

type UpdateTeamRequest struct {
	Name        string `json:"name" validate:"required"`
	Code        string `json:"code" validate:"required"`
	Description string `json:"description" validate:"required"`
	Sort        int    `json:"sort" validate:"required"`
}

func (r UpdateTeamRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, maxNameLen)),
		validation.Field(
			&r.Code,
			validation.Required,
			validation.Length(1, maxCodeLen),
			validation.Match(codeRegexp),
		),
		validation.Field(&r.Sort, validation.Required),
	)
}

type ListTeamsResponse struct {
	Items []Team `json:"items" validate:"required"`
}

type Team struct {
	ID          int64  `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Code        string `json:"code" validate:"required"`
	Description string `json:"description" validate:"required"`
	Sort        int    `json:"sort" validate:"required"`
}

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
// @Param request body CreateTeamRequest true "Team data"
// @Success 200 {object} Team
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /teams [post]
func createTeam(manager teamManager) http.HandlerFunc {
	return handle(func(ctx context.Context, req CreateTeamRequest) (Team, error) {
		if !checkPermissions(ctx, role.TeamCreate) {
			return Team{}, ErrForbidden
		}

		entity := team.Team{
			Name:        req.Name,
			Code:        req.Code,
			Description: req.Description,
			Sort:        req.Sort,
		}

		result, err := manager.Create(ctx, entity)
		if err != nil {
			return Team{}, err
		}

		return toTeamResponse(result), nil
	})
}

// @Summary Get team by ID
// @Description Get team dictionary entry by ID
// @Tags teams
// @Accept json
// @Produce json
// @Param id path int true "Team ID"
// @Success 200 {object} Team
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /teams/{id} [get]
func getTeam(manager teamManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			handleError(r.Context(), w, ErrRequestPathEntityID)
			return
		}

		handle(func(ctx context.Context, _ NoBody) (Team, error) {
			if !checkPermissions(ctx, role.TeamGet) {
				return Team{}, ErrForbidden
			}

			result, err := manager.GetByID(ctx, id)
			if err != nil {
				return Team{}, err
			}

			return toTeamResponse(result), nil
		})(w, r)
	}
}

// @Summary Update team
// @Description Update team dictionary entry by ID
// @Tags teams
// @Accept json
// @Produce json
// @Param id path int true "Team ID"
// @Param request body UpdateTeamRequest true "Team data"
// @Success 200 {object} Team
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /teams/{id} [put]
func updateTeam(manager teamManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			handleError(r.Context(), w, ErrRequestPathEntityID)
			return
		}

		handle(func(ctx context.Context, req UpdateTeamRequest) (Team, error) {
			if !checkPermissions(ctx, role.TeamUpdate) {
				return Team{}, ErrForbidden
			}

			entity := team.Team{
				ID:          id,
				Name:        req.Name,
				Code:        req.Code,
				Description: req.Description,
				Sort:        req.Sort,
			}

			result, err := manager.Update(ctx, entity)
			if err != nil {
				return Team{}, err
			}

			return toTeamResponse(result), nil
		})(w, r)
	}
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
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /teams/{id} [delete]
func deleteTeam(manager teamManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			handleError(r.Context(), w, ErrRequestPathEntityID)
			return
		}

		handle(func(ctx context.Context, _ NoBody) (NoBody, error) {
			if !checkPermissions(ctx, role.TeamDelete) {
				return NoBody{}, ErrForbidden
			}

			if err := manager.Delete(ctx, id); err != nil {
				if errors.Is(err, team.ErrHasDependantIncidents) {
					return NoBody{}, ErrUnableToDelete("есть инциденты, назначенные на эту команду")
				} else if errors.Is(err, team.ErrHasDependantUsers) {
					return NoBody{}, ErrUnableToDelete("есть пользователи, привязанные к этой команде")
				}

				return NoBody{}, err
			}

			return NoBody{}, nil
		})(w, r)
	}
}

// listTeams godoc
// @Summary List teams
// @Description Get all teams in dictionary
// @Tags teams
// @Accept json
// @Produce json
// @Success 200 {object} ListTeamsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /teams [get]
func listTeams(manager teamManager) http.HandlerFunc {
	return handle(func(ctx context.Context, _ NoBody) (ListTeamsResponse, error) {
		if !checkPermissions(ctx, role.TeamGet) {
			return ListTeamsResponse{}, ErrForbidden
		}

		items, err := manager.List(ctx)
		if err != nil {
			return ListTeamsResponse{}, err
		}

		resp := make([]Team, 0, len(items))
		for _, item := range items {
			resp = append(resp, toTeamResponse(item))
		}

		return ListTeamsResponse{Items: resp}, nil
	})
}

func toTeamResponse(item team.Team) Team {
	return Team{
		ID:          item.ID,
		Name:        item.Name,
		Code:        item.Code,
		Description: item.Description,
		Sort:        item.Sort,
	}
}
