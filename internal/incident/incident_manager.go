package incident

import (
	"context"
	"fmt"

	"github.com/cyradin/fixik/internal/db"
	"github.com/cyradin/fixik/internal/dict"
	"github.com/cyradin/fixik/internal/user"
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

type userProvider interface {
	GetByID(ctx context.Context, id int64) (user.User, error)
	GetByIDMany(ctx context.Context, ids []int64) ([]user.User, error)
}

type IncidentManager struct {
	repo             incidentRepo
	statusProvider   entityProvider
	priorityProvider entityProvider
	teamProvider     entityProvider
	userProvider     userProvider
}

func NewIncidentManager(
	repo incidentRepo,
	statusProvider entityProvider,
	priorityProvider entityProvider,
	teamProvider entityProvider,
	userProvider userProvider,
) *IncidentManager {
	return &IncidentManager{
		repo:             repo,
		statusProvider:   statusProvider,
		priorityProvider: priorityProvider,
		teamProvider:     teamProvider,
		userProvider:     userProvider,
	}
}

func (m *IncidentManager) Create(ctx context.Context, incident CreateIncident) (Incident, error) {
	dbIncident := db.Incident{
		Title:       incident.Title,
		Description: incident.Description,
		StatusID:    incident.StatusID,
		PriorityID:  incident.PriorityID,
		TeamID:      incident.TeamID,
		UserID:      incident.UserID,
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

	priority, err := m.priorityProvider.GetByID(ctx, result.PriorityID)
	if err != nil {
		return Incident{}, fmt.Errorf("get priority: %w", err)
	}

	var team *dict.Entity

	if result.TeamID != nil {
		t, err := m.teamProvider.GetByID(ctx, *result.TeamID)
		if err != nil {
			return Incident{}, fmt.Errorf("get team: %w", err)
		}

		team = &t
	}

	var usr *user.User

	if result.UserID != nil {
		u, err := m.userProvider.GetByID(ctx, *result.UserID)
		if err != nil {
			return Incident{}, fmt.Errorf("get user: %w", err)
		}

		usr = &u
	}

	var author *user.User

	if result.AuthorID != nil {
		u, err := m.userProvider.GetByID(ctx, *result.AuthorID)
		if err != nil {
			return Incident{}, fmt.Errorf("get user: %w", err)
		}

		author = &u
	}

	return m.fromDB(result, status, priority, team, usr, author), nil
}

func (m *IncidentManager) Update(ctx context.Context, incident UpdateIncident) (Incident, error) {
	current, err := m.repo.GetByID(ctx, incident.ID)
	if err != nil {
		return Incident{}, fmt.Errorf("get incident: %w", err)
	}

	if incident.Title != nil {
		current.Title = *incident.Title
	}

	if incident.Description != nil {
		current.Description = *incident.Description
	}

	if incident.StatusID != nil {
		current.StatusID = *incident.StatusID
	}

	if incident.PriorityID != nil {
		current.PriorityID = *incident.PriorityID
	}

	if incident.TeamID != nil {
		if *incident.TeamID == 0 {
			// удалить команду
			incident.TeamID = nil
		} else {
			current.TeamID = incident.TeamID
		}
	}

	if incident.UserID != nil {
		if *incident.UserID == 0 {
			incident.UserID = nil
		} else {
			current.UserID = incident.UserID
		}
	}

	if incident.AuthorID != nil {
		if *incident.AuthorID == 0 {
			incident.AuthorID = nil
		} else {
			current.AuthorID = incident.AuthorID
		}
	}

	if err := m.repo.Update(ctx, &current); err != nil {
		return Incident{}, fmt.Errorf("update incident: %w", err)
	}

	return m.GetByID(ctx, incident.ID)
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

	userIDs := make([]int64, 0, len(rows))
	userIDSet := make(map[int64]struct{})

	for _, r := range rows {
		if r.UserID != nil {
			id := *r.UserID
			if _, exists := userIDSet[id]; !exists {
				userIDs = append(userIDs, id)
				userIDSet[id] = struct{}{}
			}
		}
	}

	usersMap := make(map[int64]user.User)

	if len(userIDs) > 0 {
		users, err := m.userProvider.GetByIDMany(ctx, userIDs)
		if err != nil {
			return nil, fmt.Errorf("get users: %w", err)
		}

		for _, u := range users {
			usersMap[u.ID] = u
		}
	}

	statusMap, err := m.loadStatuses(ctx)
	if err != nil {
		return nil, err
	}

	priorityMap, err := m.loadPriorities(ctx)
	if err != nil {
		return nil, err
	}

	teamMap, err := m.loadTeams(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]Incident, 0, len(rows))

	for _, r := range rows {
		item, err := m.transformFromDB(
			r,
			statusMap,
			priorityMap,
			teamMap,
			usersMap,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, nil
}

func (m *IncidentManager) loadStatuses(ctx context.Context) (map[int64]dict.Entity, error) {
	items, err := m.statusProvider.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list statuses: %w", err)
	}

	result := make(map[int64]dict.Entity, len(items))
	for _, e := range items {
		result[e.ID] = e
	}

	return result, nil
}

func (m *IncidentManager) loadPriorities(ctx context.Context) (map[int64]dict.Entity, error) {
	items, err := m.priorityProvider.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list priorities: %w", err)
	}

	result := make(map[int64]dict.Entity, len(items))
	for _, e := range items {
		result[e.ID] = e
	}

	return result, nil
}

func (m *IncidentManager) loadTeams(ctx context.Context) (map[int64]dict.Entity, error) {
	items, err := m.teamProvider.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list teams: %w", err)
	}

	result := make(map[int64]dict.Entity, len(items))
	for _, e := range items {
		result[e.ID] = e
	}

	return result, nil
}

func (m *IncidentManager) transformFromDB(
	r db.Incident,
	statusMap map[int64]dict.Entity,
	priorityMap map[int64]dict.Entity,
	teamMap map[int64]dict.Entity,
	userMap map[int64]user.User,
) (Incident, error) {
	status, ok := statusMap[r.StatusID]
	if !ok {
		return Incident{}, fmt.Errorf("status %d not found", r.StatusID)
	}

	priority, ok := priorityMap[r.PriorityID]
	if !ok {
		return Incident{}, fmt.Errorf("priority %d not found", r.PriorityID)
	}

	var team *dict.Entity

	if r.TeamID != nil {
		t, ok := teamMap[*r.TeamID]
		if !ok {
			return Incident{}, fmt.Errorf("team %d not found", *r.TeamID)
		}

		team = &t
	}

	var usr *user.User

	if r.UserID != nil {
		u, ok := userMap[*r.UserID]

		if !ok {
			return Incident{}, fmt.Errorf("user %d not found", *r.UserID)
		}

		usr = &u
	}

	var author *user.User

	if r.AuthorID != nil {
		u, ok := userMap[*r.AuthorID]

		if !ok {
			return Incident{}, fmt.Errorf("user %d not found", *r.AuthorID)
		}

		author = &u
	}

	return m.fromDB(r, status, priority, team, usr, author), nil
}

func (m *IncidentManager) fromDB(incident db.Incident, status dict.Entity, priority dict.Entity, team *dict.Entity, user *user.User, author *user.User) Incident {
	return Incident{
		ID:          incident.ID,
		Title:       incident.Title,
		Description: incident.Description,
		Status:      status,
		Priority:    priority,
		Team:        team,
		User:        user,
		Author:      author,
	}
}
