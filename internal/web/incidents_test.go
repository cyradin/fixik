package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cyradin/fixik/internal/dict"
	"github.com/cyradin/fixik/internal/incident"
	"github.com/cyradin/fixik/internal/status"
	"github.com/cyradin/fixik/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestCreateIncident(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockIncidentManager{
			createFn: func(ctx context.Context, i incident.CreateIncident) (incident.Incident, error) {
				return testIncident(i.Title, i.Description), nil
			},
		}

		req := CreateIncidentRequest{
			Title:       "db down",
			Description: "database unavailable",
			StatusID:    1,
			PriorityID:  1,
		}

		rr := testRequest(t, createIncident(m), http.MethodPost, "/incidents", req)
		require.Equal(t, http.StatusOK, rr.Code)

		var resp IncidentResponse

		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		require.Equal(t, "db down", resp.Title)
	})

	t.Run("validation error", func(t *testing.T) {
		t.Parallel()

		m := &mockIncidentManager{}
		req := CreateIncidentRequest{}

		rr := testRequest(t, createIncident(m), http.MethodPost, "/incidents", req)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockIncidentManager{
			createFn: func(ctx context.Context, i incident.CreateIncident) (incident.Incident, error) {
				return incident.Incident{}, errors.New("fail")
			},
		}

		req := CreateIncidentRequest{
			Title:       "db down",
			Description: "database unavailable",
			StatusID:    1,
			PriorityID:  1,
		}

		rr := testRequest(t, createIncident(m), http.MethodPost, "/incidents", req)
		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestGetIncident(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockIncidentManager{
			getFn: func(ctx context.Context, id int64) (incident.Incident, error) {
				i := testIncident("test", "desc")
				i.ID = id

				return i, nil
			},
		}

		r := chi.NewRouter()
		r.Get("/incidents/{id}", getIncident(m))

		req := httptest.NewRequest(http.MethodGet, "/incidents/1", nil)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var resp IncidentResponse

		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		require.Equal(t, int64(1), resp.ID)
	})

	t.Run("invalid id", func(t *testing.T) {
		t.Parallel()

		m := &mockIncidentManager{}
		r := chi.NewRouter()
		r.Get("/incidents/{id}", getIncident(m))

		req := httptest.NewRequest(http.MethodGet, "/incidents/abc", nil)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockIncidentManager{
			getFn: func(ctx context.Context, id int64) (incident.Incident, error) {
				return incident.Incident{}, errors.New("fail")
			},
		}

		r := chi.NewRouter()
		r.Get("/incidents/{id}", getIncident(m))

		req := httptest.NewRequest(http.MethodGet, "/incidents/1", nil)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestUpdateIncident(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockIncidentManager{
			updateFn: func(ctx context.Context, u incident.UpdateIncident) (incident.Incident, error) {
				i := testIncident("updated", "desc")
				i.ID = u.ID

				return i, nil
			},
		}

		r := chi.NewRouter()
		r.Patch("/incidents/{id}", updateIncident(m))

		title := "updated"

		req := httptest.NewRequest(
			http.MethodPatch,
			"/incidents/1",
			jsonBody(t, UpdateIncidentRequest{Title: &title}),
		)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		t.Parallel()

		m := &mockIncidentManager{}
		r := chi.NewRouter()
		r.Patch("/incidents/{id}", updateIncident(m))

		title := "updated"

		req := httptest.NewRequest(
			http.MethodPatch,
			"/incidents/abc",
			jsonBody(t, UpdateIncidentRequest{Title: &title}),
		)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockIncidentManager{
			updateFn: func(ctx context.Context, u incident.UpdateIncident) (incident.Incident, error) {
				return incident.Incident{}, errors.New("fail")
			},
		}

		r := chi.NewRouter()
		r.Patch("/incidents/{id}", updateIncident(m))

		title := "updated"

		req := httptest.NewRequest(
			http.MethodPatch,
			"/incidents/1",
			jsonBody(t, UpdateIncidentRequest{Title: &title}),
		)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestDeleteIncident(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		deleted := false

		m := &mockIncidentManager{
			deleteFn: func(ctx context.Context, id int64) error {
				deleted = true
				return nil
			},
		}

		r := chi.NewRouter()
		r.Delete("/incidents/{id}", deleteIncident(m))

		req := httptest.NewRequest(http.MethodDelete, "/incidents/1", nil)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		require.True(t, deleted)
	})

	t.Run("invalid id", func(t *testing.T) {
		t.Parallel()

		m := &mockIncidentManager{}
		r := chi.NewRouter()
		r.Delete("/incidents/{id}", deleteIncident(m))

		req := httptest.NewRequest(http.MethodDelete, "/incidents/abc", nil)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockIncidentManager{
			deleteFn: func(ctx context.Context, id int64) error {
				return errors.New("fail")
			},
		}

		r := chi.NewRouter()
		r.Delete("/incidents/{id}", deleteIncident(m))

		req := httptest.NewRequest(http.MethodDelete, "/incidents/1", nil)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestListIncidents(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockIncidentManager{
			listFn: func(ctx context.Context, limit, offset int) (incident.IncidentList, error) {
				return incident.IncidentList{
					Items: []incident.Incident{
						testIncident("a", "a"),
						testIncident("b", "b"),
					},
					Total: 2,
				}, nil
			},
		}

		rr := testRequest(t, listIncidents(m), http.MethodGet, "/incidents?limit=10&offset=0", nil)

		require.Equal(t, http.StatusOK, rr.Code)

		var resp IncidentListResponse

		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)

		require.Len(t, resp.Items, 2)
		require.Equal(t, 2, resp.Pagination.Total)
		require.Equal(t, 10, resp.Pagination.Limit)
		require.Equal(t, 0, resp.Pagination.Offset)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockIncidentManager{
			listFn: func(ctx context.Context, limit, offset int) (incident.IncidentList, error) {
				return incident.IncidentList{}, errors.New("fail")
			},
		}

		rr := testRequest(t, listIncidents(m), http.MethodGet, "/incidents?limit=10&offset=0", nil)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestCreateIncidentComment(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		var createdID int64
		m := &mockCommentManager{
			createFn: func(ctx context.Context, incidentID int64, text string) (incident.Comment, error) {
				createdID = incidentID
				return incident.Comment{
					ID:         1,
					IncidentID: incidentID,
					Author:     testUser(),
					Text:       text,
					CreatedAt:  now(),
					UpdatedAt:  now(),
				}, nil
			},
		}

		r := chi.NewRouter()
		r.Post("/incidents/{id}/comments", createIncidentComment(m))

		reqBody := IncidentCommentCreateRequest{Text: "test comment"}
		buf := new(bytes.Buffer)
		require.NoError(t, json.NewEncoder(buf).Encode(reqBody))

		req := httptest.NewRequest(http.MethodPost, "/incidents/1/comments", buf)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, int64(1), createdID)

		var resp IncidentComment
		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		require.Equal(t, "test comment", resp.Text)
	})

	t.Run("validation error", func(t *testing.T) {
		t.Parallel()

		m := &mockCommentManager{}
		r := chi.NewRouter()
		r.Post("/incidents/{id}/comments", createIncidentComment(m))

		reqBody := IncidentCommentCreateRequest{} // пустой текст
		buf := new(bytes.Buffer)
		require.NoError(t, json.NewEncoder(buf).Encode(reqBody))

		req := httptest.NewRequest(http.MethodPost, "/incidents/1/comments", buf)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("manager error", func(t *testing.T) {
		t.Parallel()

		m := &mockCommentManager{
			createFn: func(ctx context.Context, incidentID int64, text string) (incident.Comment, error) {
				return incident.Comment{}, errors.New("fail")
			},
		}

		r := chi.NewRouter()
		r.Post("/incidents/{id}/comments", createIncidentComment(m))

		reqBody := IncidentCommentCreateRequest{Text: "fail"}
		buf := new(bytes.Buffer)
		require.NoError(t, json.NewEncoder(buf).Encode(reqBody))

		req := httptest.NewRequest(http.MethodPost, "/incidents/1/comments", buf)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestGetIncidentComments(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		m := &mockCommentManager{
			listFn: func(ctx context.Context, incidentID int64, limit, offset int) (incident.CommentList, error) {
				return incident.CommentList{
					Items: []incident.Comment{
						{
							ID:         1,
							IncidentID: incidentID,
							Author:     testUser(),
							Text:       "c1",
							CreatedAt:  now(),
							UpdatedAt:  now(),
						},
						{
							ID:         2,
							IncidentID: incidentID,
							Author:     testUser(),
							Text:       "c2",
							CreatedAt:  now(),
							UpdatedAt:  now(),
						},
					},
					Total: 2,
				}, nil
			},
		}

		r := chi.NewRouter()
		r.Get("/incidents/{id}/comments", getIncidentComments(m))

		req := httptest.NewRequest(http.MethodGet, "/incidents/1/comments?limit=10&offset=0", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var resp IncidentCommentListResponse
		err := json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		require.Len(t, resp.Items, 2)
		require.Equal(t, 2, resp.Pagination.Total)
	})
}

type mockCommentManager struct {
	createFn func(ctx context.Context, incidentID int64, text string) (incident.Comment, error)
	listFn   func(ctx context.Context, incidentID int64, limit, offset int) (incident.CommentList, error)
}

func (m *mockCommentManager) Create(ctx context.Context, incidentID int64, text string) (incident.Comment, error) {
	return m.createFn(ctx, incidentID, text)
}

func (m *mockCommentManager) ListByIncident(ctx context.Context, incidentID int64, limit, offset int) (incident.CommentList, error) {
	return m.listFn(ctx, incidentID, limit, offset)
}

func testUser() user.User {
	return user.User{
		ID:   1,
		Name: "Test User",
	}
}

func now() time.Time {
	return time.Now()
}

type mockIncidentManager struct {
	createFn func(ctx context.Context, i incident.CreateIncident) (incident.Incident, error)
	getFn    func(ctx context.Context, id int64) (incident.Incident, error)
	updateFn func(ctx context.Context, i incident.UpdateIncident) (incident.Incident, error)
	deleteFn func(ctx context.Context, id int64) error
	listFn   func(ctx context.Context, limit, offset int) (incident.IncidentList, error)
}

func (m *mockIncidentManager) Create(ctx context.Context, i incident.CreateIncident) (incident.Incident, error) {
	return m.createFn(ctx, i)
}

func (m *mockIncidentManager) GetByID(ctx context.Context, id int64) (incident.Incident, error) {
	return m.getFn(ctx, id)
}

func (m *mockIncidentManager) Update(ctx context.Context, i incident.UpdateIncident) (incident.Incident, error) {
	return m.updateFn(ctx, i)
}

func (m *mockIncidentManager) Delete(ctx context.Context, id int64) error {
	return m.deleteFn(ctx, id)
}

func (m *mockIncidentManager) List(ctx context.Context, limit, offset int) (incident.IncidentList, error) {
	return m.listFn(ctx, limit, offset)
}

func testIncident(title, desc string) incident.Incident {
	return incident.Incident{
		ID:          1,
		Title:       title,
		Description: desc,
		Status: status.Status{
			ID:   1,
			Code: "open",
			Name: "Open",
		},
		Priority: dict.Entity{
			ID:   1,
			Code: "critical",
			Name: "Critical",
		},
	}
}
