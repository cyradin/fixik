package user

import (
	"context"
	"errors"
	"testing"

	"github.com/cyradin/fixik/internal/db"
	"github.com/cyradin/fixik/internal/dict"
	"github.com/stretchr/testify/require"
)

func TestUserManager_Create(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name string
		cmd  CreateUser
		mock func(*userRepoMock)
		err  bool
	}{
		{
			name: "success",
			cmd: CreateUser{
				Username: "alice",
				Email:    "alice@example.com",
				Password: "pass",
				TeamID:   1,
				RoleIDs:  []int64{1, 2},
			},
			mock: func(m *userRepoMock) {
				m.createFn = func(ctx context.Context, u *db.User) error {
					u.ID = 10
					return nil
				}
				m.getByIDFn = func(ctx context.Context, id int64) (db.User, error) {
					return db.User{
						ID:       id,
						Username: "alice",
						Email:    "alice@example.com",
						Password: "hashed",
						TeamID:   1,
						RoleIDs:  []int64{1, 2},
					}, nil
				}
			},
		},
		{
			name: "repo error",
			cmd: CreateUser{
				Username: "bob",
			},
			mock: func(m *userRepoMock) {
				m.createFn = func(ctx context.Context, u *db.User) error {
					return errors.New("db error")
				}
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &userRepoMock{}
			tt.mock(repo)

			manager := NewUserManager(repo, &roleProviderMock{
				getByIDFn: func(ctx context.Context, id dict.EntityID) (dict.Entity, error) {
					return dict.Entity{ID: id}, nil
				},
			})

			u, err := manager.Create(ctx, tt.cmd)
			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.cmd.Username, u.Username)
			require.Len(t, u.RoleIDs, len(tt.cmd.RoleIDs))
		})
	}
}

func TestUserManager_GetByID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	dbUser := db.User{
		ID:       1,
		Username: "alice",
		Email:    "alice@example.com",
		Password: "hashed",
		TeamID:   1,
		RoleIDs:  []int64{1, 2},
	}

	tests := []struct {
		name string
		mock func(*userRepoMock, *roleProviderMock)
		err  bool
	}{
		{
			name: "repo error",
			mock: func(repo *userRepoMock, role *roleProviderMock) {
				repo.getByIDFn = func(ctx context.Context, id int64) (db.User, error) {
					return db.User{}, errors.New("repo error")
				}
			},
			err: true,
		},
		{
			name: "role error",
			mock: func(repo *userRepoMock, role *roleProviderMock) {
				repo.getByIDFn = func(ctx context.Context, id int64) (db.User, error) {
					return dbUser, nil
				}
				role.getByIDFn = func(ctx context.Context, id dict.EntityID) (dict.Entity, error) {
					if id == 2 {
						return dict.Entity{}, errors.New("role error")
					}
					return dict.Entity{ID: id}, nil
				}
			},
			err: true,
		},
		{
			name: "success",
			mock: func(repo *userRepoMock, role *roleProviderMock) {
				repo.getByIDFn = func(ctx context.Context, id int64) (db.User, error) {
					return dbUser, nil
				}
				role.getByIDFn = func(ctx context.Context, id dict.EntityID) (dict.Entity, error) {
					return dict.Entity{ID: id}, nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &userRepoMock{}
			role := &roleProviderMock{}
			tt.mock(repo, role)

			manager := NewUserManager(repo, role)
			u, err := manager.GetByID(ctx, 1)

			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, dbUser.ID, u.ID)
			require.Len(t, u.RoleIDs, len(dbUser.RoleIDs))
		})
	}
}

func TestUserManager_Update(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name string
		cmd  UpdateUser
		mock func(*userRepoMock)
		err  bool
	}{
		{
			name: "success",
			cmd: UpdateUser{
				ID:       1,
				Username: new("new"),
				Email:    new("new@example.com"),
				Password: new("newpass"),
				TeamID:   new(int64(1)),
				RoleIDs:  &[]int64{2, 3},
			},
			mock: func(m *userRepoMock) {
				m.getByIDFn = func(ctx context.Context, id int64) (db.User, error) {
					return db.User{ID: id, RoleIDs: []int64{1}}, nil
				}
				m.updateFn = func(ctx context.Context, u *db.User) error {
					return nil
				}
			},
		},
		{
			name: "repo error",
			cmd:  UpdateUser{ID: 1},
			mock: func(m *userRepoMock) {
				m.getByIDFn = func(ctx context.Context, id int64) (db.User, error) {
					return db.User{}, errors.New("repo error")
				}
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &userRepoMock{}
			tt.mock(repo)

			manager := NewUserManager(repo, &roleProviderMock{
				getByIDFn: func(ctx context.Context, id dict.EntityID) (dict.Entity, error) {
					return dict.Entity{ID: id}, nil
				},
			})

			_, err := manager.Update(ctx, tt.cmd)
			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestUserManager_Delete(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name string
		id   int64
		mock func(*userRepoMock)
		err  bool
	}{
		{
			name: "success",
			id:   1,
			mock: func(m *userRepoMock) {
				m.deleteFn = func(ctx context.Context, id int64) error {
					return nil
				}
			},
		},
		{
			name: "repo error",
			id:   1,
			mock: func(m *userRepoMock) {
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
			repo := &userRepoMock{}
			tt.mock(repo)

			manager := NewUserManager(repo, &roleProviderMock{})
			err := manager.Delete(ctx, tt.id)

			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

// Mocks
type userRepoMock struct {
	createFn  func(context.Context, *db.User) error
	getByIDFn func(context.Context, int64) (db.User, error)
	updateFn  func(context.Context, *db.User) error
	deleteFn  func(context.Context, int64) error
}

func (m *userRepoMock) Create(ctx context.Context, u *db.User) error { return m.createFn(ctx, u) }

func (m *userRepoMock) GetByID(ctx context.Context, id int64) (db.User, error) {
	return m.getByIDFn(ctx, id)
}

func (m *userRepoMock) Update(ctx context.Context, u *db.User) error { return m.updateFn(ctx, u) }

func (m *userRepoMock) Delete(ctx context.Context, id int64) error { return m.deleteFn(ctx, id) }

type roleProviderMock struct {
	getByIDFn func(context.Context, dict.EntityID) (dict.Entity, error)
}

func (m *roleProviderMock) GetByID(ctx context.Context, id dict.EntityID) (dict.Entity, error) {
	return m.getByIDFn(ctx, id)
}

func (m *roleProviderMock) List(ctx context.Context) ([]dict.Entity, error) {
	return nil, nil
}
