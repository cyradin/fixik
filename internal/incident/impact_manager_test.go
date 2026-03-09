package incident

import (
	"context"
	"errors"
	"testing"

	"github.com/cyradin/fixik/internal/db"
	"github.com/stretchr/testify/require"
)

func TestImpactManager_Create(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name    string
		impact  Impact
		mock    func(*impactRepoMock)
		wantErr bool
	}{
		{
			name: "success",
			impact: Impact{
				Code: "high",
				Name: "High",
			},
			mock: func(m *impactRepoMock) {
				m.createFn = func(ctx context.Context, s *db.Impact) error {
					s.ID = 1
					return nil
				}
			},
		},
		{
			name: "repo error",
			impact: Impact{
				Code: "high",
				Name: "High",
			},
			mock: func(m *impactRepoMock) {
				m.createFn = func(ctx context.Context, s *db.Impact) error {
					return errors.New("db error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &impactRepoMock{}
			tt.mock(repo)

			m := NewImpactManager(repo)

			res, err := m.Create(ctx, tt.impact)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, int64(1), res.ID)
			require.Equal(t, tt.impact.Code, res.Code)
			require.Equal(t, tt.impact.Name, res.Name)
		})
	}
}

func TestImpactManager_GetByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name    string
		id      ImpactID
		mock    func(*impactRepoMock)
		want    Impact
		wantErr bool
	}{
		{
			name: "success",
			id:   1,
			mock: func(m *impactRepoMock) {
				m.getByIDFn = func(ctx context.Context, id int64) (db.Impact, error) {
					return db.Impact{
						ID:   id,
						Code: "high",
						Name: "High",
					}, nil
				}
			},
			want: Impact{
				ID:   1,
				Code: "high",
				Name: "High",
			},
		},
		{
			name: "repo error",
			id:   1,
			mock: func(m *impactRepoMock) {
				m.getByIDFn = func(ctx context.Context, id int64) (db.Impact, error) {
					return db.Impact{}, errors.New("db error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &impactRepoMock{}
			tt.mock(repo)

			m := NewImpactManager(repo)

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

func TestImpactManager_List(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name    string
		mock    func(*impactRepoMock)
		want    []Impact
		wantErr bool
	}{
		{
			name: "success",
			mock: func(m *impactRepoMock) {
				m.listFn = func(ctx context.Context) ([]db.Impact, error) {
					return []db.Impact{
						{ID: 1, Code: "low", Name: "Low"},
						{ID: 2, Code: "high", Name: "High"},
					}, nil
				}
			},
			want: []Impact{
				{ID: 1, Code: "low", Name: "Low"},
				{ID: 2, Code: "high", Name: "High"},
			},
		},
		{
			name: "repo error",
			mock: func(m *impactRepoMock) {
				m.listFn = func(ctx context.Context) ([]db.Impact, error) {
					return nil, errors.New("db error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &impactRepoMock{}
			tt.mock(repo)

			m := NewImpactManager(repo)

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

func TestImpactManager_Update(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name    string
		impact  Impact
		mock    func(*impactRepoMock)
		wantErr bool
	}{
		{
			name: "success",
			impact: Impact{
				ID:   1,
				Code: "medium",
				Name: "Medium",
			},
			mock: func(m *impactRepoMock) {
				m.updateFn = func(ctx context.Context, s *db.Impact) error {
					return nil
				}
			},
		},
		{
			name: "repo error",
			impact: Impact{
				ID:   1,
				Code: "medium",
				Name: "Medium",
			},
			mock: func(m *impactRepoMock) {
				m.updateFn = func(ctx context.Context, s *db.Impact) error {
					return errors.New("db error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &impactRepoMock{}
			tt.mock(repo)

			m := NewImpactManager(repo)

			res, err := m.Update(ctx, tt.impact)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.impact, res)
		})
	}
}

func TestImpactManager_Delete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name    string
		id      ImpactID
		mock    func(*impactRepoMock)
		wantErr bool
	}{
		{
			name: "success",
			id:   1,
			mock: func(m *impactRepoMock) {
				m.deleteFn = func(ctx context.Context, id int64) error {
					return nil
				}
			},
		},
		{
			name: "repo error",
			id:   1,
			mock: func(m *impactRepoMock) {
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

			repo := &impactRepoMock{}
			tt.mock(repo)

			m := NewImpactManager(repo)

			err := m.Delete(ctx, tt.id)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

type impactRepoMock struct {
	createFn  func(ctx context.Context, s *db.Impact) error
	getByIDFn func(ctx context.Context, id int64) (db.Impact, error)
	listFn    func(ctx context.Context) ([]db.Impact, error)
	updateFn  func(ctx context.Context, s *db.Impact) error
	deleteFn  func(ctx context.Context, id int64) error
}

func (m *impactRepoMock) Create(ctx context.Context, s *db.Impact) error {
	return m.createFn(ctx, s)
}

func (m *impactRepoMock) GetByID(ctx context.Context, id int64) (db.Impact, error) {
	return m.getByIDFn(ctx, id)
}

func (m *impactRepoMock) List(ctx context.Context) ([]db.Impact, error) {
	return m.listFn(ctx)
}

func (m *impactRepoMock) Update(ctx context.Context, s *db.Impact) error {
	return m.updateFn(ctx, s)
}

func (m *impactRepoMock) Delete(ctx context.Context, id int64) error {
	return m.deleteFn(ctx, id)
}
