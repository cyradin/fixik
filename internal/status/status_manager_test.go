package status

import (
	"context"
	"errors"
	"testing"

	"github.com/cyradin/fixik/internal/db"
	"github.com/stretchr/testify/require"
)

func TestStatusManager_Create(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name   string
		status Status
		mock   func(*statusRepoMock)
		err    bool
	}{
		{
			name: "success",
			status: Status{
				Code:        "done",
				Name:        "Done",
				Description: "Completed status",
				Sort:        100,
				IsFinal:     true,
			},
			mock: func(m *statusRepoMock) {
				m.createFn = func(ctx context.Context, s *db.Status) error {
					s.ID = 1
					return nil
				}
			},
		},
		{
			name: "repo error",
			status: Status{
				Code: "done",
				Name: "Done",
			},
			mock: func(m *statusRepoMock) {
				m.createFn = func(ctx context.Context, s *db.Status) error {
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

			m := NewStatusManager(repo)

			res, err := m.Create(ctx, tt.status)

			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, int64(1), res.ID)
			require.Equal(t, tt.status.Code, res.Code)
			require.Equal(t, tt.status.Name, res.Name)
			require.Equal(t, tt.status.Description, res.Description)
			require.Equal(t, tt.status.Sort, res.Sort)
			require.Equal(t, tt.status.IsFinal, res.IsFinal)
		})
	}
}

func TestStatusManager_GetByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		id   StatusID
		mock func(*statusRepoMock)
		want Status
		err  bool
	}{
		{
			name: "success",
			id:   1,
			mock: func(m *statusRepoMock) {
				m.getByIDFn = func(ctx context.Context, id int64) (db.Status, error) {
					return db.Status{
						ID:          id,
						Code:        "done",
						Name:        "Done",
						Description: "Completed status",
						Sort:        100,
						IsFinal:     true,
					}, nil
				}
			},
			want: Status{
				ID:          1,
				Code:        "done",
				Name:        "Done",
				Description: "Completed status",
				Sort:        100,
				IsFinal:     true,
			},
		},
		{
			name: "repo error",
			id:   1,
			mock: func(m *statusRepoMock) {
				m.getByIDFn = func(ctx context.Context, id int64) (db.Status, error) {
					return db.Status{}, errors.New("db error")
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

			m := NewStatusManager(repo)

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

func TestStatusManager_List(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		mock func(*statusRepoMock)
		want []Status
		err  bool
	}{
		{
			name: "success",
			mock: func(m *statusRepoMock) {
				m.listFn = func(ctx context.Context) ([]db.Status, error) {
					return []db.Status{
						{ID: 1, Code: "done", Name: "Done", IsFinal: true},
						{ID: 2, Code: "in_progress", Name: "In Progress", IsFinal: false},
					}, nil
				}
			},
			want: []Status{
				{ID: 1, Code: "done", Name: "Done", IsFinal: true},
				{ID: 2, Code: "in_progress", Name: "In Progress", IsFinal: false},
			},
		},
		{
			name: "repo error",
			mock: func(m *statusRepoMock) {
				m.listFn = func(ctx context.Context) ([]db.Status, error) {
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

			m := NewStatusManager(repo)

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

func TestStatusManager_Update(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name   string
		status Status
		mock   func(*statusRepoMock)
		err    bool
	}{
		{
			name: "success",
			status: Status{
				ID:      1,
				Code:    "done",
				Name:    "Done",
				IsFinal: true,
			},
			mock: func(m *statusRepoMock) {
				m.updateFn = func(ctx context.Context, s *db.Status) error {
					return nil
				}
			},
		},
		{
			name: "repo error",
			status: Status{
				ID:   1,
				Code: "done",
				Name: "Done",
			},
			mock: func(m *statusRepoMock) {
				m.updateFn = func(ctx context.Context, s *db.Status) error {
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

			m := NewStatusManager(repo)

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

func TestStatusManager_Delete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		id   StatusID
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

			m := NewStatusManager(repo)

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
	createFn  func(ctx context.Context, s *db.Status) error
	getByIDFn func(ctx context.Context, id int64) (db.Status, error)
	listFn    func(ctx context.Context) ([]db.Status, error)
	updateFn  func(ctx context.Context, s *db.Status) error
	deleteFn  func(ctx context.Context, id int64) error
}

func (m *statusRepoMock) Create(ctx context.Context, s *db.Status) error {
	return m.createFn(ctx, s)
}

func (m *statusRepoMock) GetByID(ctx context.Context, id int64) (db.Status, error) {
	return m.getByIDFn(ctx, id)
}

func (m *statusRepoMock) List(ctx context.Context) ([]db.Status, error) {
	return m.listFn(ctx)
}

func (m *statusRepoMock) Update(ctx context.Context, s *db.Status) error {
	return m.updateFn(ctx, s)
}

func (m *statusRepoMock) Delete(ctx context.Context, id int64) error {
	return m.deleteFn(ctx, id)
}
