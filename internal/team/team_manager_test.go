package team

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/cyradin/fixik/internal/db"
	"github.com/stretchr/testify/require"
)

func TestTeamManager_Create(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		team Team
		mock func(*teamRepoMock)
		err  bool
	}{
		{
			name: "success",
			team: Team{
				Code:        "open",
				Name:        "Open",
				Description: "description",
				Sort:        10,
			},
			mock: func(m *teamRepoMock) {
				m.createFn = func(ctx context.Context, s *db.Team) error {
					s.ID = 1
					return nil
				}
			},
		},
		{
			name: "repo error",
			team: Team{
				Code: "open",
				Name: "Open",
			},
			mock: func(m *teamRepoMock) {
				m.createFn = func(ctx context.Context, s *db.Team) error {
					return errors.New("db error")
				}
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &teamRepoMock{}
			tt.mock(repo)

			m := NewTeamManager(repo, &incidentsCounterMock{}, &usersCounterMock{}, &txExecutorMock{})

			res, err := m.Create(ctx, tt.team)

			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, int64(1), res.ID)
			require.Equal(t, tt.team.Code, res.Code)
			require.Equal(t, tt.team.Name, res.Name)
			require.Equal(t, tt.team.Description, res.Description)
		})
	}
}

func TestTeamManager_GetByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		id   ID
		mock func(*teamRepoMock)
		want Team
		err  bool
	}{
		{
			name: "success",
			id:   1,
			mock: func(m *teamRepoMock) {
				m.getByIDFn = func(ctx context.Context, id int64) (db.Team, error) {
					return db.Team{
						ID:          id,
						Code:        "open",
						Name:        "Open",
						Description: "description",
						Sort:        100,
					}, nil
				}
			},
			want: Team{
				ID:          1,
				Code:        "open",
				Name:        "Open",
				Description: "description",
				Sort:        100,
			},
		},
		{
			name: "repo error",
			id:   1,
			mock: func(m *teamRepoMock) {
				m.getByIDFn = func(ctx context.Context, id int64) (db.Team, error) {
					return db.Team{}, errors.New("db error")
				}
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &teamRepoMock{}
			tt.mock(repo)

			m := NewTeamManager(repo, &incidentsCounterMock{}, &usersCounterMock{}, &txExecutorMock{})

			res, err := m.GetByID(ctx, tt.id)

			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, res)
		})
	}
}

func TestTeamManager_List(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		mock func(*teamRepoMock)
		want []Team
		err  bool
	}{
		{
			name: "success",
			mock: func(m *teamRepoMock) {
				m.listFn = func(ctx context.Context) ([]db.Team, error) {
					return []db.Team{
						{ID: 1, Code: "open", Name: "Open"},
						{ID: 2, Code: "closed", Name: "Closed"},
					}, nil
				}
			},
			want: []Team{
				{ID: 1, Code: "open", Name: "Open"},
				{ID: 2, Code: "closed", Name: "Closed"},
			},
		},
		{
			name: "repo error",
			mock: func(m *teamRepoMock) {
				m.listFn = func(ctx context.Context) ([]db.Team, error) {
					return nil, errors.New("db error")
				}
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &teamRepoMock{}
			tt.mock(repo)

			m := NewTeamManager(repo, &incidentsCounterMock{}, &usersCounterMock{}, &txExecutorMock{})

			res, err := m.List(ctx)

			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, res)
		})
	}
}

func TestTeamManager_Update(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name   string
		status Team
		mock   func(*teamRepoMock)
		err    bool
	}{
		{
			name: "success",
			status: Team{
				ID:   1,
				Code: "closed",
				Name: "Closed",
			},
			mock: func(m *teamRepoMock) {
				m.updateFn = func(ctx context.Context, s *db.Team) error {
					return nil
				}
			},
		},
		{
			name: "repo error",
			status: Team{
				ID:   1,
				Code: "closed",
				Name: "Closed",
			},
			mock: func(m *teamRepoMock) {
				m.updateFn = func(ctx context.Context, s *db.Team) error {
					return errors.New("db error")
				}
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &teamRepoMock{}
			tt.mock(repo)

			m := NewTeamManager(repo, &incidentsCounterMock{}, &usersCounterMock{}, &txExecutorMock{})

			res, err := m.Update(ctx, tt.status)

			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.status, res)
		})
	}
}

