package user

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/cyradin/fixik/internal/db"
	"github.com/cyradin/fixik/internal/role"
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
				Role:     role.User,
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
						Role:     role.User,
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

			manager := NewUserManager(repo, &incidentsCounterMock{}, &txExecutorMock{})

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

func TestUserManager_GetByUsername(t *testing.T) {
	t.Parallel()

	dbUser := db.User{
		ID:       1,
		Name:     "Алиса",
		Username: "alice",
		Email:    "alice@example.com",
		Password: "hashed",
		TeamID:   new(int64(1)),
		Role:     role.Admin,
	}

	tests := []struct {
		name string
		mock func(*userRepoMock)
		err  bool
	}{
		{
			name: "repo error",
			mock: func(repo *userRepoMock) {
				repo.getByUsernameFn = func(ctx context.Context, username string) (db.User, error) {
					return db.User{}, errors.New("repo error")
				}
			},
			err: true,
		},
		{
			name: "success",
			mock: func(repo *userRepoMock) {
				repo.getByUsernameFn = func(ctx context.Context, username string) (db.User, error) {
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

			manager := NewUserManager(repo, &incidentsCounterMock{}, &txExecutorMock{})
			u, err := manager.GetByUsername(t.Context(), "alice")

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
				Role:     role.Admin,
			}, u)
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
		Role:     role.Admin,
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

			manager := NewUserManager(repo, &incidentsCounterMock{}, &txExecutorMock{})
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
				Role:     role.Admin,
			}, u)
		})
	}
}

