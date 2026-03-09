package incident

import (
	"context"
	"errors"
	"testing"

	"github.com/cyradin/fixik/internal/db"
	"github.com/stretchr/testify/require"
)

func TestPriorityManager_Create(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name     string
		priority Priority
		mock     func(*priorityRepoMock)
		wantErr  bool
	}{
		{
			name: "success",
			priority: Priority{
				Code: "high",
				Name: "High",
			},
			mock: func(m *priorityRepoMock) {
				m.createFn = func(ctx context.Context, p *db.Priority) error {
					p.ID = 1
					return nil
				}
			},
		},
		{
			name: "repo error",
			priority: Priority{
				Code: "high",
				Name: "High",
			},
			mock: func(m *priorityRepoMock) {
				m.createFn = func(ctx context.Context, p *db.Priority) error {
					return errors.New("db error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &priorityRepoMock{}
			tt.mock(repo)

			m := NewPriorityManager(repo)

			res, err := m.Create(ctx, tt.priority)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, int64(1), res.ID)
			require.Equal(t, tt.priority.Code, res.Code)
			require.Equal(t, tt.priority.Name, res.Name)
		})
	}
}

func TestPriorityManager_GetByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name    string
		id      PriorityID
		mock    func(*priorityRepoMock)
		want    Priority
		wantErr bool
	}{
		{
			name: "success",
			id:   1,
			mock: func(m *priorityRepoMock) {
				m.getByIDFn = func(ctx context.Context, id int64) (db.Priority, error) {
					return db.Priority{
						ID:   id,
						Code: "high",
						Name: "High",
					}, nil
				}
			},
			want: Priority{
				ID:   1,
				Code: "high",
				Name: "High",
			},
		},
		{
			name: "repo error",
			id:   1,
			mock: func(m *priorityRepoMock) {
				m.getByIDFn = func(ctx context.Context, id int64) (db.Priority, error) {
					return db.Priority{}, errors.New("db error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &priorityRepoMock{}
			tt.mock(repo)

			m := NewPriorityManager(repo)

			res, err := m.GetByID(ctx, tt.id)

			if tt.wantErr {
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
		name    string
		mock    func(*priorityRepoMock)
		want    []Priority
		wantErr bool
	}{
		{
			name: "success",
			mock: func(m *priorityRepoMock) {
				m.listFn = func(ctx context.Context) ([]db.Priority, error) {
					return []db.Priority{
						{ID: 1, Code: "low", Name: "Low"},
						{ID: 2, Code: "high", Name: "High"},
					}, nil
				}
			},
			want: []Priority{
				{ID: 1, Code: "low", Name: "Low"},
				{ID: 2, Code: "high", Name: "High"},
			},
		},
		{
			name: "repo error",
			mock: func(m *priorityRepoMock) {
				m.listFn = func(ctx context.Context) ([]db.Priority, error) {
					return nil, errors.New("db error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &priorityRepoMock{}
			tt.mock(repo)

			m := NewPriorityManager(repo)

			res, err := m.List(ctx)

			if tt.wantErr {
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
		name     string
		priority Priority
		mock     func(*priorityRepoMock)
		wantErr  bool
	}{
		{
			name: "success",
			priority: Priority{
				ID:   1,
				Code: "medium",
				Name: "Medium",
			},
			mock: func(m *priorityRepoMock) {
				m.updateFn = func(ctx context.Context, p *db.Priority) error {
					return nil
				}
			},
		},
		{
			name: "repo error",
			priority: Priority{
				ID:   1,
				Code: "medium",
				Name: "Medium",
			},
			mock: func(m *priorityRepoMock) {
				m.updateFn = func(ctx context.Context, p *db.Priority) error {
					return errors.New("db error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &priorityRepoMock{}
			tt.mock(repo)

			m := NewPriorityManager(repo)

			res, err := m.Update(ctx, tt.priority)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.priority, res)
		})
	}
}

func TestPriorityManager_Delete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name    string
		id      PriorityID
		mock    func(*priorityRepoMock)
		wantErr bool
	}{
		{
			name: "success",
			id:   1,
			mock: func(m *priorityRepoMock) {
				m.deleteFn = func(ctx context.Context, id int64) error {
					return nil
				}
			},
		},
		{
			name: "repo error",
			id:   1,
			mock: func(m *priorityRepoMock) {
				m.deleteFn = func(ctx context.Context, id int64) error {
					return errors.New("db error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &priorityRepoMock{}
			tt.mock(repo)

			m := NewPriorityManager(repo)

			err := m.Delete(ctx, tt.id)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

type priorityRepoMock struct {
	createFn  func(ctx context.Context, s *db.Priority) error
	getByIDFn func(ctx context.Context, id int64) (db.Priority, error)
	listFn    func(ctx context.Context) ([]db.Priority, error)
	updateFn  func(ctx context.Context, s *db.Priority) error
	deleteFn  func(ctx context.Context, id int64) error
}

func (m *priorityRepoMock) Create(ctx context.Context, s *db.Priority) error {
	return m.createFn(ctx, s)
}

func (m *priorityRepoMock) GetByID(ctx context.Context, id int64) (db.Priority, error) {
	return m.getByIDFn(ctx, id)
}

func (m *priorityRepoMock) List(ctx context.Context) ([]db.Priority, error) {
	return m.listFn(ctx)
}

func (m *priorityRepoMock) Update(ctx context.Context, s *db.Priority) error {
	return m.updateFn(ctx, s)
}

func (m *priorityRepoMock) Delete(ctx context.Context, id int64) error {
	return m.deleteFn(ctx, id)
}
