package incident

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cyradin/fixik/internal/db"
	"github.com/cyradin/fixik/internal/user"
	"github.com/stretchr/testify/require"
)

func TestCommentManager_Create(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		ctx  context.Context
		mock func(*commentRepoMock)
		err  bool
	}{
		{
			name: "unauthorized",
			ctx:  ctx,
			mock: func(repo *commentRepoMock) {},
			err:  true,
		},
		{
			name: "repo error",
			ctx:  user.WithContext(ctx, user.User{ID: 1}),
			mock: func(repo *commentRepoMock) {
				repo.createFn = func(ctx context.Context, c *db.Comment) error {
					return errors.New("db error")
				}
			},
			err: true,
		},
		{
			name: "success",
			ctx:  user.WithContext(ctx, user.User{ID: 1}),
			mock: func(repo *commentRepoMock) {
				repo.createFn = func(ctx context.Context, c *db.Comment) error {
					c.ID = 10
					c.CreatedAt = time.Now()
					c.UpdatedAt = time.Now()
					return nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &commentRepoMock{}
			userProvider := &userProviderMock{}

			tt.mock(repo)

			manager := NewCommentManager(repo, userProvider)

			res, err := manager.Create(tt.ctx, 100, "hello")

			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, int64(10), res.ID)
			require.Equal(t, int64(100), res.IncidentID)
			require.Equal(t, int64(1), res.Author.ID)
			require.Equal(t, "hello", res.Text)
		})
	}
}

func TestCommentManager_ListByIncident(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	dbComments := []db.Comment{
		{
			ID:         1,
			IncidentID: 100,
			AuthorID:   10,
			Text:       "one",
		},
		{
			ID:         2,
			IncidentID: 100,
			AuthorID:   11,
			Text:       "two",
		},
	}

	tests := []struct {
		name string
		mock func(*commentRepoMock, *userProviderMock)
		err  bool
	}{
		{
			name: "repo error",
			mock: func(repo *commentRepoMock, up *userProviderMock) {
				repo.listFn = func(context.Context, int64, int, int) ([]db.Comment, error) {
					return nil, errors.New("repo error")
				}
			},
			err: true,
		},
		{
			name: "user provider error",
			mock: func(repo *commentRepoMock, up *userProviderMock) {
				repo.listFn = func(context.Context, int64, int, int) ([]db.Comment, error) {
					return dbComments, nil
				}

				up.getByIDManyFn = func(context.Context, []int64) ([]user.User, error) {
					return nil, errors.New("user error")
				}
			},
			err: true,
		},
		{
			name: "user not found",
			mock: func(repo *commentRepoMock, up *userProviderMock) {
				repo.listFn = func(context.Context, int64, int, int) ([]db.Comment, error) {
					return dbComments, nil
				}

				up.getByIDManyFn = func(context.Context, []int64) ([]user.User, error) {
					return []user.User{
						{ID: 10},
						// 11 missing
					}, nil
				}
			},
			err: true,
		},
		{
			name: "success",
			mock: func(repo *commentRepoMock, up *userProviderMock) {
				repo.listFn = func(context.Context, int64, int, int) ([]db.Comment, error) {
					return dbComments, nil
				}

				up.getByIDManyFn = func(ctx context.Context, ids []int64) ([]user.User, error) {
					users := make([]user.User, 0, len(ids))
					for _, id := range ids {
						users = append(users, user.User{ID: id})
					}
					return users, nil
				}
			},
		},
		{
			name: "empty list",
			mock: func(repo *commentRepoMock, up *userProviderMock) {
				repo.listFn = func(context.Context, int64, int, int) ([]db.Comment, error) {
					return []db.Comment{}, nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &commentRepoMock{}
			userProvider := &userProviderMock{}

			tt.mock(repo, userProvider)

			manager := NewCommentManager(repo, userProvider)

			res, err := manager.ListByIncident(ctx, 100, 100, 100)

			if tt.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			if tt.name == "empty list" {
				require.Empty(t, res)
				return
			}

			require.Len(t, res, 2)
			require.Equal(t, int64(1), res[0].ID)
			require.Equal(t, int64(10), res[0].Author.ID)
			require.Equal(t, "one", res[0].Text)

			require.Equal(t, int64(2), res[1].ID)
			require.Equal(t, int64(11), res[1].Author.ID)
			require.Equal(t, "two", res[1].Text)
		})
	}
}

type commentRepoMock struct {
	createFn func(context.Context, *db.Comment) error
	updateFn func(context.Context, *db.Comment) error
	deleteFn func(context.Context, int64) error
	listFn   func(context.Context, int64, int, int) ([]db.Comment, error)
}

func (m *commentRepoMock) Create(ctx context.Context, c *db.Comment) error {
	return m.createFn(ctx, c)
}

func (m *commentRepoMock) Update(ctx context.Context, c *db.Comment) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, c)
	}
	return nil
}

func (m *commentRepoMock) Delete(ctx context.Context, id int64) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *commentRepoMock) ListByIncident(ctx context.Context, incidentID int64, limit, offset int) ([]db.Comment, error) {
	return m.listFn(ctx, incidentID, limit, offset)
}
