package web

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/cyradin/fixik/internal/container"
	"github.com/cyradin/fixik/internal/incident"
	"github.com/cyradin/fixik/internal/role"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type incidentManager interface {
	Create(ctx context.Context, i incident.CreateIncident) (incident.Incident, error)
	GetByID(ctx context.Context, id int64) (incident.Incident, error)
	Update(ctx context.Context, i incident.UpdateIncident) (incident.Incident, error)
	List(ctx context.Context, limit, offset int) (incident.IncidentList, error)
	Delete(ctx context.Context, id int64) error
}

type commentCreater interface {
	Create(ctx context.Context, incidentID int64, text string) (incident.Comment, error)
}

type commentByIncidentLister interface {
	ListByIncident(ctx context.Context, incidentID int64, limit, offset int) (incident.CommentList, error)
}

type CreateIncidentRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	StatusID    int64  `json:"statusId" validate:"required"`
	PriorityID  int64  `json:"priorityId" validate:"required"`
	TeamID      *int64 `json:"teamId"`
	UserID      *int64 `json:"userId"`
	AuthorID    *int64 `json:"authorId"`
}

func (r CreateIncidentRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Title, validation.Required),
		validation.Field(&r.Description, validation.Required),
		validation.Field(&r.StatusID, validation.Required, validation.Min(1)),
		validation.Field(&r.PriorityID, validation.Required, validation.Min(1)),
		validation.Field(&r.TeamID, validation.Min(1)),
		validation.Field(&r.UserID, validation.Min(1)),
		validation.Field(&r.AuthorID, validation.Min(1)),
	)
}

type UpdateIncidentRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	StatusID    *int64  `json:"statusId"`
	PriorityID  *int64  `json:"priorityId"`
	TeamID      *int64  `json:"teamId"`
	UserID      *int64  `json:"userId"`
	AuthorID    *int64  `json:"authorId"`
}

func (r UpdateIncidentRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Title),
		validation.Field(&r.Description),
		validation.Field(&r.StatusID, validation.Min(1)),
		validation.Field(&r.PriorityID, validation.Min(1)),
		validation.Field(&r.TeamID, validation.Min(1)),
		validation.Field(&r.UserID, validation.Min(1)),
		validation.Field(&r.AuthorID, validation.Min(1)),
	)
}

type IncidentResponse struct {
	ID            int64            `json:"id" validate:"required"`
	Title         string           `json:"title" validate:"required"`
	Description   string           `json:"description" validate:"required"`
	Status        DictEntityShort  `json:"status" validate:"required"`
	Priority      DictEntityShort  `json:"priority" validate:"required"`
	Team          *DictEntityShort `json:"team"`
	User          *UserResponse    `json:"user"`
	Author        *UserResponse    `json:"author"`
	CreatedAt     string           `json:"createdAt" validate:"required"`
	UpdatedAt     string           `json:"updatedAt" validate:"required"`
	CommentsCount int              `json:"commentsCount" validate:"required"`
}

type IncidentListResponse struct {
	Items      []IncidentResponse `json:"items" validate:"required"`
	Pagination Pagination         `json:"pagination" validate:"required"`
}

type Pagination struct {
	Limit  int `json:"limit" validate:"required"`
	Offset int `json:"offset" validate:"required"`
	Total  int `json:"total" validate:"required"`
}

