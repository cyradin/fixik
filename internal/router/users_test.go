package router

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cyradin/fixik/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockUserManager{
			createFn: func(ctx context.Context, u user.CreateUser) (user.User, error) {
				return user.User{ID: 1, Username: u.Username, Email: u.Email, TeamID: u.TeamID}, nil
			},
		}
		req := CreateUserRequest{Username: "abc", Email: "a@b.com", Password: "123456", TeamID: 1}
		rr := testRequest(t, createUser(m), http.MethodPost, "/users", req)
		require.Equal(t, http.StatusOK, rr.Code)

		var resp UserResponse

		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		require.Equal(t, int64(1), resp.ID)
	})

	t.Run("validation error", func(t *testing.T) {
		t.Parallel()

		m := &mockUserManager{}
		req := CreateUserRequest{}
		rr := testRequest(t, createUser(m), http.MethodPost, "/users", req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockUserManager{
			createFn: func(ctx context.Context, u user.CreateUser) (user.User, error) {
				return user.User{}, errors.New("fail")
			},
		}
		req := CreateUserRequest{Username: "abc", Email: "a@b.com", Password: "123456", TeamID: 1}
		rr := testRequest(t, createUser(m), http.MethodPost, "/users", req)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestGetUser(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockUserManager{
			getFn: func(ctx context.Context, id int64) (user.User, error) {
				return user.User{ID: id, Username: "abc", Email: "a@b.com"}, nil
			},
		}
		r := chi.NewRouter()
		r.Get("/users/{id}", getUser(m))

		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Code)

		var resp UserResponse

		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		require.Equal(t, int64(1), resp.ID)
	})

	t.Run("invalid id", func(t *testing.T) {
		t.Parallel()

		m := &mockUserManager{}
		r := chi.NewRouter()
		r.Get("/users/{id}", getUser(m))

		req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockUserManager{
			getFn: func(ctx context.Context, id int64) (user.User, error) {
				return user.User{}, errors.New("not found")
			},
		}
		r := chi.NewRouter()
		r.Get("/users/{id}", getUser(m))

		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockUserManager{
			updateFn: func(ctx context.Context, u user.UpdateUser) (user.User, error) {
				return user.User{ID: u.ID, Username: *u.Username}, nil
			},
		}
		r := chi.NewRouter()
		r.Patch("/users/{id}", updateUser(m))

		username := "updated"
		reqBody := UpdateUserRequest{Username: &username}
		req := httptest.NewRequest(http.MethodPatch, "/users/1", jsonBody(t, reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		t.Parallel()

		m := &mockUserManager{}
		r := chi.NewRouter()
		r.Patch("/users/{id}", updateUser(m))

		username := "updated"
		reqBody := UpdateUserRequest{Username: &username}
		req := httptest.NewRequest(http.MethodPatch, "/users/abc", jsonBody(t, reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockUserManager{
			updateFn: func(ctx context.Context, u user.UpdateUser) (user.User, error) {
				return user.User{}, errors.New("fail")
			},
		}
		r := chi.NewRouter()
		r.Patch("/users/{id}", updateUser(m))

		username := "updated"
		reqBody := UpdateUserRequest{Username: &username}
		req := httptest.NewRequest(http.MethodPatch, "/users/1", jsonBody(t, reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		deleted := false
		m := &mockUserManager{
			deleteFn: func(ctx context.Context, id int64) error {
				deleted = true
				return nil
			},
		}
		r := chi.NewRouter()
		r.Delete("/users/{id}", deleteUser(m))

		req := httptest.NewRequest(http.MethodDelete, "/users/42", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Code)
		require.True(t, deleted)
	})

	t.Run("invalid id", func(t *testing.T) {
		t.Parallel()

		m := &mockUserManager{}
		r := chi.NewRouter()
		r.Delete("/users/{id}", deleteUser(m))

		req := httptest.NewRequest(http.MethodDelete, "/users/abc", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockUserManager{
			deleteFn: func(ctx context.Context, id int64) error { return errors.New("fail") },
		}
		r := chi.NewRouter()
		r.Delete("/users/{id}", deleteUser(m))

		req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestListUsers(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockUserManager{
			listFn: func(ctx context.Context, limit, offset int) ([]user.User, error) {
				return []user.User{
					{ID: 1, Username: "A"},
					{ID: 2, Username: "B"},
				}, nil
			},
		}
		rr := testRequest(t, listUsers(m), http.MethodGet, "/users?limit=10&offset=0", nil)
		require.Equal(t, http.StatusOK, rr.Code)

		var resp ListUsersResponse

		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		require.Len(t, resp.Items, 2)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockUserManager{
			listFn: func(ctx context.Context, limit, offset int) ([]user.User, error) { return nil, errors.New("fail") },
		}
		rr := testRequest(t, listUsers(m), http.MethodGet, "/users?limit=10&offset=0", nil)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

type mockUserManager struct {
	createFn func(ctx context.Context, u user.CreateUser) (user.User, error)
	getFn    func(ctx context.Context, id int64) (user.User, error)
	updateFn func(ctx context.Context, u user.UpdateUser) (user.User, error)
	deleteFn func(ctx context.Context, id int64) error
	listFn   func(ctx context.Context, limit, offset int) ([]user.User, error)
}

func (m *mockUserManager) Create(ctx context.Context, u user.CreateUser) (user.User, error) {
	return m.createFn(ctx, u)
}

func (m *mockUserManager) GetByID(ctx context.Context, id int64) (user.User, error) {
	return m.getFn(ctx, id)
}

func (m *mockUserManager) Update(ctx context.Context, u user.UpdateUser) (user.User, error) {
	return m.updateFn(ctx, u)
}

func (m *mockUserManager) Delete(ctx context.Context, id int64) error {
	return m.deleteFn(ctx, id)
}

func (m *mockUserManager) List(ctx context.Context, limit, offset int) ([]user.User, error) {
	return m.listFn(ctx, limit, offset)
}
