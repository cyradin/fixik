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
	dbIncident := db.Incident{
		ID:          cmd.ID,
		Title:       cmd.Title,
		Description: cmd.Description,
		ImpactID:    cmd.ImpactID,
		StatusID:    cmd.StatusID,
		PriorityID:  cmd.PriorityID,
	}

	if err := m.repo.Update(ctx, &dbIncident); err != nil {
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
