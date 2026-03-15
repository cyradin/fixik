package user

import (
	"context"
	"errors"
	"testing"

	"github.com/cyradin/fixik/internal/db"
	"github.com/stretchr/testify/require"
)

func TestUserManager_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		user CreateUser
		mock func(*userRepoMock)
		err  bool
	}{
		{
			name: "success",
			user: CreateUser{
				Name:     "Алиса",
				Username: "alice",
				Email:    "alice@example.com",
				Password: "pass",
				TeamID:   new(int64(1)),
				Role:     RoleUser,
			},
			mock: func(m *userRepoMock) {
				m.createFn = func(ctx context.Context, u *db.User) error {
					u.ID = 10
					return nil
				}
				m.getByIDFn = func(ctx context.Context, id int64) (db.User, error) {
					return db.User{
						ID:       id,
						Name:     "Алиса",
						Username: "alice",
						Email:    "alice@example.com",
						Password: "hashed",
						TeamID:   new(int64(1)),
						Role:     RoleUser,
					}, nil
				}
			},
		},
		{
			name: "repo error",
			user: CreateUser{
				Name:     "Боб",
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

			manager := NewUserManager(repo)

			u, err := manager.Create(t.Context(), tt.user)
			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.user.Name, u.Name)
			require.Equal(t, tt.user.Username, u.Username)
			require.Equal(t, tt.user.Email, u.Email)
			require.Equal(t, tt.user.TeamID, u.TeamID)
			require.Equal(t, tt.user.Role, u.Role)
		})
	}
}

func TestUserManager_GetByID(t *testing.T) {
	t.Parallel()

	dbUser := db.User{
		ID:       1,
		Name:     "Алиса",
		Username: "alice",
		Email:    "alice@example.com",
		Password: "hashed",
		TeamID:   new(int64(1)),
		Role:     RoleAdmin,
	}

	tests := []struct {
		name string
		mock func(*userRepoMock)
		err  bool
	}{
		{
			name: "repo error",
			mock: func(repo *userRepoMock) {
				repo.getByIDFn = func(ctx context.Context, id int64) (db.User, error) {
					return db.User{}, errors.New("repo error")
				}
			},
			err: true,
		},
		{
			name: "success",
			mock: func(repo *userRepoMock) {
				repo.getByIDFn = func(ctx context.Context, id int64) (db.User, error) {
					return dbUser, nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &userRepoMock{}
			tt.mock(repo)

			manager := NewUserManager(repo)
			u, err := manager.GetByID(t.Context(), 1)

			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, User{
				ID:       1,
				Name:     "Алиса",
				Username: "alice",
				Email:    "alice@example.com",
				TeamID:   new(int64(1)),
				Role:     RoleAdmin,
			}, u)
		})
	}
}

func TestUserManager_Update(t *testing.T) {
	t.Parallel()

	newName := "Новый Боб"

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
				Name:     &newName,
				Username: new(string),
				Email:    new(string),
				Password: new(string),
				TeamID:   new(int64),
				Role:     new(RoleManager),
			},
			mock: func(m *userRepoMock) {
				m.getByIDFn = func(ctx context.Context, id int64) (db.User, error) {
					return db.User{ID: id, Name: "Боб", Role: RoleManager}, nil
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

			manager := NewUserManager(repo)

			_, err := manager.Update(t.Context(), tt.cmd)
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

			manager := NewUserManager(repo)
			err := manager.Delete(t.Context(), tt.id)

			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestUserManager_List(t *testing.T) {
	t.Parallel()

	dbUsers := []db.User{
		{
			ID:       1,
			Name:     "Алиса",
			Username: "alice",
			Role:     RoleAdmin,
		},
		{
			ID:       2,
			Name:     "Боб",
			Username: "bob",
			Role:     RoleUser,
		},
	}

	tests := []struct {
		name string
		mock func(*userRepoMock)
		err  bool
	}{
		{
			name: "repo error",
			mock: func(repo *userRepoMock) {
				repo.listFn = func(ctx context.Context, limit, offset int) ([]db.User, error) {
					return nil, errors.New("repo error")
				}
			},
			err: true,
		},
		{
			name: "success",
			mock: func(repo *userRepoMock) {
				repo.listFn = func(ctx context.Context, limit, offset int) ([]db.User, error) {
					return dbUsers, nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &userRepoMock{}

			tt.mock(repo)

			manager := NewUserManager(repo)

			users, err := manager.List(t.Context(), 10, 0)

			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Len(t, users, 2)
			require.Equal(t, "Алиса", users[0].Name)
			require.Equal(t, "Боб", users[1].Name)
			require.Equal(t, users[0].ID, int64(1))
			require.Equal(t, users[1].ID, int64(2))
		})
	}
}

// Mocks
type userRepoMock struct {
	createFn  func(context.Context, *db.User) error
	getByIDFn func(context.Context, int64) (db.User, error)
	listFn    func(context.Context, int, int) ([]db.User, error)
	updateFn  func(context.Context, *db.User) error
	deleteFn  func(context.Context, int64) error
}

func (m *userRepoMock) Create(ctx context.Context, u *db.User) error {
	return m.createFn(ctx, u)
}

func (m *userRepoMock) GetByID(ctx context.Context, id int64) (db.User, error) {
	return m.getByIDFn(ctx, id)
}

func (m *userRepoMock) List(ctx context.Context, limit, offset int) ([]db.User, error) {
	return m.listFn(ctx, limit, offset)
}

func (m *userRepoMock) Update(ctx context.Context, u *db.User) error {
	return m.updateFn(ctx, u)
}

func (m *userRepoMock) Delete(ctx context.Context, id int64) error {
	return m.deleteFn(ctx, id)
}
