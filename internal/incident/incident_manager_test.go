package incident

import (
	"context"
	"errors"
	"testing"

	"github.com/cyradin/fixik/internal/db"
	"github.com/cyradin/fixik/internal/priority"
	"github.com/cyradin/fixik/internal/status"
	"github.com/cyradin/fixik/internal/team"
	"github.com/cyradin/fixik/internal/user"
	"github.com/stretchr/testify/require"
)

func TestIncidentManager_Create(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		cmd  CreateIncident
		mock func(*incidentRepoMock, *statusProviderMock, *priorityProviderMock, *userProviderMock)
		err  bool
	}{
		{
			name: "success",
			cmd: CreateIncident{
				Title:       "title",
				Description: "desc",
				StatusID:    1,
				PriorityID:  2,
			},
			mock: func(
				repo *incidentRepoMock,
				statusProvider *statusProviderMock,
				priorityProvider *priorityProviderMock,
				user *userProviderMock,
			) {
				repo.createFn = func(ctx context.Context, i *db.Incident) error {
					i.ID = 10
					return nil
				}

				repo.getByIDFn = func(ctx context.Context, id int64) (db.Incident, error) {
					return db.Incident{
						ID:          id,
						Title:       "title",
						Description: "desc",
						StatusID:    1,
						PriorityID:  2,
					}, nil
				}

				statusProvider.getByIDFn = func(ctx context.Context, id int64) (status.Status, error) {
					return status.Status{ID: 1}, nil
				}

				priorityProvider.getByIDFn = func(ctx context.Context, id int64) (priority.Priority, error) {
					return priority.Priority{ID: 2}, nil
				}
			},
		},
		{
			name: "repo error",
			cmd: CreateIncident{
				Title: "title",
			},
			mock: func(
				repo *incidentRepoMock,
				statusProvider *statusProviderMock,
				priorityProvider *priorityProviderMock,
				user *userProviderMock,
			) {
				repo.createFn = func(ctx context.Context, i *db.Incident) error {
					return errors.New("db error")
				}
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &incidentRepoMock{}
			status := &statusProviderMock{}
			priority := &priorityProviderMock{}
			user := &userProviderMock{}

			tt.mock(repo, status, priority, user)

			manager := NewIncidentManager(repo, status, priority, nil, user)

			_, err := manager.Create(ctx, tt.cmd)

			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestIncidentManager_GetByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	dbIncident := db.Incident{
		ID:          1,
		Title:       "title",
		Description: "desc",
		StatusID:    1,
		PriorityID:  2,
		UserID:      new(int64(100)),
		AuthorID:    new(int64(100)),
	}

	tests := []struct {
		name string
		mock func(
			*incidentRepoMock,
			*statusProviderMock,
			*priorityProviderMock,
			*userProviderMock,
		)
		erroneous bool
	}{
		{
			name: "repo error",
			mock: func(
				repo *incidentRepoMock,
				status *statusProviderMock,
				priorityProvider *priorityProviderMock,
				user *userProviderMock,
			) {
				repo.getByIDFn = func(ctx context.Context, id int64) (db.Incident, error) {
					return db.Incident{}, errors.New("repo error")
				}
			},
			erroneous: true,
		},
		{
			name: "status error",
			mock: func(
				repo *incidentRepoMock,
				statusProvider *statusProviderMock,
				priorityProvider *priorityProviderMock,
				user *userProviderMock,
			) {
				repo.getByIDFn = func(ctx context.Context, id int64) (db.Incident, error) {
					return dbIncident, nil
				}

				statusProvider.getByIDFn = func(ctx context.Context, id int64) (status.Status, error) {
					return status.Status{}, errors.New("status error")
				}
			},
			erroneous: true,
		},
		{
			name: "priority error",
			mock: func(
				repo *incidentRepoMock,
				statusProvider *statusProviderMock,
				priorityProvider *priorityProviderMock,
				user *userProviderMock,
			) {
				repo.getByIDFn = func(ctx context.Context, id int64) (db.Incident, error) {
					return dbIncident, nil
				}

				statusProvider.getByIDFn = func(ctx context.Context, id int64) (status.Status, error) {
					return status.Status{ID: 1}, nil
				}

				priorityProvider.getByIDFn = func(ctx context.Context, id int64) (priority.Priority, error) {
					return priority.Priority{}, errors.New("priority error")
				}
			},
			erroneous: true,
		},
		{
			name: "success",
			mock: func(
				repo *incidentRepoMock,
				statusProvider *statusProviderMock,
				priorityProvider *priorityProviderMock,
				up *userProviderMock,
			) {
				repo.getByIDFn = func(ctx context.Context, id int64) (db.Incident, error) {
					return dbIncident, nil
				}

				statusProvider.getByIDFn = func(ctx context.Context, id int64) (status.Status, error) {
					return status.Status{ID: 1}, nil
				}

				priorityProvider.getByIDFn = func(ctx context.Context, id int64) (priority.Priority, error) {
					return priority.Priority{ID: 2}, nil
				}

				up.getByIDFn = func(ctx context.Context, id int64) (user.User, error) {
					return user.User{ID: 100}, nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &incidentRepoMock{}
			status := &statusProviderMock{}
			priority := &priorityProviderMock{}
			team := &teamProviderMock{}
			user := &userProviderMock{}

			tt.mock(repo, status, priority, user)

			manager := NewIncidentManager(repo, status, priority, team, user)

			res, err := manager.GetByID(ctx, 1)

			if tt.erroneous {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			require.Equal(t, int64(1), res.ID)
			require.Equal(t, int64(1), res.Status.ID)
			require.Equal(t, int64(2), res.Priority.ID)
			require.Equal(t, int64(100), res.User.ID)
			require.Equal(t, int64(100), res.Author.ID)
		})
	}
}

func TestIncidentManager_Update(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		cmd  UpdateIncident
		mock func(*incidentRepoMock)
		err  bool
	}{
		{
			name: "success",
			cmd: UpdateIncident{
				ID:          1,
				Title:       new("new"),
				Description: new("desc"),
				StatusID:    new(int64(1)),
				PriorityID:  new(int64(2)),
			},
			mock: func(m *incidentRepoMock) {
				m.updateFn = func(ctx context.Context, i *db.Incident) error {
					return nil
				}
				m.getByIDFn = func(ctx context.Context, id int64) (db.Incident, error) {
					return db.Incident{
						ID:          id,
						Title:       "new",
						Description: "desc",
						StatusID:    1,
						PriorityID:  2,
					}, nil
				}
			},
		},
		{
			name: "success, delete team",
			cmd: UpdateIncident{
				ID:     1,
				TeamID: new(int64(0)),
			},
			mock: func(m *incidentRepoMock) {
				m.updateFn = func(ctx context.Context, i *db.Incident) error {
					return nil
				}
				m.getByIDFn = func(ctx context.Context, id int64) (db.Incident, error) {
					return db.Incident{
						ID:          id,
						Title:       "new",
						Description: "desc",
						StatusID:    1,
						PriorityID:  2,
					}, nil
				}
			},
		},
		{
			name: "success, delete user",
			cmd: UpdateIncident{
				ID:     1,
				UserID: new(int64(0)),
			},
			mock: func(m *incidentRepoMock) {
				m.updateFn = func(ctx context.Context, i *db.Incident) error {
					return nil
				}
				m.getByIDFn = func(ctx context.Context, id int64) (db.Incident, error) {
					return db.Incident{
						ID:          id,
						Title:       "new",
						Description: "desc",
						StatusID:    1,
						PriorityID:  2,
					}, nil
				}
			},
		},
		{
			name: "success, delete author",
			cmd: UpdateIncident{
				ID:       1,
				AuthorID: new(int64(0)),
			},
			mock: func(m *incidentRepoMock) {
				m.updateFn = func(ctx context.Context, i *db.Incident) error {
					return nil
				}
				m.getByIDFn = func(ctx context.Context, id int64) (db.Incident, error) {
					return db.Incident{
						ID:          id,
						Title:       "new",
						Description: "desc",
						StatusID:    1,
						PriorityID:  2,
					}, nil
				}
			},
		},
		{
			name: "repo error",
			cmd:  UpdateIncident{ID: 1},
			mock: func(m *incidentRepoMock) {
				m.getByIDFn = func(ctx context.Context, id int64) (db.Incident, error) {
					return db.Incident{ID: id}, nil
				}

				m.updateFn = func(ctx context.Context, i *db.Incident) error {
					return errors.New("db error")
				}
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &incidentRepoMock{}
			tt.mock(repo)

			manager := NewIncidentManager(
				repo,
				&statusProviderMock{getByIDFn: func(context.Context, status.ID) (status.Status, error) { return status.Status{ID: 2}, nil }},
				&priorityProviderMock{getByIDFn: func(context.Context, priority.ID) (priority.Priority, error) { return priority.Priority{ID: 3}, nil }},
				&teamProviderMock{getByIDFn: func(context.Context, team.ID) (team.Team, error) { return team.Team{ID: 3}, nil }},
				&userProviderMock{getByIDFn: func(ctx context.Context, i int64) (user.User, error) { return user.User{ID: 1}, nil }},
			)

			_, err := manager.Update(ctx, tt.cmd)

			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestIncidentManager_Delete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		id   int64
		mock func(*incidentRepoMock)
		err  bool
	}{
		{
			name: "success",
			id:   1,
			mock: func(m *incidentRepoMock) {
				m.deleteFn = func(ctx context.Context, id int64) error {
					return nil
				}
			},
		},
		{
			name: "repo error",
			id:   1,
			mock: func(m *incidentRepoMock) {
				m.deleteFn = func(ctx context.Context, id int64) error {
					return errors.New("db error")
				}
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &incidentRepoMock{}
			tt.mock(repo)

			manager := NewIncidentManager(
				repo,
				&statusProviderMock{getByIDFn: func(context.Context, status.ID) (status.Status, error) { return status.Status{ID: 2}, nil }},
				&priorityProviderMock{getByIDFn: func(context.Context, priority.ID) (priority.Priority, error) { return priority.Priority{ID: 3}, nil }},
				&teamProviderMock{getByIDFn: func(context.Context, team.ID) (team.Team, error) { return team.Team{ID: 3}, nil }},
				&userProviderMock{getByIDFn: func(ctx context.Context, i int64) (user.User, error) { return user.User{ID: 1}, nil }},
			)

			err := manager.Delete(ctx, tt.id)

			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
func TestIncidentManager_List(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	dbIncidents := []db.Incident{
		{
			ID:          1,
			Title:       "one",
			Description: "desc1",
			StatusID:    1,
			PriorityID:  2,
			UserID:      new(int64(100)),
		},
		{
			ID:          2,
			Title:       "two",
			Description: "desc2",
			StatusID:    1,
			PriorityID:  2,
			UserID:      new(int64(101)),
		},
	}

	tests := []struct {
		name string
		mock func(
			*incidentRepoMock,
			*statusProviderMock,
			*priorityProviderMock,
			*userProviderMock,
		)
		err bool
	}{
		{
			name: "repo error",
			mock: func(
				repo *incidentRepoMock,
				status *statusProviderMock,
				priorityProvider *priorityProviderMock,
				user *userProviderMock,
			) {
				repo.listFn = func(context.Context, db.IncidentFilter, int, int) (db.IncidentListResult, error) {
					return db.IncidentListResult{}, errors.New("repo error")
				}
			},
			err: true,
		},
		{
			name: "status list error",
			mock: func(
				repo *incidentRepoMock,
				statusProvider *statusProviderMock,
				priorityProvider *priorityProviderMock,
				user *userProviderMock,
			) {
				repo.listFn = func(context.Context, db.IncidentFilter, int, int) (db.IncidentListResult, error) {
					return db.IncidentListResult{
						Items: dbIncidents,
						Total: 2,
					}, nil
				}

				statusProvider.listFn = func(context.Context) ([]status.Status, error) {
					return nil, errors.New("status error")
				}
			},
			err: true,
		},
		{
			name: "priority list error",
			mock: func(
				repo *incidentRepoMock,
				statusProvider *statusProviderMock,
				priorityProvider *priorityProviderMock,
				user *userProviderMock,
			) {
				repo.listFn = func(context.Context, db.IncidentFilter, int, int) (db.IncidentListResult, error) {
					return db.IncidentListResult{
						Items: dbIncidents,
						Total: 2,
					}, nil
				}

				statusProvider.listFn = func(context.Context) ([]status.Status, error) {
					return []status.Status{{ID: 1}}, nil
				}

				priorityProvider.listFn = func(context.Context) ([]priority.Priority, error) {
					return nil, errors.New("priority error")
				}
			},
			err: true,
		},
		{
			name: "success",
			mock: func(
				repo *incidentRepoMock,
				statusProvider *statusProviderMock,
				priorityProvider *priorityProviderMock,
				up *userProviderMock,
			) {
				repo.listFn = func(context.Context, db.IncidentFilter, int, int) (db.IncidentListResult, error) {
					return db.IncidentListResult{
						Items: dbIncidents,
						Total: 2,
					}, nil
				}

				statusProvider.listFn = func(context.Context) ([]status.Status, error) {
					return []status.Status{{ID: 1}}, nil
				}

				priorityProvider.listFn = func(context.Context) ([]priority.Priority, error) {
					return []priority.Priority{{ID: 2}}, nil
				}

				up.getByIDManyFn = func(ctx context.Context, ids []int64) ([]user.User, error) {
					users := make([]user.User, 0, len(ids))
					for _, id := range ids {
						users = append(users, user.User{ID: id})
					}

					return users, nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &incidentRepoMock{}
			status := &statusProviderMock{}
			priority := &priorityProviderMock{}
			team := &teamProviderMock{}
			user := &userProviderMock{}

			tt.mock(repo, status, priority, user)

			manager := NewIncidentManager(repo, status, priority, team, user)

			res, err := manager.List(ctx, Filter{}, 10, 0)

			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			require.Len(t, res.Items, 2)
			require.Equal(t, 2, res.Total)

			require.Equal(t, int64(1), res.Items[0].Status.ID)
			require.Equal(t, int64(2), res.Items[0].Priority.ID)
			require.Equal(t, int64(100), res.Items[0].User.ID)
			require.Equal(t, int64(101), res.Items[1].User.ID)
		})
	}
}

type incidentRepoMock struct {
	createFn  func(context.Context, *db.Incident) error
	getByIDFn func(context.Context, int64) (db.Incident, error)
	updateFn  func(context.Context, *db.Incident) error
	deleteFn  func(context.Context, int64) error
	listFn    func(context.Context, db.IncidentFilter, int, int) (db.IncidentListResult, error)
}

func (m *incidentRepoMock) Create(ctx context.Context, i *db.Incident) error {
	return m.createFn(ctx, i)
}

func (m *incidentRepoMock) GetByID(ctx context.Context, id int64) (db.Incident, error) {
	return m.getByIDFn(ctx, id)
}

func (m *incidentRepoMock) Update(ctx context.Context, i *db.Incident) error {
	return m.updateFn(ctx, i)
}

func (m *incidentRepoMock) Delete(ctx context.Context, id int64) error {
	return m.deleteFn(ctx, id)
}

func (m *incidentRepoMock) List(ctx context.Context, filter db.IncidentFilter, limit, offset int) (db.IncidentListResult, error) {
	return m.listFn(ctx, filter, limit, offset)
}

type statusProviderMock struct {
	getByIDFn func(context.Context, status.ID) (status.Status, error)
	listFn    func(context.Context) ([]status.Status, error)
}

func (m *statusProviderMock) GetByID(ctx context.Context, id status.ID) (status.Status, error) {
	return m.getByIDFn(ctx, id)
}

func (m *statusProviderMock) List(ctx context.Context) ([]status.Status, error) {
	if m.listFn != nil {
		return m.listFn(ctx)
	}

	return nil, nil
}

type priorityProviderMock struct {
	getByIDFn func(context.Context, priority.ID) (priority.Priority, error)
	listFn    func(context.Context) ([]priority.Priority, error)
}

func (m *priorityProviderMock) GetByID(ctx context.Context, id priority.ID) (priority.Priority, error) {
	return m.getByIDFn(ctx, id)
}

func (m *priorityProviderMock) List(ctx context.Context) ([]priority.Priority, error) {
	if m.listFn != nil {
		return m.listFn(ctx)
	}

	return nil, nil
}

type teamProviderMock struct {
	getByIDFn func(context.Context, team.ID) (team.Team, error)
	listFn    func(context.Context) ([]team.Team, error)
}

func (m *teamProviderMock) GetByID(ctx context.Context, id team.ID) (team.Team, error) {
	return m.getByIDFn(ctx, id)
}

func (m *teamProviderMock) List(ctx context.Context) ([]team.Team, error) {
	if m.listFn != nil {
		return m.listFn(ctx)
	}

	return nil, nil
}

type userProviderMock struct {
	getByIDFn     func(context.Context, int64) (user.User, error)
	getByIDManyFn func(context.Context, []int64) ([]user.User, error)
}

func (m *userProviderMock) GetByID(ctx context.Context, id int64) (user.User, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}

	return user.User{}, nil
}

func (m *userProviderMock) GetByIDMany(ctx context.Context, ids []int64) ([]user.User, error) {
	if m.getByIDManyFn != nil {
		return m.getByIDManyFn(ctx, ids)
	}

	return nil, nil
}
