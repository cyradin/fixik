package web

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cyradin/fixik/internal/container"
	"github.com/cyradin/fixik/internal/status"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type statusManager interface {
	Create(ctx context.Context, s status.Status) (status.Status, error)
	GetByID(ctx context.Context, id int64) (status.Status, error)
	List(ctx context.Context) ([]status.Status, error)
	Update(ctx context.Context, s status.Status) (status.Status, error)
	Delete(ctx context.Context, id int64) error
}

type CreateStatusRequest struct {
	Name        string `json:"name" validate:"required"`
	Code        string `json:"code" validate:"required"`
	Description string `json:"description" validate:"required"`
	Sort        int    `json:"sort" validate:"required"`
	IsFinal     bool   `json:"isFinal" validate:"required"`
}

func (r CreateStatusRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, maxNameLen)),
		validation.Field(&r.Code, validation.Required, validation.Length(1, maxCodeLen)),
		validation.Field(&r.Sort, validation.Required),
	)
}

type UpdateStatusRequest struct {
	Name        string `json:"name" validate:"required"`
	Code        string `json:"code" validate:"required"`
	Description string `json:"description" validate:"required"`
	Sort        int    `json:"sort" validate:"required"`
	IsFinal     bool   `json:"isFinal" validate:"required"`
}

func (r UpdateStatusRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, maxNameLen)),
		validation.Field(&r.Code, validation.Required, validation.Length(1, maxCodeLen)),
		validation.Field(&r.Sort, validation.Required),
	)
}

type ListStatusesResponse struct {
	Items []Status `json:"items" validate:"required"`
}

type Status struct {
	ID          int64  `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Code        string `json:"code" validate:"required"`
	Description string `json:"description" validate:"required"`
	Sort        int    `json:"sort" validate:"required"`
	IsFinal     bool   `json:"isFinal" validate:"required"`
}

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
// @Param request body CreateStatusRequest true "Status data"
// @Success 200 {object} Status
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /statuses [post]
func createStatus(manager statusManager) http.HandlerFunc {
	return handle(func(ctx context.Context, req CreateStatusRequest) (Status, error) {
		entity := status.Status{
			Name:        req.Name,
			Code:        req.Code,
			Description: req.Description,
			Sort:        req.Sort,
			IsFinal:     req.IsFinal,
		}

		result, err := manager.Create(ctx, entity)
		if err != nil {
			return Status{}, err
		}

		return toStatus(result), nil
	})
}

// @Summary Get status by ID
// @Description Get status dictionary entry by ID
// @Tags statuses
// @Accept json
// @Produce json
// @Param id path int true "Status ID"
// @Success 200 {object} Status
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /statuses/{id} [get]
func getStatus(manager statusManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			handleError(r.Context(), w, ErrRequestPathEntityID)
			return
		}

		handle(func(ctx context.Context, _ NoBody) (Status, error) {
			result, err := manager.GetByID(ctx, id)
			if err != nil {
				return Status{}, err
			}

			return toStatus(result), nil
		})(w, r)
	}
}

// @Summary Update status
// @Description Update status dictionary entry by ID
// @Tags statuses
// @Accept json
// @Produce json
// @Param id path int true "Status ID"
// @Param request body UpdateStatusRequest true "Status data"
// @Success 200 {object} Status
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /statuses/{id} [put]
func updateStatus(manager statusManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			handleError(r.Context(), w, ErrRequestPathEntityID)
			return
		}

		handle(func(ctx context.Context, req UpdateStatusRequest) (Status, error) {
			entity := status.Status{
				ID:          id,
				Name:        req.Name,
				Code:        req.Code,
				Description: req.Description,
				Sort:        req.Sort,
				IsFinal:     req.IsFinal,
			}

			result, err := manager.Update(ctx, entity)
			if err != nil {
				return Status{}, err
			}

			return toStatus(result), nil
		})(w, r)
	}
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
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /statuses/{id} [delete]
func deleteStatus(manager statusManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			handleError(r.Context(), w, ErrRequestPathEntityID)
			return
		}

		handle(func(ctx context.Context, _ NoBody) (NoBody, error) {
			if err := manager.Delete(ctx, id); err != nil {
				return NoBody{}, err
			}

			return NoBody{}, nil
		})(w, r)
	}
}

// listStatuses godoc
// @Summary List statuses
// @Description Get all statuses in dictionary
// @Tags statuses
// @Accept json
// @Produce json
// @Success 200 {object} ListStatusesResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /statuses [get]
func listStatuses(manager statusManager) http.HandlerFunc {
	return handle(func(ctx context.Context, _ NoBody) (ListStatusesResponse, error) {
		items, err := manager.List(ctx)
		if err != nil {
			return ListStatusesResponse{}, err
		}

		resp := make([]Status, 0, len(items))
		for _, item := range items {
			resp = append(resp, toStatus(item))
		}

		return ListStatusesResponse{Items: resp}, nil
	})
}

func toStatus(item status.Status) Status {
	return Status{
		ID:          item.ID,
		Name:        item.Name,
		Code:        item.Code,
		Description: item.Description,
		Sort:        item.Sort,
		IsFinal:     item.IsFinal,
	}
}