type DictEntityShort struct {
	ID   int64  `json:"id" validate:"required"`
	Code string `json:"code" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type IncidentCommentCreateRequest struct {
	Text string `json:"text" validate:"required"`
}

func (r IncidentCommentCreateRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Text, validation.Required),
	)
}

type IncidentCommentListResponse struct {
	Items      []IncidentComment `json:"items" validate:"required"`
	Pagination Pagination        `json:"pagination" validate:"required"`
}

type IncidentComment struct {
	ID         int64        `json:"id" validate:"required"`
	IncidentID int64        `json:"incidentId" validate:"required"`
	Author     UserResponse `json:"author" validate:"required"`
	Text       string       `json:"text" validate:"required"`
	CreatedAt  string       `json:"createdAt" validate:"required"`
	UpdatedAt  string       `json:"updatedAt" validate:"required"`
}

func incidentRoutes(c *container.Container) func(r chi.Router) {
	incidentManager := c.IncidentManager()
	commentManager := c.CommentManager()

	return func(r chi.Router) {
		r.Post("/", createIncident(incidentManager))
		r.Get("/", listIncidents(incidentManager))
		r.Get("/{id}", getIncident(incidentManager))
		r.Patch("/{id}", updateIncident(incidentManager))
		r.Delete("/{id}", deleteIncident(incidentManager))
		r.Post("/{id}/comments", createIncidentComment(commentManager))
		r.Get("/{id}/comments", getIncidentComments(commentManager))
	}
}

// @Summary Create incident
// @Description Create new incident
// @Tags incidents
// @Accept json
// @Produce json
// @Param request body CreateIncidentRequest true "Incident data"
// @Success 200 {object} IncidentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /incidents [post]
func createIncident(manager incidentManager) http.HandlerFunc {
	return handle(func(ctx context.Context, req CreateIncidentRequest) (IncidentResponse, error) {
		if !checkPermissions(ctx, role.IncidentCreate) {
			return IncidentResponse{}, ErrForbidden
		}

		if err := req.Validate(); err != nil {
			return IncidentResponse{}, err
		}

		i, err := manager.Create(ctx, incident.CreateIncident{
			Title:       req.Title,
			Description: req.Description,
			StatusID:    req.StatusID,
			PriorityID:  req.PriorityID,
			TeamID:      req.TeamID,
			UserID:      req.UserID,
			AuthorID:    req.AuthorID,
		})
		if err != nil {
			return IncidentResponse{}, err
		}

		return toIncidentResponse(i), nil
	})
}

// @Summary Get incident
// @Description Get incident by ID
// @Tags incidents
// @Accept json
// @Produce json
// @Param id path int true "Incident ID"
// @Success 200 {object} IncidentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /incidents/{id} [get]
func getIncident(manager incidentManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			handleError(r.Context(), w, ErrRequestPathEntityID)
			return
		}

		handle(func(ctx context.Context, _ NoBody) (IncidentResponse, error) {
			if !checkPermissions(ctx, role.IncidentGet) {
				return IncidentResponse{}, ErrForbidden
			}

			i, err := manager.GetByID(ctx, id)
			if err != nil {
				return IncidentResponse{}, err
			}

			return toIncidentResponse(i), nil
		})(w, r)
	}
}

// @Summary Update incident
// @Description Update incident
// @Tags incidents
// @Accept json
// @Produce json
// @Param id path int true "Incident ID"
// @Param request body UpdateIncidentRequest true "Incident data"
// @Success 200 {object} IncidentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /incidents/{id} [patch]
func updateIncident(manager incidentManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			handleError(r.Context(), w, ErrRequestPathEntityID)
			return
		}

		handle(func(ctx context.Context, req UpdateIncidentRequest) (IncidentResponse, error) {
			if !checkPermissions(ctx, role.IncidentUpdate) {
				return IncidentResponse{}, ErrForbidden
			}

			if err := req.Validate(); err != nil {
				return IncidentResponse{}, err
			}

			i, err := manager.Update(ctx, incident.UpdateIncident{
				ID:          id,
				Title:       req.Title,
				Description: req.Description,
				StatusID:    req.StatusID,
				PriorityID:  req.PriorityID,
				TeamID:      req.TeamID,
				UserID:      req.UserID,
				AuthorID:    req.AuthorID,
			})
			if err != nil {
				return IncidentResponse{}, err
			}

			return toIncidentResponse(i), nil
		})(w, r)
	}
}

// @Summary Delete incident
// @Description Delete incident
// @Tags incidents
// @Accept json
// @Produce json
// @Param id path int true "Incident ID"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /incidents/{id} [delete]
func deleteIncident(manager incidentManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			handleError(r.Context(), w, ErrRequestPathEntityID)
			return
		}

		handle(func(ctx context.Context, _ NoBody) (NoBody, error) {
			if !checkPermissions(ctx, role.IncidentDelete) {
				return NoBody{}, ErrForbidden
			}

			return NoBody{}, manager.Delete(ctx, id)
		})(w, r)
	}
}

// @Summary List incidents
// @Description Get all incidents with pagination
// @Tags incidents
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(100)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} IncidentListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /incidents [get]
func listIncidents(manager incidentManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, offset, err := decodePagination(r, 1, 100) //nolint:mnd
		if err != nil {
			handleError(r.Context(), w, ErrValidation(err.Error()))
			return
		}

		handle(func(ctx context.Context, _ NoBody) (IncidentListResponse, error) {
			if !checkPermissions(ctx, role.IncidentGet) {
				return IncidentListResponse{}, ErrForbidden
			}

			listResult, err := manager.List(ctx, limit, offset)
			if err != nil {
				return IncidentListResponse{}, err
			}

			respItems := make([]IncidentResponse, len(listResult.Items))

			for i, item := range listResult.Items {
				respItems[i] = toIncidentResponse(item)
			}

			return IncidentListResponse{
				Items: respItems,
				Pagination: Pagination{
					Limit:  limit,
					Offset: offset,
					Total:  listResult.Total,
				},
			}, nil
		})(w, r)
	}
}

func toIncidentResponse(i incident.Incident) IncidentResponse {
	var (
		team   *DictEntityShort
		user   *UserResponse
		author *UserResponse
	)

	if i.Team != nil {
		team = &DictEntityShort{
			ID:   i.Team.ID,
			Code: i.Team.Code,
			Name: i.Team.Name,
		}
	}

	if i.User != nil {
		user = new(toUserResponse(*i.User))
	}

	if i.Author != nil {
		author = new(toUserResponse(*i.Author))
	}

	return IncidentResponse{
		ID:          i.ID,
		Title:       i.Title,
		Description: i.Description,
		Status: DictEntityShort{
			ID:   i.Status.ID,
			Code: i.Status.Code,
			Name: i.Status.Name,
		},
		Priority: DictEntityShort{
			ID:   i.Priority.ID,
			Code: i.Priority.Code,
			Name: i.Priority.Name,
		},
		Team:   team,
		User:   user,
		Author: author,

		CreatedAt: i.CreatedAt.Format(time.RFC3339),
		UpdatedAt: i.UpdatedAt.Format(time.RFC3339),

		CommentsCount: i.CommentsCount,
	}
}

// @Summary Create incident comment
// @Description Create a new comment for an incident
// @Tags incidents
// @Accept json
// @Produce json
// @Param id path int true "Incident ID"
// @Param request body IncidentCommentCreateRequest true "Comment data"
// @Success 200 {object} IncidentComment
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /incidents/{id}/comments [post]
func createIncidentComment(manager commentCreater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		incidentID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			handleError(r.Context(), w, ErrRequestPathEntityID)
			return
		}

		handle(func(ctx context.Context, req IncidentCommentCreateRequest) (IncidentComment, error) {
			if !checkPermissions(ctx, role.CommentCreate) {
				return IncidentComment{}, ErrForbidden
			}

			if err := req.Validate(); err != nil {
				return IncidentComment{}, err
			}

			c, err := manager.Create(ctx, incidentID, req.Text)
			if err != nil {
				return IncidentComment{}, err
			}

			return toIncidentCommentResponse(c), nil
		})(w, r)
	}
}

// @Summary Get incident comments
// @Description Get all comments for an incident with pagination
// @Tags incidents
// @Accept json
// @Produce json
// @Param id path int true "Incident ID"
// @Param limit query int false "Limit" default(100)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} IncidentCommentListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /incidents/{id}/comments [get]
func getIncidentComments(manager commentByIncidentLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		incidentID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			handleError(r.Context(), w, ErrRequestPathEntityID)
			return
		}

		limit, offset, err := decodePagination(r, 1, 100) //nolint:mnd
		if err != nil {
			handleError(r.Context(), w, ErrValidation(err.Error()))
			return
		}

		handle(func(ctx context.Context, _ NoBody) (IncidentCommentListResponse, error) {
			if !checkPermissions(ctx, role.CommentGet) {
				return IncidentCommentListResponse{}, ErrForbidden
			}

			result, err := manager.ListByIncident(ctx, incidentID, limit, offset)
			if err != nil {
				return IncidentCommentListResponse{}, err
			}

			resp := make([]IncidentComment, len(result.Items))
			for i, item := range result.Items {
				resp[i] = toIncidentCommentResponse(item)
			}

			return IncidentCommentListResponse{
				Items: resp,
				Pagination: Pagination{
					Limit:  limit,
					Offset: offset,
					Total:  result.Total,
				},
			}, nil
		})(w, r)
	}
}

func toIncidentCommentResponse(c incident.Comment) IncidentComment {
	return IncidentComment{
		ID:         c.ID,
		IncidentID: c.IncidentID,
		Author:     toUserResponse(c.Author),
		Text:       c.Text,
		CreatedAt:  c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  c.UpdatedAt.Format(time.RFC3339),
	}
}
