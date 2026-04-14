package priority

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/cyradin/fixik/internal/db"
	"github.com/stretchr/testify/require"
)

func TestPriorityManager_Create(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name   string
		entity Priority
		mock   func(*priorityRepoMock)
		err    bool
	}{
		{
			name: "success",
			entity: Priority{
				Code:        "open",
				Name:        "Open",
				Description: "description",
				Sort:        10,
			},
			mock: func(m *priorityRepoMock) {
				m.createFn = func(ctx context.Context, s *db.DictEntity) error {
					s.ID = 1
					return nil
				}
			},
		},
		{
			name: "repo error",
			entity: Priority{
				Code: "open",
				Name: "Open",
			},
			mock: func(m *priorityRepoMock) {
				m.createFn = func(ctx context.Context, s *db.DictEntity) error {
					return errors.New("db error")
				}
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &priorityRepoMock{}
			tt.mock(repo)

			m := NewPriorityManager(repo, &incidentsCounterMock{}, &txExecutorMock{})

			res, err := m.Create(ctx, tt.entity)

			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, int64(1), res.ID)
			require.Equal(t, tt.entity.Code, res.Code)
			require.Equal(t, tt.entity.Name, res.Name)
			require.Equal(t, tt.entity.Description, res.Description)
		})
	}
}

func TestPriorityManager_GetByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		id   ID
		mock func(*priorityRepoMock)
		want Priority
		err  bool
	}{
		{
			name: "success",
			id:   1,
			mock: func(m *priorityRepoMock) {
				m.getByIDFn = func(ctx context.Context, id int64) (db.DictEntity, error) {
					return db.DictEntity{
						ID:          id,
						Code:        "open",
						Name:        "Open",
						Description: "description",
						Sort:        100,
					}, nil
				}
			},
			want: Priority{
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
			mock: func(m *priorityRepoMock) {
				m.getByIDFn = func(ctx context.Context, id int64) (db.DictEntity, error) {
					return db.DictEntity{}, errors.New("db error")
				}
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &priorityRepoMock{}
			tt.mock(repo)

			m := NewPriorityManager(repo, &incidentsCounterMock{}, &txExecutorMock{})

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

func TestPriorityManager_List(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		mock func(*priorityRepoMock)
		want []Priority
		err  bool
	}{
		{
			name: "success",
			mock: func(m *priorityRepoMock) {
				m.listFn = func(ctx context.Context) ([]db.DictEntity, error) {
					return []db.DictEntity{
						{ID: 1, Code: "open", Name: "Open"},
						{ID: 2, Code: "closed", Name: "Closed"},
					}, nil
				}
			},
			want: []Priority{
				{ID: 1, Code: "open", Name: "Open"},
				{ID: 2, Code: "closed", Name: "Closed"},
			},
		},
		{
			name: "repo error",
			mock: func(m *priorityRepoMock) {
				m.listFn = func(ctx context.Context) ([]db.DictEntity, error) {
					return nil, errors.New("db error")
				}
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &priorityRepoMock{}
			tt.mock(repo)

			m := NewPriorityManager(repo, &incidentsCounterMock{}, &txExecutorMock{})

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

func TestPriorityManager_Update(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name   string
		status Priority
		mock   func(*priorityRepoMock)
		err    bool
	}{
		{
			name: "success",
			status: Priority{
				ID:   1,
				Code: "closed",
				Name: "Closed",
			},
			mock: func(m *priorityRepoMock) {
				m.updateFn = func(ctx context.Context, s *db.DictEntity) error {
					return nil
				}
			},
		},
		{
			name: "repo error",
			status: Priority{
				ID:   1,
				Code: "closed",
				Name: "Closed",
			},
			mock: func(m *priorityRepoMock) {
				m.updateFn = func(ctx context.Context, s *db.DictEntity) error {
					return errors.New("db error")
				}
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &priorityRepoMock{}
			tt.mock(repo)

			m := NewPriorityManager(repo, &incidentsCounterMock{}, &txExecutorMock{})

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

func TestPriorityManager_Delete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		id   ID
		mock func(*priorityRepoMock, *incidentsCounterMock, *txExecutorMock)
		err  error
	}{
		{
			name: "success",
			id:   1,
			mock: func(r *priorityRepoMock, c *incidentsCounterMock, tx *txExecutorMock) {
				c.countFn = func(ctx context.Context, id int64) (int, error) {
					return 0, nil
				}

				r.deleteFn = func(ctx context.Context, id int64) error {
					return nil
				}
			},
		},
		{
			name: "has dependencies",
			id:   1,
			mock: func(r *priorityRepoMock, c *incidentsCounterMock, tx *txExecutorMock) {
				c.countFn = func(ctx context.Context, id int64) (int, error) {
					return 10, nil
				}
			},
			err: ErrHasDependantEntities,
		},
		{
			name: "count error",
			id:   1,
			mock: func(r *priorityRepoMock, c *incidentsCounterMock, tx *txExecutorMock) {
				c.countFn = func(ctx context.Context, id int64) (int, error) {
					return 0, errors.New("db error")
				}
			},
			err: fmt.Errorf("count incidents"),
		},
		{
			name: "repo error",
			id:   1,
			mock: func(r *priorityRepoMock, c *incidentsCounterMock, tx *txExecutorMock) {
				c.countFn = func(ctx context.Context, id int64) (int, error) {
					return 0, nil
				}

				r.deleteFn = func(ctx context.Context, id int64) error {
					return errors.New("db error")
				}
			},
			err: fmt.Errorf("repository delete"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &priorityRepoMock{}
			counter := &incidentsCounterMock{}
			tx := &txExecutorMock{}

			tt.mock(repo, counter, tx)

			m := NewPriorityManager(repo, counter, tx)

			err := m.Delete(ctx, tt.id)

			if tt.err != nil {
				require.Error(t, err)

				if errors.Is(tt.err, ErrHasDependantEntities) {
					require.ErrorIs(t, err, ErrHasDependantEntities)
				}

				return
			}

			require.NoError(t, err)
		})
	}
}

type priorityRepoMock struct {
	createFn  func(ctx context.Context, s *db.DictEntity) error
	getByIDFn func(ctx context.Context, id int64) (db.DictEntity, error)
	listFn    func(ctx context.Context) ([]db.DictEntity, error)
	updateFn  func(ctx context.Context, s *db.DictEntity) error
	deleteFn  func(ctx context.Context, id int64) error
}

func (m *priorityRepoMock) Create(ctx context.Context, s *db.DictEntity) error {
	return m.createFn(ctx, s)
}

func (m *priorityRepoMock) GetByID(ctx context.Context, id int64) (db.DictEntity, error) {
	return m.getByIDFn(ctx, id)
}

func (m *priorityRepoMock) List(ctx context.Context) ([]db.DictEntity, error) {
	return m.listFn(ctx)
}

func (m *priorityRepoMock) Update(ctx context.Context, s *db.DictEntity) error {
	return m.updateFn(ctx, s)
}

func (m *priorityRepoMock) Delete(ctx context.Context, id int64) error {
	return m.deleteFn(ctx, id)
}

type incidentsCounterMock struct {
	countFn func(ctx context.Context, id int64) (int, error)
}

func (m *incidentsCounterMock) CountByPriority(ctx context.Context, id int64) (int, error) {
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
