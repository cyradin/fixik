package incident

import (
	"context"
	"fmt"

	"github.com/cyradin/fixik/internal/db"
	"github.com/cyradin/fixik/internal/dict"
)

type incidentRepo interface {
	Create(ctx context.Context, i *db.Incident) error
	GetByID(ctx context.Context, id int64) (db.Incident, error)
	Update(ctx context.Context, i *db.Incident) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]db.Incident, error)
}

type entityProvider interface {
	GetByID(ctx context.Context, id dict.EntityID) (dict.Entity, error)
	List(ctx context.Context) ([]dict.Entity, error)
}

type IncidentManager struct {
	repo             incidentRepo
	impactProvider   entityProvider
	statusProvider   entityProvider
	priorityProvider entityProvider
}

func NewIncidentManager(
	repo incidentRepo,
	impactProvider entityProvider,
	statusProvider entityProvider,
	priorityProvider entityProvider,
) *IncidentManager {
	return &IncidentManager{
		repo:             repo,
		impactProvider:   impactProvider,
		statusProvider:   statusProvider,
		priorityProvider: priorityProvider,
	}
}

func (m *IncidentManager) Create(ctx context.Context, cmd CreateIncident) (Incident, error) {
	dbIncident := db.Incident{
		Title:       cmd.Title,
		Description: cmd.Description,
		ImpactID:    cmd.ImpactID,
		StatusID:    cmd.StatusID,
		PriorityID:  cmd.PriorityID,
	}

	if err := m.repo.Create(ctx, &dbIncident); err != nil {
		return Incident{}, fmt.Errorf("create incident: %w", err)
	}

	return m.GetByID(ctx, dbIncident.ID)
}

func (m *IncidentManager) GetByID(ctx context.Context, id int64) (Incident, error) {
	result, err := m.repo.GetByID(ctx, id)
	if err != nil {
		return Incident{}, fmt.Errorf("get incident by id: %w", err)
	}

	status, err := m.statusProvider.GetByID(ctx, result.StatusID)
	if err != nil {
		return Incident{}, fmt.Errorf("get status: %w", err)
	}

	impact, err := m.impactProvider.GetByID(ctx, result.ImpactID)
	if err != nil {
		return Incident{}, fmt.Errorf("get impact: %w", err)
	}

	priority, err := m.priorityProvider.GetByID(ctx, result.PriorityID)
	if err != nil {
		return Incident{}, fmt.Errorf("get priority: %w", err)
	}

	return m.fromDB(result, status, impact, priority), nil
}

func (m *IncidentManager) Update(ctx context.Context, cmd UpdateIncident) (Incident, error) {
	current, err := m.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return Incident{}, fmt.Errorf("get incident: %w", err)
	}

	if cmd.Title != nil {
		current.Title = *cmd.Title
	}

	if cmd.Description != nil {
		current.Description = *cmd.Description
	}

	if cmd.ImpactID != nil {
		current.ImpactID = *cmd.ImpactID
	}

	if cmd.StatusID != nil {
		current.StatusID = *cmd.StatusID
	}

	if cmd.PriorityID != nil {
		current.PriorityID = *cmd.PriorityID
	}

	if err := m.repo.Update(ctx, &current); err != nil {
		return Incident{}, fmt.Errorf("update incident: %w", err)
	}

	return m.GetByID(ctx, cmd.ID)
}

func (m *IncidentManager) Delete(ctx context.Context, id int64) error {
	if err := m.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete incident: %w", err)
	}

	return nil
}

func (m *IncidentManager) List(ctx context.Context, limit, offset int) ([]Incident, error) {
	rows, err := m.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list incidents: %w", err)
	}

	statuses, err := m.statusProvider.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list statuses: %w", err)
	}

	impacts, err := m.impactProvider.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list impacts: %w", err)
	}

	priorities, err := m.priorityProvider.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list priorities: %w", err)
	}

	statusMap := make(map[int64]dict.Entity, len(statuses))
	for _, e := range statuses {
		statusMap[e.ID] = e
	}

	impactMap := make(map[int64]dict.Entity, len(impacts))
	for _, e := range impacts {
		impactMap[e.ID] = e
	}

	priorityMap := make(map[int64]dict.Entity, len(priorities))
	for _, e := range priorities {
		priorityMap[e.ID] = e
	}

	result := make([]Incident, 0, len(rows))

	for _, r := range rows {
		status, ok := statusMap[r.StatusID]
		if !ok {
			return nil, fmt.Errorf("status %d not found", r.StatusID)
		}

		impact, ok := impactMap[r.ImpactID]
		if !ok {
			return nil, fmt.Errorf("impact %d not found", r.ImpactID)
		}

		priority, ok := priorityMap[r.PriorityID]
		if !ok {
			return nil, fmt.Errorf("priority %d not found", r.PriorityID)
		}

		result = append(result, m.fromDB(r, status, impact, priority))
	}

	return result, nil
}

func (m *IncidentManager) fromDB(incident db.Incident, status dict.Entity, impact dict.Entity, priority dict.Entity) Incident {
	return Incident{
		ID:          incident.ID,
		Title:       incident.Title,
		Description: incident.Description,
		Impact:      impact,
		Status:      status,
		Priority:    priority,
	}
}
