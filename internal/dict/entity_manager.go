package dict

import (
	"context"
	"fmt"

	"github.com/cyradin/fixik/internal/db"
)

type entityRepo interface {
	Create(ctx context.Context, s *db.DictEntity) error
	GetByID(ctx context.Context, id int64) (db.DictEntity, error)
	List(ctx context.Context) ([]db.DictEntity, error)
	Update(ctx context.Context, s *db.DictEntity) error
	Delete(ctx context.Context, id int64) error
}

type EntityManager struct {
	repo entityRepo
}

func newEntityManager(repo entityRepo) *EntityManager {
	return &EntityManager{
		repo: repo,
	}
}

func NewImpactManager(repo entityRepo) *EntityManager {
	return newEntityManager(repo)
}

func NewPriorityManager(repo entityRepo) *EntityManager {
	return newEntityManager(repo)
}

func NewStatusManager(repo entityRepo) *EntityManager {
	return newEntityManager(repo)
}

func NewTeamManager(repo entityRepo) *EntityManager {
	return newEntityManager(repo)
}

func NewRoleManager(repo entityRepo) *EntityManager {
	return newEntityManager(repo)
}

func (m *EntityManager) Create(ctx context.Context, status Entity) (Entity, error) {
	dbEntity := m.toDB(status)

	if err := m.repo.Create(ctx, &dbEntity); err != nil {
		return Entity{}, fmt.Errorf("create status: %w", err)
	}

	result := m.fromDB(dbEntity)

	return result, nil
}

func (m *EntityManager) GetByID(ctx context.Context, id EntityID) (Entity, error) {
	result, err := m.repo.GetByID(ctx, id)
	if err != nil {
		return Entity{}, fmt.Errorf("get status by id: %w", err)
	}

	return m.fromDB(result), nil
}

func (m *EntityManager) List(ctx context.Context) ([]Entity, error) {
	items, err := m.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("get status list: %w", err)
	}

	result := make([]Entity, 0, len(items))
	for _, item := range items {
		result = append(result, m.fromDB(item))
	}

	return result, nil
}

func (m *EntityManager) Update(ctx context.Context, status Entity) (Entity, error) {
	dbEntity := m.toDB(status)

	if err := m.repo.Update(ctx, &dbEntity); err != nil {
		return Entity{}, fmt.Errorf("update status: %w", err)
	}

	result := m.fromDB(dbEntity)

	return result, nil
}

func (m *EntityManager) Delete(ctx context.Context, id EntityID) error {
	if err := m.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete status: %w", err)
	}

	return nil
}

func (m *EntityManager) toDB(status Entity) db.DictEntity {
	return db.DictEntity{
		ID:   status.ID,
		Name: status.Name,
		Code: status.Code,
	}
}

func (m *EntityManager) fromDB(status db.DictEntity) Entity {
	return Entity{
		ID:   status.ID,
		Name: status.Name,
		Code: status.Code,
	}
}