func TestUserManager_GetByIDMany(t *testing.T) {
	t.Parallel()

	dbUsers := []db.User{
		{ID: 1, Name: "Алиса", Username: "alice", Role: role.Admin, TeamID: new(int64(1))},
		{ID: 2, Name: "Боб", Username: "bob", Role: role.User, TeamID: new(int64(2))},
		{ID: 3, Name: "Каролина", Username: "carol", Role: role.Manager, TeamID: new(int64(3))},
	}

	tests := []struct {
		name string
		ids  []int64
		mock func(*userRepoMock)
		err  bool
		want []User
	}{
		{
			name: "success all",
			ids:  []int64{1, 2, 3},
			mock: func(m *userRepoMock) {
				m.getByIDManyFn = func(ctx context.Context, ids []int64) ([]db.User, error) {
					return dbUsers, nil
				}
			},
			want: []User{
				{ID: 1, Name: "Алиса", Username: "alice", Role: role.Admin, TeamID: new(int64(1))},
				{ID: 2, Name: "Боб", Username: "bob", Role: role.User, TeamID: new(int64(2))},
				{ID: 3, Name: "Каролина", Username: "carol", Role: role.Manager, TeamID: new(int64(3))},
			},
		},
		{
			name: "repo error",
			ids:  []int64{1, 2},
			mock: func(m *userRepoMock) {
				m.getByIDManyFn = func(ctx context.Context, ids []int64) ([]db.User, error) {
					return nil, errors.New("db error")
				}
			},
			err: true,
		},
		{
			name: "empty input",
			ids:  []int64{},
			mock: func(m *userRepoMock) {
				m.getByIDManyFn = func(ctx context.Context, ids []int64) ([]db.User, error) {
					return []db.User{}, nil
				}
			},
			want: []User{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &userRepoMock{}
			tt.mock(repo)

			manager := NewUserManager(repo, &incidentsCounterMock{}, &txExecutorMock{})
			users, err := manager.GetByIDMany(t.Context(), tt.ids)

			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			for i := range users {
				require.Equal(t, tt.want[i].ID, users[i].ID)
				require.Equal(t, tt.want[i].Name, users[i].Name)
				require.Equal(t, tt.want[i].Username, users[i].Username)
				require.Equal(t, tt.want[i].Role, users[i].Role)
				require.Equal(t, tt.want[i].TeamID, users[i].TeamID)
			}
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
				Role:     new(role.Manager),
			},
			mock: func(m *userRepoMock) {
				m.getByIDFn = func(ctx context.Context, id int64) (db.User, error) {
					return db.User{ID: id, Name: "Боб", Role: role.Manager}, nil
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

			manager := NewUserManager(repo, &incidentsCounterMock{}, &txExecutorMock{})

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

	ctx := context.Background()

	tests := []struct {
		name string
		id   int64
		mock func(*userRepoMock, *incidentsCounterMock, *txExecutorMock)
		err  error
	}{
		{
			name: "success",
			id:   1,
			mock: func(r *userRepoMock, c *incidentsCounterMock, tx *txExecutorMock) {
				tx.execFn = func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				}

				c.countFn = func(ctx context.Context, id int64) (int, error) {
					require.Equal(t, int64(1), id)
					return 0, nil
				}

				r.deleteFn = func(ctx context.Context, id int64) error {
					require.Equal(t, int64(1), id)
					return nil
				}
			},
		},
		{
			name: "has dependencies",
			id:   1,
			mock: func(r *userRepoMock, c *incidentsCounterMock, tx *txExecutorMock) {
				tx.execFn = func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				}

				c.countFn = func(ctx context.Context, id int64) (int, error) {
					return 2, nil
				}
			},
			err: ErrHasDependantEntities,
		},
		{
			name: "count error",
			id:   1,
			mock: func(r *userRepoMock, c *incidentsCounterMock, tx *txExecutorMock) {
				tx.execFn = func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				}

				c.countFn = func(ctx context.Context, id int64) (int, error) {
					return 0, errors.New("db error")
				}
			},
			err: fmt.Errorf("count incidents"),
		},
		{
			name: "repo error",
			id:   1,
			mock: func(r *userRepoMock, c *incidentsCounterMock, tx *txExecutorMock) {
				tx.execFn = func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				}

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

			repo := &userRepoMock{}
			counter := &incidentsCounterMock{}
			tx := &txExecutorMock{}

			tt.mock(repo, counter, tx)

			manager := NewUserManager(repo, counter, tx)

			err := manager.Delete(ctx, tt.id)

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

func TestUserManager_List(t *testing.T) {
	t.Parallel()

	dbUsers := []db.User{
		{
			ID:       1,
			Name:     "Алиса",
			Username: "alice",
			Role:     role.Admin,
		},
		{
			ID:       2,
			Name:     "Боб",
			Username: "bob",
			Role:     role.User,
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

			manager := NewUserManager(repo, &incidentsCounterMock{}, &txExecutorMock{})

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
	createFn        func(context.Context, *db.User) error
	getByIDFn       func(context.Context, int64) (db.User, error)
	getByUsernameFn func(context.Context, string) (db.User, error)
	getByIDManyFn   func(context.Context, []int64) ([]db.User, error)
	listFn          func(context.Context, int, int) ([]db.User, error)
	updateFn        func(context.Context, *db.User) error
	deleteFn        func(context.Context, int64) error
}

func (m *userRepoMock) Create(ctx context.Context, u *db.User) error {
	return m.createFn(ctx, u)
}

func (m *userRepoMock) GetByID(ctx context.Context, id int64) (db.User, error) {
	return m.getByIDFn(ctx, id)
}

func (m *userRepoMock) GetByUsername(ctx context.Context, username string) (db.User, error) {
	return m.getByUsernameFn(ctx, username)
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

func (m *userRepoMock) GetByIDMany(ctx context.Context, ids []int64) ([]db.User, error) {
	return m.getByIDManyFn(ctx, ids)
}

type incidentsCounterMock struct {
	countFn func(ctx context.Context, id int64) (int, error)
}

func (m *incidentsCounterMock) CountByUser(ctx context.Context, id int64) (int, error) {
	return m.countFn(ctx, id)
}

type txExecutorMock struct {
	execFn func(ctx context.Context, fn func(ctx context.Context) error) error
}

func (m *txExecutorMock) Exec(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.execFn(ctx, fn)
}
