package incident

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cyradin/fixik/pkg/tests"
)

type IncidentRepositorySuite struct {
	tests.PostgresSuite
	repo *Repository
}

func TestIncidentRepositorySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(IncidentRepositorySuite))
}

func (s *IncidentRepositorySuite) SetupTest() {
	s.repo = NewRepository(s.Postgres())
}

func (s *IncidentRepositorySuite) TestCreate() {
	ctx := s.T().Context()
	inc := &Incident{
		Title:       "DB down",
		Description: "Database unavailable",
		Impact:      "high",
		Priority:    "critical",
		Status:      "open",
	}

	err := s.repo.Create(ctx, inc)
	s.Require().NoError(err)
	s.NotZero(inc.ID)
	s.NotZero(inc.CreatedAt)
	s.NotZero(inc.UpdatedAt)
}

func (s *IncidentRepositorySuite) TestGetByID_Found() {
	inc := s.loadIncidents([]Incident{
		{Title: "Service down", Description: "Manual insert", Impact: "medium", Priority: "medium", Status: "open"},
	})[0]

	fromDB, err := s.repo.GetByID(s.T().Context(), inc.ID)
	s.Require().NoError(err)
	s.Equal(inc.ID, fromDB.ID)
	s.Equal(inc.Title, fromDB.Title)
	s.Equal(inc.Description, fromDB.Description)
	s.Equal(inc.Impact, fromDB.Impact)
	s.Equal(inc.Priority, fromDB.Priority)
	s.Equal(inc.Status, fromDB.Status)
}

func (s *IncidentRepositorySuite) TestGetByID_NotFound() {
	_, err := s.repo.GetByID(s.T().Context(), 999999)
	s.Require().Error(err)
}

func (s *IncidentRepositorySuite) TestUpdate_Found() {
	inc := s.loadIncidents([]Incident{
		{Title: "Service unavailable", Description: "Initial description", Impact: "medium", Priority: "medium", Status: "open"},
	})[0]

	oldUpdatedAt := inc.UpdatedAt
	inc.Title = "Updated title"
	inc.Status = "in_progress"

	err := s.repo.Update(s.T().Context(), &inc)
	s.Require().NoError(err)
	s.True(inc.UpdatedAt.After(oldUpdatedAt))

	fromDB, err := s.repo.GetByID(s.T().Context(), inc.ID)
	s.Require().NoError(err)
	s.Equal("Updated title", fromDB.Title)
	s.Equal("in_progress", fromDB.Status)
}

func (s *IncidentRepositorySuite) TestUpdate_NotFound() {
	inc := &Incident{
		ID:     999999,
		Title:  "Nonexistent",
		Status: "open",
	}
	err := s.repo.Update(s.T().Context(), inc)
	s.Require().Error(err)
}

func (s *IncidentRepositorySuite) TestDelete_Found() {
	inc := s.loadIncidents([]Incident{
		{Title: "Temp incident", Description: "To be deleted", Impact: "low", Priority: "low", Status: "open"},
	})[0]

	err := s.repo.Delete(s.T().Context(), inc.ID)
	s.Require().NoError(err)

	fromDB, err := s.repo.GetByID(s.T().Context(), inc.ID)
	s.Require().NoError(err)
	s.Equal(inc.ID, fromDB.ID)
}

func (s *IncidentRepositorySuite) TestDelete_NotFound() {
	err := s.repo.Delete(s.T().Context(), 999999)
	s.Require().NoError(err)
}

func (s *IncidentRepositorySuite) loadIncidents(fixtures []Incident) []Incident {
	s.T().Helper()
	ctx := s.T().Context()

	inserted := make([]Incident, len(fixtures))
	for i, inc := range fixtures {
		row := s.Postgres().QueryRow(ctx,
			`INSERT INTO incidents (title, description, impact, urgency, priority, status)
             VALUES ($1, $2, $3, $4, $5, $6)
             RETURNING id, created_at, updated_at`,
			inc.Title, inc.Description, inc.Impact, inc.Priority, inc.Status,
		)
		err := row.Scan(&inc.ID, &inc.CreatedAt, &inc.UpdatedAt)
		s.Require().NoError(err)

		inserted[i] = inc
	}

	return inserted
}
