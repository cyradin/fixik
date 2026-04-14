package web

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/cyradin/fixik/internal/container"
	"github.com/cyradin/fixik/internal/priority"
	"github.com/cyradin/fixik/internal/role"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	maxNameLen = 100
	maxCodeLen = 50
)

type priorityManager interface {
	Create(ctx context.Context, e priority.Priority) (priority.Priority, error)
	GetByID(ctx context.Context, id int64) (priority.Priority, error)
	List(ctx context.Context) ([]priority.Priority, error)
	Update(ctx context.Context, s priority.Priority) (priority.Priority, error)
	Delete(ctx context.Context, id int64) error
}

type CreatePriorityRequest struct {
	Name        string `json:"name" validate:"required"`
	Code        string `json:"code" validate:"required"`
	Description string `json:"description" validate:"required"`
	Sort        int    `json:"sort" validate:"required"`
}

func (r CreatePriorityRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, maxNameLen)),
		validation.Field(&r.Code, validation.Required, validation.Length(1, maxCodeLen)),
		validation.Field(&r.Sort, validation.Required),
	)
}

type UpdatePriorityRequest struct {
	Name        string `json:"name" validate:"required"`
	Code        string `json:"code" validate:"required"`
	Description string `json:"description" validate:"required"`
	Sort        int    `json:"sort" validate:"required"`
}

func (r UpdatePriorityRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, maxNameLen)),
		validation.Field(&r.Code, validation.Required, validation.Length(1, maxCodeLen)),
		validation.Field(&r.Sort, validation.Required),
	)
}

type ListPrioritiesResponse struct {
	Items []Priority `json:"items" validate:"required"`
}

type Priority struct {
	ID          int64  `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Code        string `json:"code" validate:"required"`
	Description string `json:"description" validate:"required"`
	Sort        int    `json:"sort" validate:"required"`
}

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
// @Param request body CreatePriorityRequest true "Priority data"
// @Success 200 {object} Priority
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /priorities [post]
func createPriority(manager priorityManager) http.HandlerFunc {
	return handle(func(ctx context.Context, req CreatePriorityRequest) (Priority, error) {
		if !checkPermissions(ctx, role.PriorityCreate) {
			return Priority{}, ErrForbidden
		}

		entity := priority.Priority{
			Name:        req.Name,
			Code:        req.Code,
			Description: req.Description,
			Sort:        req.Sort,
		}

		result, err := manager.Create(ctx, entity)
		if err != nil {
			return Priority{}, err
		}

		return toPriorityResponse(result), nil
	})
}

// @Summary Get priority by ID
// @Description Get priority dictionary entry by ID
// @Tags priorities
// @Accept json
// @Produce json
// @Param id path int true "Priority ID"
// @Success 200 {object} Priority
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /priorities/{id} [get]
func getPriority(manager priorityManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			handleError(r.Context(), w, ErrRequestPathEntityID)
			return
		}

		handle(func(ctx context.Context, _ NoBody) (Priority, error) {
			if !checkPermissions(ctx, role.PriorityGet) {
				return Priority{}, ErrForbidden
			}

			result, err := manager.GetByID(ctx, id)
			if err != nil {
				return Priority{}, err
			}

			return toPriorityResponse(result), nil
		})(w, r)
	}
}

// @Summary Update priority
// @Description Update priority dictionary entry by ID
// @Tags priorities
// @Accept json
// @Produce json
// @Param id path int true "Priority ID"
// @Param request body UpdatePriorityRequest true "Priority data"
// @Success 200 {object} Priority
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /priorities/{id} [put]
func updatePriority(manager priorityManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			handleError(r.Context(), w, ErrRequestPathEntityID)
			return
		}

		handle(func(ctx context.Context, req UpdatePriorityRequest) (Priority, error) {
			if !checkPermissions(ctx, role.PriorityUpdate) {
				return Priority{}, ErrForbidden
			}

			entity := priority.Priority{
				ID:          id,
				Name:        req.Name,
				Code:        req.Code,
				Description: req.Description,
				Sort:        req.Sort,
			}

			result, err := manager.Update(ctx, entity)
			if err != nil {
				return Priority{}, err
			}

			return toPriorityResponse(result), nil
		})(w, r)
	}
}

// deletePriority godoc
// @Summary Delete priority
// @Description Delete priority dictionary entry by ID
// @Tags priorities
// @Accept json
// @Produce json
// @Param id path int true "Priority ID"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /priorities/{id} [delete]
func deletePriority(manager priorityManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			handleError(r.Context(), w, ErrRequestPathEntityID)
			return
		}

		handle(func(ctx context.Context, _ NoBody) (NoBody, error) {
			if !checkPermissions(ctx, role.PriorityDelete) {
				return NoBody{}, ErrForbidden
			}

			if err := manager.Delete(ctx, id); err != nil {
				if errors.Is(err, priority.ErrHasDependantEntities) {
					return NoBody{}, ErrUnableToDelete("есть инциденты с таким приоритетом")
				}

				return NoBody{}, err
			}

			return NoBody{}, nil
		})(w, r)
	}
}

// listPriorities godoc
// @Summary List priorities
// @Description Get all priorities in dictionary
// @Tags priorities
// @Accept json
// @Produce json
// @Success 200 {object} ListPrioritiesResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /priorities [get]
func listPriorities(manager priorityManager) http.HandlerFunc {
	return handle(func(ctx context.Context, _ NoBody) (ListPrioritiesResponse, error) {
		if !checkPermissions(ctx, role.PriorityGet) {
			return ListPrioritiesResponse{}, ErrForbidden
		}

		items, err := manager.List(ctx)
		if err != nil {
			return ListPrioritiesResponse{}, err
		}

		resp := make([]Priority, 0, len(items))
		for _, item := range items {
			resp = append(resp, toPriorityResponse(item))
		}

		return ListPrioritiesResponse{Items: resp}, nil
	})
}

func toPriorityResponse(item priority.Priority) Priority {
	return Priority{
		ID:          item.ID,
		Name:        item.Name,
		Code:        item.Code,
		Description: item.Description,
		Sort:        item.Sort,
	}
}
