package web

import (
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

func TestCreateTeam(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{
			createFn: func(ctx context.Context, e dict.Entity) (dict.Entity, error) {
				return dict.Entity{
					ID:          1,
					Name:        e.Name,
					Code:        e.Code,
					Description: e.Description,
					Sort:        e.Sort,
				}, nil
			},
		}
		req := CreateTeamRequest{
			Name:        "Test",
			Code:        "TST",
			Description: "desc",
			Sort:        42,
		}
		rr := testRequest(t, createTeam(m), http.MethodPost, "/dummy", req)
		require.Equal(t, http.StatusOK, rr.Code)

		var resp Team

		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		require.Equal(t, int64(1), resp.ID)
		require.Equal(t, "desc", resp.Description)
		require.Equal(t, 42, resp.Sort)
	})

	t.Run("validation error", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{}
		req := CreateTeamRequest{Name: "", Code: ""}
		rr := testRequest(t, createTeam(m), http.MethodPost, "/dummy", req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{
			createFn: func(ctx context.Context, e dict.Entity) (dict.Entity, error) {
				return dict.Entity{}, errors.New("create failed")
			},
		}
		req := CreateTeamRequest{Name: "Name", Code: "C", Sort: 123}
		rr := testRequest(t, createTeam(m), http.MethodPost, "/dummy", req)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestGetTeam(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{
			getFn: func(ctx context.Context, id int64) (dict.Entity, error) {
				return dict.Entity{ID: id, Name: "A"}, nil
			},
		}
		r := chi.NewRouter()
		r.Get("/entities/{id}", getTeam(m))

		req := httptest.NewRequest(http.MethodGet, "/entities/1", nil)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Code)

		var resp Team

		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		require.Equal(t, int64(1), resp.ID)
	})

	t.Run("invalid id", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{}
		r := chi.NewRouter()

		r.Get("/entities/{id}", getTeam(m))

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
		r.Get("/entities/{id}", getTeam(m))

		req := httptest.NewRequest(http.MethodGet, "/entities/1", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestUpdateTeam(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{
			updateFn: func(ctx context.Context, e dict.Entity) (dict.Entity, error) { return e, nil },
		}
		r := chi.NewRouter()
		r.Put("/entities/{id}", updateTeam(m))

		reqBody := UpdateTeamRequest{Name: "Updated", Code: "UPD", Sort: 123}
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
		r.Put("/entities/{id}", updateTeam(m))

		reqBody := UpdateTeamRequest{Name: "", Code: ""}
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
		r.Put("/entities/{id}", updateTeam(m))

		reqBody := UpdateTeamRequest{Name: "X", Code: "X"}
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
		r.Put("/entities/{id}", updateTeam(m))

		reqBody := UpdateTeamRequest{Name: "X", Code: "X", Sort: 123}
		req := httptest.NewRequest(http.MethodPut, "/entities/1", jsonBody(t, reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestDeleteTeam(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		deleted := false
		m := &mockManager{deleteFn: func(ctx context.Context, id int64) error { deleted = true; return nil }}
		r := chi.NewRouter()
		r.Delete("/entities/{id}", deleteTeam(m))

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
		r.Delete("/entities/{id}", deleteTeam(m))

		req := httptest.NewRequest(http.MethodDelete, "/entities/abc", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{deleteFn: func(ctx context.Context, id int64) error { return errors.New("fail") }}
		r := chi.NewRouter()
		r.Delete("/entities/{id}", deleteTeam(m))

		req := httptest.NewRequest(http.MethodDelete, "/entities/1", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestListTeams(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{
			listFn: func(ctx context.Context) ([]dict.Entity, error) {
				return []dict.Entity{
					{ID: 1, Name: "A", Code: "A", Sort: 1},
					{ID: 2, Name: "B", Code: "B", Sort: 2},
				}, nil
			},
		}
		rr := testRequest(t, listTeams(m), http.MethodGet, "/dummy", nil)
		require.Equal(t, http.StatusOK, rr.Code)

		var resp ListTeamsResponse

		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		require.Len(t, resp.Items, 2)
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{
			listFn: func(ctx context.Context) ([]dict.Entity, error) { return []dict.Entity{}, nil },
		}
		rr := testRequest(t, listTeams(m), http.MethodGet, "/dummy", nil)
		require.Equal(t, http.StatusOK, rr.Code)

		var resp ListTeamsResponse

		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		require.Len(t, resp.Items, 0)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockManager{
			listFn: func(ctx context.Context) ([]dict.Entity, error) { return nil, errors.New("fail") },
		}
		rr := testRequest(t, listTeams(m), http.MethodGet, "/dummy", nil)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
