package team

import (
	"context"
	"fmt"

	"github.com/cyradin/fixik/internal/db"
)

var (
	ErrHasDependantIncidents = fmt.Errorf("has dependant incidents")
	ErrHasDependantUsers     = fmt.Errorf("has dependant users")
)

type teamRepo interface {
	Create(ctx context.Context, s *db.Team) error
	GetByID(ctx context.Context, id int64) (db.Team, error)
	List(ctx context.Context) ([]db.Team, error)
	Update(ctx context.Context, s *db.Team) error
	Delete(ctx context.Context, id int64) error
}

type incidentsByTeamCounter interface {
	CountByTeam(ctx context.Context, statusID int64) (int, error)
}

type usersByTeamCounter interface {
	CountByTeam(ctx context.Context, teamID int64) (int, error)
}

type txExecutor interface {
	Exec(ctx context.Context, callback func(ctx context.Context) error) error
}

type TeamManager struct {
	repo             teamRepo
	incidentsCounter incidentsByTeamCounter
	usersCounter     usersByTeamCounter
	txExecutor       txExecutor
}

func NewTeamManager(
	repo teamRepo,
	incidentsCounter incidentsByTeamCounter,
	usersCounter usersByTeamCounter,
	txExecutor txExecutor,
) *TeamManager {
	return &TeamManager{
		repo:             repo,
		incidentsCounter: incidentsCounter,
		usersCounter:     usersCounter,
		txExecutor:       txExecutor,
	}
}

func (m *TeamManager) Create(ctx context.Context, e Team) (Team, error) {
	dbTeam := m.toDB(e)

	if err := m.repo.Create(ctx, &dbTeam); err != nil {
		return Team{}, fmt.Errorf("repository create: %w", err)
	}

	result := m.fromDB(dbTeam)

	return result, nil
}

func (m *TeamManager) GetByID(ctx context.Context, id ID) (Team, error) {
	result, err := m.repo.GetByID(ctx, id)
	if err != nil {
		return Team{}, fmt.Errorf("repository get by id: %w", err)
	}

	return m.fromDB(result), nil
}

func (m *TeamManager) List(ctx context.Context) ([]Team, error) {
	items, err := m.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("repository get list: %w", err)
	}

	result := make([]Team, 0, len(items))
	for _, item := range items {
		result = append(result, m.fromDB(item))
	}

	return result, nil
}

func (m *TeamManager) Update(ctx context.Context, e Team) (Team, error) {
	dbTeam := m.toDB(e)

	if err := m.repo.Update(ctx, &dbTeam); err != nil {
		return Team{}, fmt.Errorf("repository update: %w", err)
	}

	result := m.fromDB(dbTeam)

	return result, nil
}

func (m *TeamManager) Delete(ctx context.Context, id ID) error {
	return m.txExecutor.Exec(ctx, func(ctx context.Context) error {
		cnt, err := m.incidentsCounter.CountByTeam(ctx, id)
		if err != nil {
			return fmt.Errorf("count incidents: %w", err)
		}

		if cnt > 0 {
			return ErrHasDependantIncidents
		}

		cnt, err = m.usersCounter.CountByTeam(ctx, id)
		if err != nil {
			return fmt.Errorf("count users: %w", err)
		}

		if cnt > 0 {
			return ErrHasDependantUsers
		}

		if err := m.repo.Delete(ctx, id); err != nil {
			return fmt.Errorf("repository delete: %w", err)
		}

		return nil
	})
}

func (m *TeamManager) toDB(team Team) db.Team {
	return db.Team{
		ID:          team.ID,
		Name:        team.Name,
		Code:        team.Code,
		Sort:        team.Sort,
		Description: team.Description,
	}
}

func (m *TeamManager) fromDB(team db.Team) Team {
	return Team{
		ID:          team.ID,
		Name:        team.Name,
		Code:        team.Code,
		Sort:        team.Sort,
		Description: team.Description,
	}
}
