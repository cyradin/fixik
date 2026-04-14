package priority

import (
	"context"
	"fmt"

	"github.com/cyradin/fixik/internal/db"
)

var ErrHasDependantEntities = fmt.Errorf("has dependant entities")

type priorityRepo interface {
	Create(ctx context.Context, s *db.DictEntity) error
	GetByID(ctx context.Context, id int64) (db.DictEntity, error)
	List(ctx context.Context) ([]db.DictEntity, error)
	Update(ctx context.Context, s *db.DictEntity) error
	Delete(ctx context.Context, id int64) error
}

type incidentsByPriorityCounter interface {
	CountByPriority(ctx context.Context, priorityID int64) (int, error)
}

type txExecutor interface {
	Exec(ctx context.Context, callback func(ctx context.Context) error) error
}

type PriorityManager struct {
	repo             priorityRepo
	incidentsCounter incidentsByPriorityCounter
	txExecutor       txExecutor
}

func NewPriorityManager(
	repo priorityRepo,
	incidentsCounter incidentsByPriorityCounter,
	txExecutor txExecutor,
) *PriorityManager {
	return &PriorityManager{
		repo:             repo,
		incidentsCounter: incidentsCounter,
		txExecutor:       txExecutor,
	}
}

func (m *PriorityManager) Create(ctx context.Context, e Priority) (Priority, error) {
	dbEntity := m.toDB(e)

	if err := m.repo.Create(ctx, &dbEntity); err != nil {
		return Priority{}, fmt.Errorf("repository create: %w", err)
	}

	result := m.fromDB(dbEntity)

	return result, nil
}

func (m *PriorityManager) GetByID(ctx context.Context, id ID) (Priority, error) {
	result, err := m.repo.GetByID(ctx, id)
	if err != nil {
		return Priority{}, fmt.Errorf("repository get by id: %w", err)
	}

	return m.fromDB(result), nil
}

func (m *PriorityManager) List(ctx context.Context) ([]Priority, error) {
	items, err := m.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("repository get list: %w", err)
	}

	result := make([]Priority, 0, len(items))
	for _, item := range items {
		result = append(result, m.fromDB(item))
	}

	return result, nil
}

func (m *PriorityManager) Update(ctx context.Context, e Priority) (Priority, error) {
	dbEntity := m.toDB(e)

	if err := m.repo.Update(ctx, &dbEntity); err != nil {
		return Priority{}, fmt.Errorf("repository update: %w", err)
	}

	result := m.fromDB(dbEntity)

	return result, nil
}

func (m *PriorityManager) Delete(ctx context.Context, id ID) error {
	return m.txExecutor.Exec(ctx, func(ctx context.Context) error {
		cnt, err := m.incidentsCounter.CountByPriority(ctx, id)
		if err != nil {
			return fmt.Errorf("count incidents: %w", err)
		}

		if cnt > 0 {
			return ErrHasDependantEntities
		}

		if err := m.repo.Delete(ctx, id); err != nil {
			return fmt.Errorf("repository delete: %w", err)
		}

		return nil
	})
}

func (m *PriorityManager) toDB(entity Priority) db.DictEntity {
	return db.DictEntity{
		ID:          entity.ID,
		Name:        entity.Name,
		Code:        entity.Code,
		Sort:        entity.Sort,
		Description: entity.Description,
	}
}

func (m *PriorityManager) fromDB(entity db.DictEntity) Priority {
	return Priority{
		ID:          entity.ID,
		Name:        entity.Name,
		Code:        entity.Code,
		Sort:        entity.Sort,
		Description: entity.Description,
	}
}
