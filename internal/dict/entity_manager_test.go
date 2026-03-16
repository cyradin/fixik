package dict

import (
	"context"
	"errors"
	"testing"

	"github.com/cyradin/fixik/internal/db"
	"github.com/stretchr/testify/require"
)

func TestEntityManager_Create(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name   string
		entity Entity
		mock   func(*statusRepoMock)
		err    bool
	}{
		{
			name: "success",
			entity: Entity{
				Code:        "open",
				Name:        "Open",
				Description: "description",
			},
			mock: func(m *statusRepoMock) {
				m.createFn = func(ctx context.Context, s *db.DictEntity) error {
					s.ID = 1
					return nil
				}
			},
		},
		{
			name: "repo error",
			entity: Entity{
				Code: "open",
				Name: "Open",
			},
			mock: func(m *statusRepoMock) {
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

			repo := &statusRepoMock{}
			tt.mock(repo)

			m := newEntityManager(repo)

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

func TestEntityManager_GetByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		id   EntityID
		mock func(*statusRepoMock)
		want Entity
		err  bool
	}{
		{
			name: "success",
			id:   1,
			mock: func(m *statusRepoMock) {
				m.getByIDFn = func(ctx context.Context, id int64) (db.DictEntity, error) {
					return db.DictEntity{
						ID:          id,
						Code:        "open",
						Name:        "Open",
						Description: "description",
					}, nil
				}
			},
			want: Entity{
				ID:          1,
				Code:        "open",
				Name:        "Open",
				Description: "description",
			},
		},
		{
			name: "repo error",
			id:   1,
			mock: func(m *statusRepoMock) {
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

			repo := &statusRepoMock{}
			tt.mock(repo)

			m := newEntityManager(repo)

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

func TestEntityManager_List(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		mock func(*statusRepoMock)
		want []Entity
		err  bool
	}{
		{
			name: "success",
			mock: func(m *statusRepoMock) {
				m.listFn = func(ctx context.Context) ([]db.DictEntity, error) {
					return []db.DictEntity{
						{ID: 1, Code: "open", Name: "Open"},
						{ID: 2, Code: "closed", Name: "Closed"},
					}, nil
				}
			},
			want: []Entity{
				{ID: 1, Code: "open", Name: "Open"},
				{ID: 2, Code: "closed", Name: "Closed"},
			},
		},
		{
			name: "repo error",
			mock: func(m *statusRepoMock) {
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

			repo := &statusRepoMock{}
			tt.mock(repo)

			m := newEntityManager(repo)

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

func TestEntityManager_Update(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name   string
		status Entity
		mock   func(*statusRepoMock)
		err    bool
	}{
		{
			name: "success",
			status: Entity{
				ID:   1,
				Code: "closed",
				Name: "Closed",
			},
			mock: func(m *statusRepoMock) {
				m.updateFn = func(ctx context.Context, s *db.DictEntity) error {
					return nil
				}
			},
		},
		{
			name: "repo error",
			status: Entity{
				ID:   1,
				Code: "closed",
				Name: "Closed",
			},
			mock: func(m *statusRepoMock) {
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

			repo := &statusRepoMock{}
			tt.mock(repo)

			m := newEntityManager(repo)

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

func TestEntityManager_Delete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		id   EntityID
		mock func(*statusRepoMock)
		err  bool
	}{
		{
			name: "success",
			id:   1,
			mock: func(m *statusRepoMock) {
				m.deleteFn = func(ctx context.Context, id int64) error {
					return nil
				}
			},
		},
		{
			name: "repo error",
			id:   1,
			mock: func(m *statusRepoMock) {
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

			repo := &statusRepoMock{}
			tt.mock(repo)

			m := newEntityManager(repo)

			err := m.Delete(ctx, tt.id)

			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

type statusRepoMock struct {
	createFn  func(ctx context.Context, s *db.DictEntity) error
	getByIDFn func(ctx context.Context, id int64) (db.DictEntity, error)
	listFn    func(ctx context.Context) ([]db.DictEntity, error)
	updateFn  func(ctx context.Context, s *db.DictEntity) error
	deleteFn  func(ctx context.Context, id int64) error
}

func (m *statusRepoMock) Create(ctx context.Context, s *db.DictEntity) error {
	return m.createFn(ctx, s)
}

func (m *statusRepoMock) GetByID(ctx context.Context, id int64) (db.DictEntity, error) {
	return m.getByIDFn(ctx, id)
}

func (m *statusRepoMock) List(ctx context.Context) ([]db.DictEntity, error) {
	return m.listFn(ctx)
}

func (m *statusRepoMock) Update(ctx context.Context, s *db.DictEntity) error {
	return m.updateFn(ctx, s)
}

func (m *statusRepoMock) Delete(ctx context.Context, id int64) error {
	return m.deleteFn(ctx, id)
}
