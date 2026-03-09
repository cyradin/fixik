package incident

import (
	"context"
	"fmt"

	"github.com/cyradin/fixik/internal/db"
)

type priorityRepo interface {
	Create(ctx context.Context, s *db.Priority) error
	GetByID(ctx context.Context, id int64) (db.Priority, error)
	List(ctx context.Context) ([]db.Priority, error)
	Update(ctx context.Context, s *db.Priority) error
	Delete(ctx context.Context, id int64) error
}

type PriorityManager struct {
	repo priorityRepo
}

func NewPriorityManager(repo priorityRepo) *PriorityManager {
	return &PriorityManager{
		repo: repo,
	}
}

func (m *PriorityManager) Create(ctx context.Context, priority Priority) (Priority, error) {
	dbPriority := m.toDB(priority)

	if err := m.repo.Create(ctx, &dbPriority); err != nil {
		return Priority{}, fmt.Errorf("create priority: %w", err)
	}

	result := m.fromDB(dbPriority)

	return result, nil
}

func (m *PriorityManager) GetByID(ctx context.Context, id PriorityID) (Priority, error) {
	result, err := m.repo.GetByID(ctx, id)
	if err != nil {
		return Priority{}, fmt.Errorf("get priority by id: %w", err)
	}

	return m.fromDB(result), nil
}

func (m *PriorityManager) List(ctx context.Context) ([]Priority, error) {
	items, err := m.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("get priority list: %w", err)
	}

	result := make([]Priority, 0, len(items))
	for _, item := range items {
		result = append(result, m.fromDB(item))
	}

	return result, nil
}

func (m *PriorityManager) Update(ctx context.Context, priority Priority) (Priority, error) {
	dbPriority := m.toDB(priority)

	if err := m.repo.Update(ctx, &dbPriority); err != nil {
		return Priority{}, fmt.Errorf("update priority: %w", err)
	}

	result := m.fromDB(dbPriority)

	return result, nil
}

func (m *PriorityManager) Delete(ctx context.Context, id PriorityID) error {
	if err := m.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete priority: %w", err)
	}

	return nil
}

func (m *PriorityManager) toDB(priority Priority) db.Priority {
	return db.Priority{
		ID:   priority.ID,
		Name: priority.Name,
		Code: priority.Code,
	}
}

func (m *PriorityManager) fromDB(priority db.Priority) Priority {
	return Priority{
		ID:   priority.ID,
		Name: priority.Name,
		Code: priority.Code,
	}
}
