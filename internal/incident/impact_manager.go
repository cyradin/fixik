package incident

import (
	"context"
	"fmt"

	"github.com/cyradin/fixik/internal/db"
)

type impactRepo interface {
	Create(ctx context.Context, s *db.Impact) error
	GetByID(ctx context.Context, id int64) (db.Impact, error)
	List(ctx context.Context) ([]db.Impact, error)
	Update(ctx context.Context, s *db.Impact) error
	Delete(ctx context.Context, id int64) error
}

type ImpactManager struct {
	repo impactRepo
}

func NewImpactManager(repo impactRepo) *ImpactManager {
	return &ImpactManager{
		repo: repo,
	}
}

func (m *ImpactManager) Create(ctx context.Context, impact Impact) (Impact, error) {
	dbImpact := m.toDB(impact)

	if err := m.repo.Create(ctx, &dbImpact); err != nil {
		return Impact{}, fmt.Errorf("create impact: %w", err)
	}

	result := m.fromDB(dbImpact)

	return result, nil
}

func (m *ImpactManager) GetByID(ctx context.Context, id ImpactID) (Impact, error) {
	result, err := m.repo.GetByID(ctx, id)
	if err != nil {
		return Impact{}, fmt.Errorf("get impact by id: %w", err)
	}

	return m.fromDB(result), nil
}

func (m *ImpactManager) List(ctx context.Context) ([]Impact, error) {
	items, err := m.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("get impact list: %w", err)
	}

	result := make([]Impact, 0, len(items))
	for _, item := range items {
		result = append(result, m.fromDB(item))
	}

	return result, nil
}

func (m *ImpactManager) Update(ctx context.Context, impact Impact) (Impact, error) {
	dbImpact := m.toDB(impact)

	if err := m.repo.Update(ctx, &dbImpact); err != nil {
		return Impact{}, fmt.Errorf("update impact: %w", err)
	}

	result := m.fromDB(dbImpact)

	return result, nil
}

func (m *ImpactManager) Delete(ctx context.Context, id ImpactID) error {
	if err := m.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete impact: %w", err)
	}

	return nil
}

func (m *ImpactManager) toDB(impact Impact) db.Impact {
	return db.Impact{
		ID:   impact.ID,
		Name: impact.Name,
		Code: impact.Code,
	}
}

func (m *ImpactManager) fromDB(impact db.Impact) Impact {
	return Impact{
		ID:   impact.ID,
		Name: impact.Name,
		Code: impact.Code,
	}
}
