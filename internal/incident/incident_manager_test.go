package incident

import (
	"context"
	"errors"
	"testing"

	"github.com/cyradin/fixik/internal/db"
	"github.com/cyradin/fixik/internal/dict"
	"github.com/stretchr/testify/require"
)

func TestIncidentManager_Create(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		cmd  CreateIncident
		mock func(*incidentRepoMock)
		err  bool
	}{
		{
			name: "success",
			cmd: CreateIncident{
				Title:       "title",
				Description: "desc",
				StatusID:    1,
				ImpactID:    2,
				PriorityID:  3,
			},
			mock: func(m *incidentRepoMock) {
				m.createFn = func(ctx context.Context, i *db.Incident) error {
					i.ID = 10
					return nil
				}
				m.getByIDFn = func(ctx context.Context, id int64) (db.Incident, error) {
					return db.Incident{
						ID:          id,
						Title:       "title",
						Description: "desc",
						StatusID:    1,
						ImpactID:    2,
						PriorityID:  3,
					}, nil
				}
			},
		},
		{
			name: "repo error",
			cmd: CreateIncident{
				Title: "title",
			},
			mock: func(m *incidentRepoMock) {
				m.createFn = func(ctx context.Context, i *db.Incident) error {
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
				&impactProviderMock{getByIDFn: func(context.Context, dict.EntityID) (dict.Entity, error) { return dict.Entity{ID: 1}, nil }},
				&statusProviderMock{getByIDFn: func(context.Context, dict.EntityID) (dict.Entity, error) { return dict.Entity{ID: 2}, nil }},
				&priorityProviderMock{getByIDFn: func(context.Context, dict.EntityID) (dict.Entity, error) { return dict.Entity{ID: 3}, nil }},
			)

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
		ImpactID:    2,
		PriorityID:  3,
	}

	tests := []struct {
		name string
		mock func(
			*incidentRepoMock,
			*statusProviderMock,
			*impactProviderMock,
			*priorityProviderMock,
		)
		erroneous bool
	}{
		{
			name: "repo error",
			mock: func(
				repo *incidentRepoMock,
				status *statusProviderMock,
				impact *impactProviderMock,
				priority *priorityProviderMock,
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
				status *statusProviderMock,
				impact *impactProviderMock,
				priority *priorityProviderMock,
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
			name: "impact error",
			mock: func(
				repo *incidentRepoMock,
				status *statusProviderMock,
				impact *impactProviderMock,
				priority *priorityProviderMock,
			) {
				repo.getByIDFn = func(ctx context.Context, id int64) (db.Incident, error) {
					return dbIncident, nil
				}

				status.getByIDFn = func(ctx context.Context, id dict.EntityID) (dict.Entity, error) {
					return dict.Entity{ID: 1}, nil
				}

				impact.getByIDFn = func(ctx context.Context, id dict.EntityID) (dict.Entity, error) {
					return dict.Entity{}, errors.New("impact error")
				}
			},
			erroneous: true,
		},
		{
			name: "priority error",
			mock: func(
				repo *incidentRepoMock,
				status *statusProviderMock,
				impact *impactProviderMock,
				priority *priorityProviderMock,
			) {
				repo.getByIDFn = func(ctx context.Context, id int64) (db.Incident, error) {
					return dbIncident, nil
				}

				status.getByIDFn = func(ctx context.Context, id dict.EntityID) (dict.Entity, error) {
					return dict.Entity{ID: 1}, nil
				}

				impact.getByIDFn = func(ctx context.Context, id dict.EntityID) (dict.Entity, error) {
					return dict.Entity{ID: 2}, nil
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
				status *statusProviderMock,
				impact *impactProviderMock,
				priority *priorityProviderMock,
			) {
				repo.getByIDFn = func(ctx context.Context, id int64) (db.Incident, error) {
					return dbIncident, nil
				}

				status.getByIDFn = func(ctx context.Context, id dict.EntityID) (dict.Entity, error) {
					return dict.Entity{ID: 1}, nil
				}

				impact.getByIDFn = func(ctx context.Context, id dict.EntityID) (dict.Entity, error) {
					return dict.Entity{ID: 2}, nil
				}

				priority.getByIDFn = func(ctx context.Context, id dict.EntityID) (dict.Entity, error) {
					return dict.Entity{ID: 3}, nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &incidentRepoMock{}
			status := &statusProviderMock{}
			impact := &impactProviderMock{}
			priority := &priorityProviderMock{}

			tt.mock(repo, status, impact, priority)

			manager := NewIncidentManager(repo, impact, status, priority)

			res, err := manager.GetByID(ctx, 1)

			if tt.erroneous {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, int64(1), res.ID)
			require.Equal(t, dict.EntityID(1), res.Status.ID)
			require.Equal(t, dict.EntityID(2), res.Impact.ID)
			require.Equal(t, dict.EntityID(3), res.Priority.ID)
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
				ImpactID:    new(int64(2)),
				PriorityID:  new(int64(3)),
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
						ImpactID:    2,
						PriorityID:  3,
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
				&impactProviderMock{getByIDFn: func(context.Context, dict.EntityID) (dict.Entity, error) { return dict.Entity{ID: 1}, nil }},
				&statusProviderMock{getByIDFn: func(context.Context, dict.EntityID) (dict.Entity, error) { return dict.Entity{ID: 2}, nil }},
				&priorityProviderMock{getByIDFn: func(context.Context, dict.EntityID) (dict.Entity, error) { return dict.Entity{ID: 3}, nil }},
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
				&impactProviderMock{getByIDFn: func(context.Context, dict.EntityID) (dict.Entity, error) { return dict.Entity{ID: 1}, nil }},
				&statusProviderMock{getByIDFn: func(context.Context, dict.EntityID) (dict.Entity, error) { return dict.Entity{ID: 2}, nil }},
				&priorityProviderMock{getByIDFn: func(context.Context, dict.EntityID) (dict.Entity, error) { return dict.Entity{ID: 3}, nil }},
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
			ImpactID:    2,
			PriorityID:  3,
		},
		{
			ID:          2,
			Title:       "two",
			Description: "desc2",
			StatusID:    1,
			ImpactID:    2,
			PriorityID:  3,
		},
	}

	tests := []struct {
		name string
		mock func(
			*incidentRepoMock,
			*statusProviderMock,
			*impactProviderMock,
			*priorityProviderMock,
		)
		err bool
	}{
		{
			name: "repo error",
			mock: func(
				repo *incidentRepoMock,
				status *statusProviderMock,
				impact *impactProviderMock,
				priority *priorityProviderMock,
			) {
				repo.listFn = func(context.Context, int, int) ([]db.Incident, error) {
					return nil, errors.New("repo error")
				}
			},
			err: true,
		},
		{
			name: "status list error",
			mock: func(
				repo *incidentRepoMock,
				status *statusProviderMock,
				impact *impactProviderMock,
				priority *priorityProviderMock,
			) {
				repo.listFn = func(context.Context, int, int) ([]db.Incident, error) {
					return dbIncidents, nil
				}

				status.listFn = func(context.Context) ([]dict.Entity, error) {
					return nil, errors.New("status error")
				}
			},
			err: true,
		},
		{
			name: "impact list error",
			mock: func(
				repo *incidentRepoMock,
				status *statusProviderMock,
				impact *impactProviderMock,
				priority *priorityProviderMock,
			) {
				repo.listFn = func(context.Context, int, int) ([]db.Incident, error) {
					return dbIncidents, nil
				}

				status.listFn = func(context.Context) ([]dict.Entity, error) {
					return []dict.Entity{{ID: 1}}, nil
				}

				impact.listFn = func(context.Context) ([]dict.Entity, error) {
					return nil, errors.New("impact error")
				}
			},
			err: true,
		},
		{
			name: "priority list error",
			mock: func(
				repo *incidentRepoMock,
				status *statusProviderMock,
				impact *impactProviderMock,
				priority *priorityProviderMock,
			) {
				repo.listFn = func(context.Context, int, int) ([]db.Incident, error) {
					return dbIncidents, nil
				}

				status.listFn = func(context.Context) ([]dict.Entity, error) {
					return []dict.Entity{{ID: 1}}, nil
				}

				impact.listFn = func(context.Context) ([]dict.Entity, error) {
					return []dict.Entity{{ID: 2}}, nil
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
				status *statusProviderMock,
				impact *impactProviderMock,
				priority *priorityProviderMock,
			) {
				repo.listFn = func(context.Context, int, int) ([]db.Incident, error) {
					return dbIncidents, nil
				}

				status.listFn = func(context.Context) ([]dict.Entity, error) {
					return []dict.Entity{{ID: 1}}, nil
				}

				impact.listFn = func(context.Context) ([]dict.Entity, error) {
					return []dict.Entity{{ID: 2}}, nil
				}

				priority.listFn = func(context.Context) ([]dict.Entity, error) {
					return []dict.Entity{{ID: 3}}, nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &incidentRepoMock{}
			status := &statusProviderMock{}
			impact := &impactProviderMock{}
			priority := &priorityProviderMock{}

			tt.mock(repo, status, impact, priority)

			manager := NewIncidentManager(repo, impact, status, priority)

			res, err := manager.List(ctx, 10, 0)

			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Len(t, res, 2)
			require.Equal(t, dict.EntityID(1), res[0].Status.ID)
			require.Equal(t, dict.EntityID(2), res[0].Impact.ID)
			require.Equal(t, dict.EntityID(3), res[0].Priority.ID)
		})
	}
}

type incidentRepoMock struct {
	createFn  func(context.Context, *db.Incident) error
	getByIDFn func(context.Context, int64) (db.Incident, error)
	updateFn  func(context.Context, *db.Incident) error
	deleteFn  func(context.Context, int64) error
	listFn    func(context.Context, int, int) ([]db.Incident, error)
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

func (m *incidentRepoMock) List(ctx context.Context, limit, offset int) ([]db.Incident, error) {
	return m.listFn(ctx, limit, offset)
}

type impactProviderMock struct {
	getByIDFn func(context.Context, dict.EntityID) (dict.Entity, error)
	listFn    func(context.Context) ([]dict.Entity, error)
}

func (m *impactProviderMock) GetByID(ctx context.Context, id dict.EntityID) (dict.Entity, error) {
	return m.getByIDFn(ctx, id)
}

func (m *impactProviderMock) List(ctx context.Context) ([]dict.Entity, error) {
	if m.listFn != nil {
		return m.listFn(ctx)
	}

	return nil, nil
}

type statusProviderMock struct {
	getByIDFn func(context.Context, dict.EntityID) (dict.Entity, error)
	listFn    func(context.Context) ([]dict.Entity, error)
}

func (m *statusProviderMock) GetByID(ctx context.Context, id dict.EntityID) (dict.Entity, error) {
	return m.getByIDFn(ctx, id)
}

func (m *statusProviderMock) List(ctx context.Context) ([]dict.Entity, error) {
	if m.listFn != nil {
		return m.listFn(ctx)
	}

	return nil, nil
}

type priorityProviderMock struct {
	getByIDFn func(context.Context, dict.EntityID) (dict.Entity, error)
	listFn    func(context.Context) ([]dict.Entity, error)
}

func (m *priorityProviderMock) GetByID(ctx context.Context, id dict.EntityID) (dict.Entity, error) {
	return m.getByIDFn(ctx, id)
}

func (m *priorityProviderMock) List(ctx context.Context) ([]dict.Entity, error) {
	if m.listFn != nil {
		return m.listFn(ctx)
	}

	return nil, nil
}
