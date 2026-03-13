package router

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cyradin/fixik/internal/container"
	"github.com/cyradin/fixik/internal/incident"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type incidentManager interface {
	Create(ctx context.Context, i incident.CreateIncident) (incident.Incident, error)
	GetByID(ctx context.Context, id int64) (incident.Incident, error)
	Update(ctx context.Context, i incident.UpdateIncident) (incident.Incident, error)
	List(ctx context.Context, limit, offset int) ([]incident.Incident, error)
	Delete(ctx context.Context, id int64) error
}

type CreateIncidentRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImpactID    int64  `json:"impactId"`
	StatusID    int64  `json:"statusId"`
	PriorityID  int64  `json:"priorityId"`
}

func (r CreateIncidentRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Title, validation.Required),
		validation.Field(&r.Description, validation.Required),
		validation.Field(&r.ImpactID, validation.Required, validation.Min(1)),
		validation.Field(&r.StatusID, validation.Required, validation.Min(1)),
		validation.Field(&r.PriorityID, validation.Required, validation.Min(1)),
	)
}

type UpdateIncidentRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	ImpactID    *int64  `json:"impactId"`
	StatusID    *int64  `json:"statusId"`
	PriorityID  *int64  `json:"priorityId"`
}

func (r UpdateIncidentRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Title),
		validation.Field(&r.Description),
		validation.Field(&r.ImpactID, validation.Min(1)),
		validation.Field(&r.StatusID, validation.Min(1)),
		validation.Field(&r.PriorityID, validation.Min(1)),
	)
}

type IncidentResponse struct {
	ID          int64           `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Impact      DictEntityShort `json:"impact"`
	Status      DictEntityShort `json:"status"`
	Priority    DictEntityShort `json:"priority"`
}

type DictEntityShort struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func incidentRoutes(c *container.Container) func(r chi.Router) {
	incidentManager := c.IncidentManager()

	return func(r chi.Router) {
		r.Post("/incidents", createIncident(incidentManager))
		r.Get("/incidents", listIncidents(incidentManager))
		r.Get("/incidents/{id}", getIncident(incidentManager))
		r.Patch("/incidents/{id}", updateIncident(incidentManager))
		r.Delete("/incidents/{id}", deleteIncident(incidentManager))
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
// @Failure 500 {object} ErrorResponse
// @Router /incidents [post]
func createIncident(manager incidentManager) http.HandlerFunc {
	return handle(func(ctx context.Context, req CreateIncidentRequest) (IncidentResponse, error) {
		if err := req.Validate(); err != nil {
			return IncidentResponse{}, err
		}

		i, err := manager.Create(ctx, incident.CreateIncident{
			Title:       req.Title,
			Description: req.Description,
			ImpactID:    req.ImpactID,
			StatusID:    req.StatusID,
			PriorityID:  req.PriorityID,
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
// @Failure 500 {object} ErrorResponse
// @Router /incidents/{id} [get]
func getIncident(manager incidentManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		handle(func(ctx context.Context, _ NoBody) (IncidentResponse, error) {
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
// @Failure 500 {object} ErrorResponse
// @Router /incidents/{id} [patch]
func updateIncident(manager incidentManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		handle(func(ctx context.Context, req UpdateIncidentRequest) (IncidentResponse, error) {
			if err := req.Validate(); err != nil {
				return IncidentResponse{}, err
			}

			i, err := manager.Update(ctx, incident.UpdateIncident{
				ID:          id,
				Title:       req.Title,
				Description: req.Description,
				ImpactID:    req.ImpactID,
				StatusID:    req.StatusID,
				PriorityID:  req.PriorityID,
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
// @Failure 500 {object} ErrorResponse
// @Router /incidents/{id} [delete]
func deleteIncident(manager incidentManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		handle(func(ctx context.Context, _ NoBody) (NoBody, error) {
			return NoBody{}, manager.Delete(ctx, id)
		})(w, r)
	}
}

func listIncidents(manager incidentManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, offset, err := decodePagination(r, 1, 100)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		handle(func(ctx context.Context, _ NoBody) ([]IncidentResponse, error) {
			items, err := manager.List(ctx, limit, offset)
			if err != nil {
				return nil, err
			}

			resp := make([]IncidentResponse, len(items))

			for i, item := range items {
				resp[i] = toIncidentResponse(item)
			}

			return resp, nil

		})(w, r)
	}
}

func toIncidentResponse(i incident.Incident) IncidentResponse {
	return IncidentResponse{
		ID:          i.ID,
		Title:       i.Title,
		Description: i.Description,
		Impact: DictEntityShort{
			ID:   i.Impact.ID,
			Code: i.Impact.Code,
			Name: i.Impact.Name,
		},
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
	}
}
