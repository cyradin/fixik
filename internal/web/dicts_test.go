package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cyradin/fixik/internal/dict"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestCreateDictEntity(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{
			createFn: func(ctx context.Context, e dict.Entity) (dict.Entity, error) {
				return dict.Entity{ID: 1, Name: e.Name, Code: e.Code, Description: e.Description}, nil
			},
		}
		req := CreateDictEntityRequest{Name: "Test", Code: "TST"}
		rr := testRequest(t, createDictEntity(m), http.MethodPost, "/dummy", req)
		require.Equal(t, http.StatusOK, rr.Code)

		var resp DictEntity

		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		require.Equal(t, int64(1), resp.ID)
	})

	t.Run("validation error", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{}
		req := CreateDictEntityRequest{Name: "", Code: ""}
		rr := testRequest(t, createDictEntity(m), http.MethodPost, "/dummy", req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{
			createFn: func(ctx context.Context, e dict.Entity) (dict.Entity, error) {
				return dict.Entity{}, errors.New("create failed")
			},
		}
		req := CreateDictEntityRequest{Name: "Name", Code: "C"}
		rr := testRequest(t, createDictEntity(m), http.MethodPost, "/dummy", req)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestGetDictEntity(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{
			getFn: func(ctx context.Context, id int64) (dict.Entity, error) {
				return dict.Entity{ID: id, Name: "A"}, nil
			},
		}
		r := chi.NewRouter()
		r.Get("/entities/{id}", getDictEntity(m))

		req := httptest.NewRequest(http.MethodGet, "/entities/1", nil)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Code)

		var resp DictEntity

		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		require.Equal(t, int64(1), resp.ID)
	})

	t.Run("invalid id", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{}
		r := chi.NewRouter()

		r.Get("/entities/{id}", getDictEntity(m))

		req := httptest.NewRequest(http.MethodGet, "/entities/abc", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{
			getFn: func(ctx context.Context, id int64) (dict.Entity, error) {
				return dict.Entity{}, errors.New("not found")
			},
		}
		r := chi.NewRouter()
		r.Get("/entities/{id}", getDictEntity(m))

		req := httptest.NewRequest(http.MethodGet, "/entities/1", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestUpdateDictEntity(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{
			updateFn: func(ctx context.Context, e dict.Entity) (dict.Entity, error) { return e, nil },
		}
		r := chi.NewRouter()
		r.Put("/entities/{id}", updateDictEntity(m))

		reqBody := UpdateDictEntityRequest{Name: "Updated", Code: "UPD"}
		req := httptest.NewRequest(http.MethodPut, "/entities/1", jsonBody(t, reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{}
		r := chi.NewRouter()
		r.Put("/entities/{id}", updateDictEntity(m))

		reqBody := UpdateDictEntityRequest{Name: "", Code: ""}
		req := httptest.NewRequest(http.MethodPut, "/entities/1", jsonBody(t, reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{}
		r := chi.NewRouter()
		r.Put("/entities/{id}", updateDictEntity(m))

		reqBody := UpdateDictEntityRequest{Name: "X", Code: "X"}
		req := httptest.NewRequest(http.MethodPut, "/entities/abc", jsonBody(t, reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{
			updateFn: func(ctx context.Context, e dict.Entity) (dict.Entity, error) {
				return dict.Entity{}, errors.New("fail")
			},
		}
		r := chi.NewRouter()
		r.Put("/entities/{id}", updateDictEntity(m))

		reqBody := UpdateDictEntityRequest{Name: "X", Code: "X"}
		req := httptest.NewRequest(http.MethodPut, "/entities/1", jsonBody(t, reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestDeleteDictEntity(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		deleted := false
		m := &mockManager{deleteFn: func(ctx context.Context, id int64) error { deleted = true; return nil }}
		r := chi.NewRouter()
		r.Delete("/entities/{id}", deleteDictEntity(m))

		req := httptest.NewRequest(http.MethodDelete, "/entities/42", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Code)
		require.True(t, deleted)
	})

	t.Run("invalid id", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{}
		r := chi.NewRouter()
		r.Delete("/entities/{id}", deleteDictEntity(m))

		req := httptest.NewRequest(http.MethodDelete, "/entities/abc", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{deleteFn: func(ctx context.Context, id int64) error { return errors.New("fail") }}
		r := chi.NewRouter()
		r.Delete("/entities/{id}", deleteDictEntity(m))

		req := httptest.NewRequest(http.MethodDelete, "/entities/1", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestListDictEntities(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{
			listFn: func(ctx context.Context) ([]dict.Entity, error) {
				return []dict.Entity{
					{ID: 1, Name: "A", Code: "A"},
					{ID: 2, Name: "B", Code: "B"},
				}, nil
			},
		}
		rr := testRequest(t, listDictEntities(m), http.MethodGet, "/dummy", nil)
		require.Equal(t, http.StatusOK, rr.Code)

		var resp ListDictEntitiesResponse

		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		require.Len(t, resp.Items, 2)
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{
			listFn: func(ctx context.Context) ([]dict.Entity, error) { return []dict.Entity{}, nil },
		}
		rr := testRequest(t, listDictEntities(m), http.MethodGet, "/dummy", nil)
		require.Equal(t, http.StatusOK, rr.Code)

		var resp ListDictEntitiesResponse

		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		require.Len(t, resp.Items, 0)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{
			listFn: func(ctx context.Context) ([]dict.Entity, error) { return nil, errors.New("fail") },
		}
		rr := testRequest(t, listDictEntities(m), http.MethodGet, "/dummy", nil)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

type mockManager struct {
	createFn func(ctx context.Context, e dict.Entity) (dict.Entity, error)
	getFn    func(ctx context.Context, id int64) (dict.Entity, error)
	listFn   func(ctx context.Context) ([]dict.Entity, error)
	updateFn func(ctx context.Context, e dict.Entity) (dict.Entity, error)
	deleteFn func(ctx context.Context, id int64) error
}

func (m *mockManager) Create(ctx context.Context, e dict.Entity) (dict.Entity, error) {
	return m.createFn(ctx, e)
}

func (m *mockManager) GetByID(ctx context.Context, id int64) (dict.Entity, error) {
	return m.getFn(ctx, id)
}

func (m *mockManager) List(ctx context.Context) ([]dict.Entity, error) {
	return m.listFn(ctx)
}

func (m *mockManager) Update(ctx context.Context, e dict.Entity) (dict.Entity, error) {
	return m.updateFn(ctx, e)
}

func (m *mockManager) Delete(ctx context.Context, id int64) error { return m.deleteFn(ctx, id) }

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

	req := httptest.NewRequest(method, url, buf)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	return rr
}
