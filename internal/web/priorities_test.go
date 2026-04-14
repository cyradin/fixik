package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cyradin/fixik/internal/priority"
	"github.com/cyradin/fixik/internal/role"
	"github.com/cyradin/fixik/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestCreatePriority(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockPriorityManager{
			createFn: func(ctx context.Context, e priority.Priority) (priority.Priority, error) {
				return priority.Priority{
					ID:          1,
					Name:        e.Name,
					Code:        e.Code,
					Description: e.Description,
					Sort:        e.Sort,
				}, nil
			},
		}
		req := CreatePriorityRequest{
			Name:        "Test",
			Code:        "TST",
			Description: "desc",
			Sort:        42,
		}
		rr := testRequest(t, createPriority(m), http.MethodPost, "/dummy", req)
		require.Equal(t, http.StatusOK, rr.Code)

		var resp Priority

		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		require.Equal(t, int64(1), resp.ID)
		require.Equal(t, "desc", resp.Description)
		require.Equal(t, 42, resp.Sort)
	})

	t.Run("validation error", func(t *testing.T) {
		t.Parallel()

		m := &mockPriorityManager{}
		req := CreatePriorityRequest{Name: "", Code: ""}
		rr := testRequest(t, createPriority(m), http.MethodPost, "/dummy", req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockPriorityManager{
			createFn: func(ctx context.Context, e priority.Priority) (priority.Priority, error) {
				return priority.Priority{}, errors.New("create failed")
			},
		}
		req := CreatePriorityRequest{Name: "Name", Code: "C", Sort: 123}
		rr := testRequest(t, createPriority(m), http.MethodPost, "/dummy", req)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestGetPriority(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockPriorityManager{
			getFn: func(ctx context.Context, id int64) (priority.Priority, error) {
				return priority.Priority{ID: id, Name: "A"}, nil
			},
		}
		r := chi.NewRouter()
		r.Get("/entities/{id}", getPriority(m))

		req := httptest.NewRequest(http.MethodGet, "/entities/1", nil).
			WithContext(user.WithContext(t.Context(), user.User{
				Role: role.Admin,
			}))
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Code)

		var resp Priority

		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		require.Equal(t, int64(1), resp.ID)
	})

	t.Run("invalid id", func(t *testing.T) {
		t.Parallel()

		m := &mockPriorityManager{}
		r := chi.NewRouter()

		r.Get("/entities/{id}", getPriority(m))

		req := httptest.NewRequest(http.MethodGet, "/entities/abc", nil).
			WithContext(user.WithContext(t.Context(), user.User{
				Role: role.Admin,
			}))
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockPriorityManager{
			getFn: func(ctx context.Context, id int64) (priority.Priority, error) {
				return priority.Priority{}, errors.New("not found")
			},
		}
		r := chi.NewRouter()
		r.Get("/entities/{id}", getPriority(m))

		req := httptest.NewRequest(http.MethodGet, "/entities/1", nil).
			WithContext(user.WithContext(t.Context(), user.User{
				Role: role.Admin,
			}))
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestUpdatePriority(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockPriorityManager{
			updateFn: func(ctx context.Context, e priority.Priority) (priority.Priority, error) { return e, nil },
		}
		r := chi.NewRouter()
		r.Put("/entities/{id}", updatePriority(m))

		reqBody := UpdatePriorityRequest{Name: "Updated", Code: "UPD", Sort: 123}
		req := httptest.NewRequest(http.MethodPut, "/entities/1", jsonBody(t, reqBody)).
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

		m := &mockPriorityManager{}
		r := chi.NewRouter()
		r.Put("/entities/{id}", updatePriority(m))

		reqBody := UpdatePriorityRequest{Name: "", Code: ""}
		req := httptest.NewRequest(http.MethodPut, "/entities/1", jsonBody(t, reqBody)).
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

		m := &mockPriorityManager{}
		r := chi.NewRouter()
		r.Put("/entities/{id}", updatePriority(m))

		reqBody := UpdatePriorityRequest{Name: "X", Code: "X"}
		req := httptest.NewRequest(http.MethodPut, "/entities/abc", jsonBody(t, reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockPriorityManager{
			updateFn: func(ctx context.Context, e priority.Priority) (priority.Priority, error) {
				return priority.Priority{}, errors.New("fail")
			},
		}
		r := chi.NewRouter()
		r.Put("/entities/{id}", updatePriority(m))

		reqBody := UpdatePriorityRequest{Name: "X", Code: "X", Sort: 123}
		req := httptest.NewRequest(http.MethodPut, "/entities/1", jsonBody(t, reqBody)).
			WithContext(user.WithContext(t.Context(), user.User{
				Role: role.Admin,
			}))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestDeletePriority(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		deleted := false
		m := &mockPriorityManager{deleteFn: func(ctx context.Context, id int64) error { deleted = true; return nil }}
		r := chi.NewRouter()
		r.Delete("/entities/{id}", deletePriority(m))

		req := httptest.NewRequest(http.MethodDelete, "/entities/42", nil).
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

		m := &mockPriorityManager{}
		r := chi.NewRouter()
		r.Delete("/entities/{id}", deletePriority(m))

		req := httptest.NewRequest(http.MethodDelete, "/entities/abc", nil).
			WithContext(user.WithContext(t.Context(), user.User{
				Role: role.Admin,
			}))
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockPriorityManager{deleteFn: func(ctx context.Context, id int64) error { return errors.New("fail") }}
		r := chi.NewRouter()
		r.Delete("/entities/{id}", deletePriority(m))

		req := httptest.NewRequest(http.MethodDelete, "/entities/1", nil).
			WithContext(user.WithContext(t.Context(), user.User{
				Role: role.Admin,
			}))
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestListPriorities(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockPriorityManager{
			listFn: func(ctx context.Context) ([]priority.Priority, error) {
				return []priority.Priority{
					{ID: 1, Name: "A", Code: "A", Sort: 1},
					{ID: 2, Name: "B", Code: "B", Sort: 2},
				}, nil
			},
		}
		rr := testRequest(t, listPriorities(m), http.MethodGet, "/dummy", nil)
		require.Equal(t, http.StatusOK, rr.Code)

		var resp ListPrioritiesResponse

		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		require.Len(t, resp.Items, 2)
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()

		m := &mockPriorityManager{
			listFn: func(ctx context.Context) ([]priority.Priority, error) { return []priority.Priority{}, nil },
		}
		rr := testRequest(t, listPriorities(m), http.MethodGet, "/dummy", nil)
		require.Equal(t, http.StatusOK, rr.Code)

		var resp ListPrioritiesResponse

		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		require.Len(t, resp.Items, 0)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockPriorityManager{
			listFn: func(ctx context.Context) ([]priority.Priority, error) { return nil, errors.New("fail") },
		}
		rr := testRequest(t, listPriorities(m), http.MethodGet, "/dummy", nil)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

type mockPriorityManager struct {
	createFn func(ctx context.Context, e priority.Priority) (priority.Priority, error)
	getFn    func(ctx context.Context, id int64) (priority.Priority, error)
	listFn   func(ctx context.Context) ([]priority.Priority, error)
	updateFn func(ctx context.Context, e priority.Priority) (priority.Priority, error)
	deleteFn func(ctx context.Context, id int64) error
}

func (m *mockPriorityManager) Create(ctx context.Context, e priority.Priority) (priority.Priority, error) {
	return m.createFn(ctx, e)
}

func (m *mockPriorityManager) GetByID(ctx context.Context, id int64) (priority.Priority, error) {
	return m.getFn(ctx, id)
}

func (m *mockPriorityManager) List(ctx context.Context) ([]priority.Priority, error) {
	return m.listFn(ctx)
}

func (m *mockPriorityManager) Update(ctx context.Context, e priority.Priority) (priority.Priority, error) {
	return m.updateFn(ctx, e)
}

func (m *mockPriorityManager) Delete(ctx context.Context, id int64) error { return m.deleteFn(ctx, id) }

func jsonBody(t *testing.T, v any) *bytes.Buffer {
	t.Helper()

	b, err := json.Marshal(v)
	require.NoError(t, err)

	return bytes.NewBuffer(b)
}

func testRequest(t *testing.T, handler http.HandlerFunc, method, url string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var buf *bytes.Buffer
	if body != nil {
		buf = jsonBody(t, body)
	} else {
		buf = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, url, buf).
		WithContext(user.WithContext(t.Context(), user.User{
			Role: role.Admin,
		}))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	return rr
}
