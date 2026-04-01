package status

import (
	"context"
	"fmt"

	"github.com/cyradin/fixik/internal/db"
)

type statusRepo interface {
	Create(ctx context.Context, s *db.Status) error
	GetByID(ctx context.Context, id int64) (db.Status, error)
	List(ctx context.Context) ([]db.Status, error)
	Update(ctx context.Context, s *db.Status) error
	Delete(ctx context.Context, id int64) error
}

type StatusManager struct {
	repo statusRepo
}

func NewStatusManager(repo statusRepo) *StatusManager {
	return &StatusManager{repo: repo}
}

func (m *StatusManager) Create(ctx context.Context, e Status) (Status, error) {
	dbStatus := m.toDB(e)

	if err := m.repo.Create(ctx, &dbStatus); err != nil {
		return Status{}, fmt.Errorf("repository create: %w", err)
	}

	result := m.fromDB(dbStatus)

	return result, nil
}

func (m *StatusManager) GetByID(ctx context.Context, id StatusID) (Status, error) {
	result, err := m.repo.GetByID(ctx, id)
	if err != nil {
		return Status{}, fmt.Errorf("repository get by id: %w", err)
	}

	return m.fromDB(result), nil
}

func (m *StatusManager) List(ctx context.Context) ([]Status, error) {
	items, err := m.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("repository get list: %w", err)
	}

	result := make([]Status, 0, len(items))
	for _, item := range items {
		result = append(result, m.fromDB(item))
	}

	return result, nil
}

func (m *StatusManager) Update(ctx context.Context, e Status) (Status, error) {
	dbStatus := m.toDB(e)

	if err := m.repo.Update(ctx, &dbStatus); err != nil {
		return Status{}, fmt.Errorf("repository update: %w", err)
	}

	result := m.fromDB(dbStatus)

	return result, nil
}

func (m *StatusManager) Delete(ctx context.Context, id StatusID) error {
	if err := m.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("repository delete: %w", err)
	}

	return nil
}

func (m *StatusManager) toDB(status Status) db.Status {
	return db.Status{
		ID:          status.ID,
		Name:        status.Name,
		Code:        status.Code,
		Sort:        status.Sort,
		IsFinal:     status.IsFinal,
		Description: status.Description,
	}
}

func (m *StatusManager) fromDB(status db.Status) Status {
	return Status{
		ID:          status.ID,
		Name:        status.Name,
		Code:        status.Code,
		Sort:        status.Sort,
		IsFinal:     status.IsFinal,
		Description: status.Description,
	}
}