func TestTeamManager_Delete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	expectedErr := fmt.Errorf("error")

	tests := []struct {
		name string
		id   ID
		mock func(*teamRepoMock, *incidentsCounterMock, *usersCounterMock, *txExecutorMock)
		err  error
	}{
		{
			name: "success",
			id:   1,
			mock: func(r *teamRepoMock, c *incidentsCounterMock, u *usersCounterMock, tx *txExecutorMock) {
				c.countFn = func(ctx context.Context, id int64) (int, error) {
					return 0, nil
				}
				u.countFn = func(ctx context.Context, id int64) (int, error) {
					return 0, nil
				}

				r.deleteFn = func(ctx context.Context, id int64) error {
					return nil
				}
			},
		},
		{
			name: "has dependencies (incidents)",
			id:   1,
			mock: func(r *teamRepoMock, c *incidentsCounterMock, u *usersCounterMock, tx *txExecutorMock) {
				c.countFn = func(ctx context.Context, id int64) (int, error) {
					return 5, nil
				}
				u.countFn = func(ctx context.Context, id int64) (int, error) {
					return 0, nil
				}
			},
			err: ErrHasDependantIncidents,
		},
		{
			name: "has dependencies (users)",
			id:   1,
			mock: func(r *teamRepoMock, c *incidentsCounterMock, u *usersCounterMock, tx *txExecutorMock) {
				c.countFn = func(ctx context.Context, id int64) (int, error) {
					return 0, nil
				}
				u.countFn = func(ctx context.Context, id int64) (int, error) {
					return 3, nil
				}
			},
			err: ErrHasDependantUsers,
		},
		{
			name: "count error (incidents)",
			id:   1,
			mock: func(r *teamRepoMock, c *incidentsCounterMock, u *usersCounterMock, tx *txExecutorMock) {
				c.countFn = func(ctx context.Context, id int64) (int, error) {
					return 0, expectedErr
				}
			},
			err: expectedErr,
		},
		{
			name: "repo error",
			id:   1,
			mock: func(r *teamRepoMock, c *incidentsCounterMock, u *usersCounterMock, tx *txExecutorMock) {
				c.countFn = func(ctx context.Context, id int64) (int, error) {
					return 0, nil
				}
				u.countFn = func(ctx context.Context, id int64) (int, error) {
					return 0, nil
				}
				r.deleteFn = func(ctx context.Context, id int64) error {
					return expectedErr
				}
			},
			err: expectedErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &teamRepoMock{}
			counter := &incidentsCounterMock{}
			usersCounter := &usersCounterMock{}
			tx := &txExecutorMock{}

			tt.mock(repo, counter, usersCounter, tx)

			m := NewTeamManager(repo, counter, usersCounter, tx)

			err := m.Delete(ctx, tt.id)
			require.ErrorIs(t, err, tt.err)
		})
	}
}

type teamRepoMock struct {
	createFn  func(ctx context.Context, s *db.Team) error
	getByIDFn func(ctx context.Context, id int64) (db.Team, error)
	listFn    func(ctx context.Context) ([]db.Team, error)
	updateFn  func(ctx context.Context, s *db.Team) error
	deleteFn  func(ctx context.Context, id int64) error
}

func (m *teamRepoMock) Create(ctx context.Context, s *db.Team) error {
	return m.createFn(ctx, s)
}

func (m *teamRepoMock) GetByID(ctx context.Context, id int64) (db.Team, error) {
	return m.getByIDFn(ctx, id)
}

func (m *teamRepoMock) List(ctx context.Context) ([]db.Team, error) {
	return m.listFn(ctx)
}

func (m *teamRepoMock) Update(ctx context.Context, s *db.Team) error {
	return m.updateFn(ctx, s)
}

func (m *teamRepoMock) Delete(ctx context.Context, id int64) error {
	return m.deleteFn(ctx, id)
}

type incidentsCounterMock struct {
	countFn func(ctx context.Context, id int64) (int, error)
}

func (m *incidentsCounterMock) CountByTeam(ctx context.Context, id int64) (int, error) {
	return m.countFn(ctx, id)
}

type txExecutorMock struct {
	execFn func(ctx context.Context, fn func(ctx context.Context) error) error
}

func (m *txExecutorMock) Exec(ctx context.Context, fn func(ctx context.Context) error) error {
	if m.execFn != nil {
		return m.execFn(ctx, fn)
	}

	return fn(ctx)
}

type usersCounterMock struct {
	countFn func(ctx context.Context, teamID int64) (int, error)
}

func (m *usersCounterMock) CountByTeam(ctx context.Context, teamID int64) (int, error) {
	if m.countFn != nil {
		return m.countFn(ctx, teamID)
	}

	return 0, nil
}
