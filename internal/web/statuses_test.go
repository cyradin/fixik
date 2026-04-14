package web

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cyradin/fixik/internal/role"
	"github.com/cyradin/fixik/internal/status"
	"github.com/cyradin/fixik/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestCreateStatus(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockStatusManager{
			createFn: func(ctx context.Context, s status.Status) (status.Status, error) {
				s.ID = 1
				return s, nil
			},
		}

		req := CreateStatusRequest{
			Name: "Test", Code: "TST", Description: "desc", Sort: 1, IsFinal: true,
		}

		rr := testRequest(t, createStatus(m), http.MethodPost, "/dummy", req)
		require.Equal(t, http.StatusOK, rr.Code)

		var resp Status
		require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
		require.Equal(t, int64(1), resp.ID)
		require.True(t, resp.IsFinal)
	})

	t.Run("validation error", func(t *testing.T) {
		t.Parallel()

		m := &mockStatusManager{}
		req := CreateStatusRequest{Name: "", Code: ""}

		rr := testRequest(t, createStatus(m), http.MethodPost, "/dummy", req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockStatusManager{
			createFn: func(ctx context.Context, s status.Status) (status.Status, error) {
				return status.Status{}, errors.New("fail")
			},
		}

		req := CreateStatusRequest{Name: "A", Code: "A", Sort: 1}

		rr := testRequest(t, createStatus(m), http.MethodPost, "/dummy", req)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestGetStatus(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockStatusManager{
			getFn: func(ctx context.Context, id int64) (status.Status, error) {
				return status.Status{ID: id, Name: "A"}, nil
			},
		}

		r := chi.NewRouter()
		r.Get("/statuses/{id}", getStatus(m))

		req := httptest.NewRequest(http.MethodGet, "/statuses/1", nil).
			WithContext(user.WithContext(t.Context(), user.User{
				Role: role.Admin,
			}))
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("403", func(t *testing.T) {
		t.Parallel()

		m := &mockStatusManager{
			getFn: func(ctx context.Context, id int64) (status.Status, error) {
				return status.Status{ID: id, Name: "A"}, nil
			},
		}

		r := chi.NewRouter()
		r.Get("/statuses/{id}", getStatus(m))

		req := httptest.NewRequest(http.MethodGet, "/statuses/1", nil)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusForbidden, rr.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		t.Parallel()

		m := &mockStatusManager{}
		r := chi.NewRouter()
		r.Get("/statuses/{id}", getStatus(m))

		req := httptest.NewRequest(http.MethodGet, "/statuses/abc", nil).
			WithContext(user.WithContext(t.Context(), user.User{
				Role: role.Admin,
			}))
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockStatusManager{
			getFn: func(ctx context.Context, id int64) (status.Status, error) {
				return status.Status{}, errors.New("fail")
			},
		}

		r := chi.NewRouter()
		r.Get("/statuses/{id}", getStatus(m))

		req := httptest.NewRequest(http.MethodGet, "/statuses/1", nil).
			WithContext(user.WithContext(t.Context(), user.User{
				Role: role.Admin,
			}))
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestUpdateStatus(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockStatusManager{
			updateFn: func(ctx context.Context, s status.Status) (status.Status, error) {
				return s, nil
			},
		}

		r := chi.NewRouter()
		r.Put("/statuses/{id}", updateStatus(m))

		reqBody := UpdateStatusRequest{Name: "A", Code: "A", Sort: 1, IsFinal: true}
		req := httptest.NewRequest(http.MethodPut, "/statuses/1", jsonBody(t, reqBody)).
			WithContext(user.WithContext(t.Context(), user.User{
				Role: role.Admin,
			}))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		t.Parallel()

		m := &mockStatusManager{}
		r := chi.NewRouter()
		r.Put("/statuses/{id}", updateStatus(m))

		req := httptest.NewRequest(http.MethodPut, "/statuses/1", jsonBody(t, UpdateStatusRequest{})).
			WithContext(user.WithContext(t.Context(), user.User{
				Role: role.Admin,
			}))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		t.Parallel()

		m := &mockStatusManager{}
		r := chi.NewRouter()
		r.Put("/statuses/{id}", updateStatus(m))

		req := httptest.NewRequest(http.MethodPut, "/statuses/abc", jsonBody(t, UpdateStatusRequest{})).
			WithContext(user.WithContext(t.Context(), user.User{
				Role: role.Admin,
			}))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockStatusManager{
			updateFn: func(ctx context.Context, s status.Status) (status.Status, error) {
				return status.Status{}, errors.New("fail")
			},
		}

		r := chi.NewRouter()
		r.Put("/statuses/{id}", updateStatus(m))

		req := httptest.NewRequest(http.MethodPut, "/statuses/1", jsonBody(t, UpdateStatusRequest{
			Name: "A", Code: "A", Sort: 1,
		})).
			WithContext(user.WithContext(t.Context(), user.User{
				Role: role.Admin,
			}))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestDeleteStatus(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		deleted := false
		m := &mockStatusManager{
			deleteFn: func(ctx context.Context, id int64) error {
				deleted = true
				return nil
			},
		}

		r := chi.NewRouter()
		r.Delete("/statuses/{id}", deleteStatus(m))

		req := httptest.NewRequest(http.MethodDelete, "/statuses/1", nil).
			WithContext(user.WithContext(t.Context(), user.User{
				Role: role.Admin,
			}))
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		require.True(t, deleted)
	})

	t.Run("invalid id", func(t *testing.T) {
		t.Parallel()

		m := &mockStatusManager{}
		r := chi.NewRouter()
		r.Delete("/statuses/{id}", deleteStatus(m))

		req := httptest.NewRequest(http.MethodDelete, "/statuses/abc", nil).
			WithContext(user.WithContext(t.Context(), user.User{
				Role: role.Admin,
			}))
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockStatusManager{
			deleteFn: func(ctx context.Context, id int64) error {
				return errors.New("fail")
			},
		}

		r := chi.NewRouter()
		r.Delete("/statuses/{id}", deleteStatus(m))

		req := httptest.NewRequest(http.MethodDelete, "/statuses/1", nil).
			WithContext(user.WithContext(t.Context(), user.User{
				Role: role.Admin,
			}))
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestListStatuses(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockStatusManager{
			listFn: func(ctx context.Context) ([]status.Status, error) {
				return []status.Status{
					{ID: 1, Name: "A"},
					{ID: 2, Name: "B"},
				}, nil
			},
		}

		rr := testRequest(t, listStatuses(m), http.MethodGet, "/dummy", nil)
		require.Equal(t, http.StatusOK, rr.Code)

		var resp ListStatusesResponse
		require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
		require.Len(t, resp.Items, 2)
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()

		m := &mockStatusManager{
			listFn: func(ctx context.Context) ([]status.Status, error) {
				return []status.Status{}, nil
			},
		}

		rr := testRequest(t, listStatuses(m), http.MethodGet, "/dummy", nil)
		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockStatusManager{
			listFn: func(ctx context.Context) ([]status.Status, error) {
				return nil, errors.New("fail")
			},
		}

		rr := testRequest(t, listStatuses(m), http.MethodGet, "/dummy", nil)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

type mockStatusManager struct {
	createFn func(ctx context.Context, s status.Status) (status.Status, error)
	getFn    func(ctx context.Context, id int64) (status.Status, error)
	listFn   func(ctx context.Context) ([]status.Status, error)
	updateFn func(ctx context.Context, s status.Status) (status.Status, error)
	deleteFn func(ctx context.Context, id int64) error
}

func (m *mockStatusManager) Create(ctx context.Context, s status.Status) (status.Status, error) {
	if m.createFn != nil {
		return m.createFn(ctx, s)
	}

	return status.Status{}, nil
}

func (m *mockStatusManager) GetByID(ctx context.Context, id int64) (status.Status, error) {
	if m.getFn != nil {
		return m.getFn(ctx, id)
	}

	return status.Status{}, nil
}

func (m *mockStatusManager) List(ctx context.Context) ([]status.Status, error) {
	if m.listFn != nil {
		return m.listFn(ctx)
	}

	return nil, nil
}

func (m *mockStatusManager) Update(ctx context.Context, s status.Status) (status.Status, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, s)
	}

	return status.Status{}, nil
}

func (m *mockStatusManager) Delete(ctx context.Context, id int64) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}

	return nil
}
