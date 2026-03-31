package incident

import (
	"context"
	"errors"
	"testing"

	"github.com/cyradin/fixik/internal/db"
	"github.com/cyradin/fixik/internal/dict"
	"github.com/cyradin/fixik/internal/user"
	"github.com/stretchr/testify/require"
)

func TestIncidentManager_Create(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		cmd  CreateIncident
		mock func(*incidentRepoMock, *entityProviderMock, *entityProviderMock, *userProviderMock)
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
			mock: func(repo *incidentRepoMock, status *entityProviderMock, priority *entityProviderMock, user *userProviderMock) {
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
				status.getByIDFn = func(ctx context.Context, id dict.EntityID) (dict.Entity, error) {
					return dict.Entity{ID: 1}, nil
				}
				priority.getByIDFn = func(ctx context.Context, id dict.EntityID) (dict.Entity, error) {
					return dict.Entity{ID: 2}, nil
				}
			},
		},
		{
			name: "repo error",
			cmd: CreateIncident{
				Title: "title",
			},
			mock: func(repo *incidentRepoMock, status *entityProviderMock, priority *entityProviderMock, user *userProviderMock) {
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
			status := &entityProviderMock{}
			priority := &entityProviderMock{}
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
			*entityProviderMock,
			*entityProviderMock,
			*userProviderMock,
		)
		erroneous bool
	}{
		{
			name: "repo error",
			mock: func(
				repo *incidentRepoMock,
				status *entityProviderMock,
				priority *entityProviderMock,
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
				status *entityProviderMock,
				priority *entityProviderMock,
				user *userProviderMock,
			) {
				repo.getByIDFn = func(ctx context.Context, id int64) (db.Incident, error) {
					return dbIncident, nil
				}

				status.getByIDFn = func(ctx context.Context, id dict.EntityID) (dict.Entity, error) {
					return dict.Entity{}, errors.New("status error")
				}
			},
			erroneous: true,
		},
		{
			name: "priority error",
			mock: func(
				repo *incidentRepoMock,
				status *entityProviderMock,
				priority *entityProviderMock,
				user *userProviderMock,
			) {
				repo.getByIDFn = func(ctx context.Context, id int64) (db.Incident, error) {
					return dbIncident, nil
				}

				status.getByIDFn = func(ctx context.Context, id dict.EntityID) (dict.Entity, error) {
					return dict.Entity{ID: 1}, nil
				}

				priority.getByIDFn = func(ctx context.Context, id dict.EntityID) (dict.Entity, error) {
					return dict.Entity{}, errors.New("priority error")
				}
			},
			erroneous: true,
		},
		{
			name: "success",
			mock: func(
				repo *incidentRepoMock,
				status *entityProviderMock,
				priority *entityProviderMock,
				up *userProviderMock,
			) {
				repo.getByIDFn = func(ctx context.Context, id int64) (db.Incident, error) {
					return dbIncident, nil
				}

				status.getByIDFn = func(ctx context.Context, id dict.EntityID) (dict.Entity, error) {
					return dict.Entity{ID: 1}, nil
				}

				priority.getByIDFn = func(ctx context.Context, id dict.EntityID) (dict.Entity, error) {
					return dict.Entity{ID: 2}, nil
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
			status := &entityProviderMock{}
			priority := &entityProviderMock{}
			team := &entityProviderMock{}
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
			require.Equal(t, dict.EntityID(1), res.Status.ID)
			require.Equal(t, dict.EntityID(2), res.Priority.ID)
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
				&entityProviderMock{getByIDFn: func(context.Context, dict.EntityID) (dict.Entity, error) { return dict.Entity{ID: 2}, nil }},
				&entityProviderMock{getByIDFn: func(context.Context, dict.EntityID) (dict.Entity, error) { return dict.Entity{ID: 3}, nil }},
				&entityProviderMock{getByIDFn: func(context.Context, dict.EntityID) (dict.Entity, error) { return dict.Entity{ID: 3}, nil }},
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
				&entityProviderMock{getByIDFn: func(context.Context, dict.EntityID) (dict.Entity, error) { return dict.Entity{ID: 2}, nil }},
				&entityProviderMock{getByIDFn: func(context.Context, dict.EntityID) (dict.Entity, error) { return dict.Entity{ID: 3}, nil }},
				&entityProviderMock{getByIDFn: func(context.Context, dict.EntityID) (dict.Entity, error) { return dict.Entity{ID: 3}, nil }},
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
			*entityProviderMock,
			*entityProviderMock,
			*userProviderMock,
		)
		err bool
	}{
		{
			name: "repo error",
			mock: func(
				repo *incidentRepoMock,
				status *entityProviderMock,
				priority *entityProviderMock,
				user *userProviderMock,
			) {
				repo.listFn = func(context.Context, int, int) (db.IncidentListResult, error) {
					return db.IncidentListResult{}, errors.New("repo error")
				}
			},
			err: true,
		},
		{
			name: "status list error",
			mock: func(
				repo *incidentRepoMock,
				status *entityProviderMock,
				priority *entityProviderMock,
				user *userProviderMock,
			) {
				repo.listFn = func(context.Context, int, int) (db.IncidentListResult, error) {
					return db.IncidentListResult{
						Items: dbIncidents,
						Total: 2,
					}, nil
				}

				status.listFn = func(context.Context) ([]dict.Entity, error) {
					return nil, errors.New("status error")
				}
			},
			err: true,
		},
		{
			name: "priority list error",
			mock: func(
				repo *incidentRepoMock,
				status *entityProviderMock,
				priority *entityProviderMock,
				user *userProviderMock,
			) {
				repo.listFn = func(context.Context, int, int) (db.IncidentListResult, error) {
					return db.IncidentListResult{
						Items: dbIncidents,
						Total: 2,
					}, nil
				}

				status.listFn = func(context.Context) ([]dict.Entity, error) {
					return []dict.Entity{{ID: 1}}, nil
				}

				priority.listFn = func(context.Context) ([]dict.Entity, error) {
					return nil, errors.New("priority error")
				}
			},
			err: true,
		},
		{
			name: "success",
			mock: func(
				repo *incidentRepoMock,
				status *entityProviderMock,
				priority *entityProviderMock,
				up *userProviderMock,
			) {
				repo.listFn = func(context.Context, int, int) (db.IncidentListResult, error) {
					return db.IncidentListResult{
						Items: dbIncidents,
						Total: 2,
					}, nil
				}

				status.listFn = func(context.Context) ([]dict.Entity, error) {
					return []dict.Entity{{ID: 1}}, nil
				}

				priority.listFn = func(context.Context) ([]dict.Entity, error) {
					return []dict.Entity{{ID: 2}}, nil
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
			status := &entityProviderMock{}
			priority := &entityProviderMock{}
			team := &entityProviderMock{}
			user := &userProviderMock{}

			tt.mock(repo, status, priority, user)

			manager := NewIncidentManager(repo, status, priority, team, user)

			res, err := manager.List(ctx, 10, 0)

			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			require.Len(t, res.Items, 2)

			require.Equal(t, 2, res.Total)

			require.Equal(t, dict.EntityID(1), res.Items[0].Status.ID)
			require.Equal(t, dict.EntityID(2), res.Items[0].Priority.ID)
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
	listFn    func(context.Context, int, int) (db.IncidentListResult, error)
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

func (m *incidentRepoMock) List(ctx context.Context, limit, offset int) (db.IncidentListResult, error) {
	return m.listFn(ctx, limit, offset)
}

type entityProviderMock struct {
	getByIDFn func(context.Context, dict.EntityID) (dict.Entity, error)
	listFn    func(context.Context) ([]dict.Entity, error)
}

func (m *entityProviderMock) GetByID(ctx context.Context, id dict.EntityID) (dict.Entity, error) {
	return m.getByIDFn(ctx, id)
}

func (m *entityProviderMock) List(ctx context.Context) ([]dict.Entity, error) {
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
